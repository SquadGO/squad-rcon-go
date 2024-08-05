package rcon

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"

	p "github.com/iamalone98/squad-rcon-go/internal/parser"
	"github.com/iamalone98/squad-rcon-go/internal/utils"
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

type Warn p.Warn
type Kick p.Kick
type Message p.Message
type PosAdminCam p.PosAdminCam
type UnposAdminCam p.UnposAdminCam
type SquadCreated p.SquadCreated
type Players p.Players
type Squads p.Squads

type Rcon struct {
	connected       bool
	client          net.Conn
	host            string
	port            string
	password        string
	responseBody    string
	lastCommand     string
	lastDataBuffer  []byte
	executeChan     chan string
	onClose         func(error)
	onData          func(string)
	onWarn          func(Warn)
	onKick          func(Kick)
	onMessage       func(Message)
	onPosAdminCam   func(PosAdminCam)
	onUnposAdminCam func(UnposAdminCam)
	onSquadCreated  func(SquadCreated)
	onListPlayers   func(Players)
	onListSquads    func(Squads)
}

func Dial(host, port, password string) (*Rcon, error) {
	r := &Rcon{
		connected:      false,
		lastDataBuffer: make([]byte, 0),
		executeChan:    make(chan string),
	}

	if err := r.connect(host, port); err != nil {
		return nil, err
	}

	if err := r.auth(password); err != nil {
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

func (r *Rcon) connect(host, port string) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), 5*time.Second)
	if err != nil {
		return fmt.Errorf("Connection error: %w", err)
	}

	r.client = conn
	r.connected = true

	return nil
}

func (r *Rcon) auth(password string) error {
	if _, err := r.client.Write(utils.Encode(serverDataAuth, authPacketID, password)); err != nil {
		return fmt.Errorf("Authorization error: %w", err)
	}

	return nil
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

	if r.onClose != nil {
		r.onClose(err)
	} else {
		log.Fatalln(err)
	}

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

			switch data := p.CommandParser(r.responseBody, r.lastCommand).(type) {
			case p.Players:
				{
					if r.onListPlayers != nil {
						r.onListPlayers(Players(data))
					}
				}
			case p.Squads:
				{
					if r.onListSquads != nil {
						r.onListSquads(Squads(data))
					}
				}
			}

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
				if r.onData != nil {
					r.onData(packet.Body)
				}

				switch data := p.ChatParser(packet.Body).(type) {
				case p.Warn:
					{
						if r.onWarn != nil {
							r.onWarn(Warn(data))
						}
					}
				case p.Kick:
					{
						if r.onKick != nil {
							r.onKick(Kick(data))
						}
					}
				case p.Message:
					{
						if r.onMessage != nil {
							r.onMessage(Message(data))
						}
					}
				case p.PosAdminCam:
					{
						if r.onPosAdminCam != nil {
							r.onPosAdminCam(PosAdminCam(data))
						}
					}
				case p.UnposAdminCam:
					{
						if r.onUnposAdminCam != nil {
							r.onUnposAdminCam(UnposAdminCam(data))
						}
					}
				case p.SquadCreated:
					{
						if r.onSquadCreated != nil {
							r.onSquadCreated(SquadCreated(data))
						}
					}
				}
			}

			r.lastDataBuffer = r.lastDataBuffer[size:]
		}
	}
}
