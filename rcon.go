package rcon

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/iamalone98/eventEmitter"

	"github.com/SquadGO/squad-rcon-go/v2/internal/parser"
	"github.com/SquadGO/squad-rcon-go/v2/internal/utils"
	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

const (
	serverDataAuth     = 0x03
	serverDataCommand  = 0x02
	serverDataServer   = 0x01
	serverDataResponse = 0x00

	emptyPacketID    = 100
	authPacketID     = 101
	executeCommandID = 50
)

type RconConfig struct {
	Host               string
	Port               string
	Password           string
	AutoReconnect      bool
	AutoReconnectDelay int
}

type Rcon struct {
	Emitter            eventEmitter.EventEmitter
	connected          bool
	reconnecting       bool
	client             net.Conn
	host               string
	port               string
	password           string
	responseBody       string
	lastCommand        string
	autoReconnect      bool
	autoReconnectDelay int
	lastDataBuffer     []byte
	executeChan        chan string
}

func NewRcon(config RconConfig) (*Rcon, error) {
	c := config
	r := &Rcon{
		Emitter:            eventEmitter.NewEventEmitter(),
		host:               c.Host,
		port:               c.Port,
		password:           c.Password,
		connected:          false,
		lastDataBuffer:     make([]byte, 0),
		executeChan:        make(chan string),
		autoReconnect:      c.AutoReconnect,
		autoReconnectDelay: c.AutoReconnectDelay,
	}

	if err := r.connect(); err != nil {
		return nil, err
	}

	if err := r.auth(); err != nil {
		return nil, err
	}

	go func() {
		r.byteReader()
	}()

	r.ping()

	return r, nil
}

func (r *Rcon) Close() {
	if r.connected {
		r.connected = false

		r.lastCommand = ""
		r.lastDataBuffer = make([]byte, 0)

		close(r.executeChan)
		r.client.Close()

		r.Emitter.Emit(rconEvents.CLOSE, true)

		if r.autoReconnect && r.autoReconnectDelay > 0 {
			r.reconnect(r.autoReconnectDelay)
		}
	}
}

func (r *Rcon) Execute(command string) string {
	r.client.Write(utils.Encode(serverDataCommand, executeCommandID, command))
	r.client.Write(utils.Encode(serverDataCommand, emptyPacketID, ""))

	r.lastCommand = command

	v, ok := <-r.executeChan

	if ok {
		return v
	}

	return ""
}

func (r *Rcon) connect() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", r.host, r.port), 5*time.Second)
	r.reconnecting = false

	if err != nil {
		msg := fmt.Errorf("[RCON] Connection error: %w", err)
		r.Emitter.Emit(rconEvents.ERROR, msg)
		return msg
	}

	r.client = conn
	r.connected = true

	r.Emitter.Emit(rconEvents.CONNECTED, true)

	return nil
}

func (r *Rcon) auth() error {
	if _, err := r.client.Write(utils.Encode(serverDataAuth, authPacketID, r.password)); err != nil {
		msg := fmt.Errorf("[RCON] Authorization error: %w", err)
		r.Emitter.Emit(rconEvents.ERROR, msg)
		return msg
	}

	return nil
}

func (r *Rcon) reconnect(delay int) {
	ticker := time.NewTicker(time.Duration(delay) * time.Second)
	go func() {
	loop:
		for {
			select {
			case <-ticker.C:
				if r.connected {
					break loop
				}

				if !r.reconnecting {
					r.reconnecting = true
					r.connect()
				}
			}
		}
	}()
}

func (r *Rcon) ping() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if r.connected {
					r.Execute("PING_CONNECTION")
				}
			}
		}
	}()
}

func (r *Rcon) byteReader() {
	var err error
	reader := bufio.NewReader(r.client)

	for {
		b, e := reader.ReadByte()
		if e != nil {
			if errors.Is(e, syscall.ECONNRESET) {
				err = fmt.Errorf("[RCON] Error: %w. Check password", e)
			} else if errors.Is(e, syscall.EADDRNOTAVAIL) {
				err = fmt.Errorf("[RCON] Error: %w. Connection lost", e)
			} else {
				err = fmt.Errorf("[RCON] Unknown error: %w", e)
			}

			break
		}

		r.byteParser(b)
	}

	r.Emitter.Emit(rconEvents.ERROR, err)
	r.Close()
}

func (r *Rcon) byteParser(b byte) {
	r.lastDataBuffer = append(r.lastDataBuffer, b)

	if len(r.lastDataBuffer) >= 7 {
		size := int32(binary.LittleEndian.Uint32(r.lastDataBuffer[:4])) + 4

		if r.lastDataBuffer[0] == 0 &&
			r.lastDataBuffer[1] == 1 &&
			r.lastDataBuffer[2] == 0 &&
			r.lastDataBuffer[3] == 0 &&
			r.lastDataBuffer[4] == 0 &&
			r.lastDataBuffer[5] == 0 &&
			r.lastDataBuffer[6] == 0 {

			parser.RconParser(r.responseBody, r.lastCommand, r.Emitter)

			r.lastCommand = ""
			r.executeChan <- r.responseBody
			r.responseBody = ""
			r.lastDataBuffer = make([]byte, 0)
		}

		if int32(len(r.lastDataBuffer)) == size {
			packet := utils.Decode(r.lastDataBuffer)

			if packet.Type == serverDataResponse && packet.ID != authPacketID && packet.ID != emptyPacketID {
				r.responseBody += packet.Body
			}

			if packet.Type == serverDataServer {
				parser.RconParser(packet.Body, r.lastCommand, r.Emitter)
			}

			r.lastDataBuffer = r.lastDataBuffer[size:]
		}
	}
}
