package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/iamalone98/squad-rcon-go/utils"
)

const (
	SERVERDATA_AUTH     = 0x03
	SERVERDATA_COMMAND  = 0x02
	SERVERDATA_SERVER   = 0x01
	SERVERDATA_RESPONSE = 0x00

	EMPTY_PACKET_ID = 100
	AUTH_PACKET_ID  = 101

	EXECUTE_COMMAND_ID = 50
)

type Rcon struct {
	host           string
	port           string
	password       string
	responseBody   string
	resExecuteChan chan string
	resServerChan  chan string
	lastDataBuffer []byte
	connected      bool
	client         net.Conn
	wg             sync.WaitGroup
}

func Dial(host, port, password string) (*Rcon, error) {
	r := &Rcon{
		host:           host,
		port:           port,
		password:       password,
		connected:      false,
		lastDataBuffer: make([]byte, 0),
		resExecuteChan: make(chan string),
		resServerChan:  make(chan string),
	}

	r.wg.Add(1)

	if err := r.connect(); err != nil {
		fmt.Println("Connection error:", err)

		return nil, err
	}

	fmt.Println("Connection successful")

	if err := r.auth(); err != nil {
		fmt.Println("Authorization error: ", err)

		return nil, err
	}

	fmt.Println("Authorization successful")

	go func() {
		r.byteReader()
	}()

	r.ping()
	r.wg.Wait()

	return r, nil
}

func (r *Rcon) Close() error {
	r.connected = false
	close(r.resExecuteChan)
	close(r.resServerChan)
	return r.client.Close()
}

func (r *Rcon) Execute(command string) string {
	r.client.Write(utils.Encode(SERVERDATA_COMMAND, EXECUTE_COMMAND_ID, command))
	r.client.Write(utils.Encode(SERVERDATA_COMMAND, EMPTY_PACKET_ID, ""))

	v, ok := <-r.resExecuteChan

	if ok {
		return v
	}

	return ""
}

func (r *Rcon) OnData(callback func(string)) {
	go func() {
		for data := range r.resServerChan {
			callback(data)
		}
	}()
}

func (r *Rcon) connect() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", r.host, r.port), 5*time.Second)
	if err != nil {
		return err
	}

	r.client = conn

	return nil
}

func (r *Rcon) auth() error {
	if _, err := r.client.Write(utils.Encode(SERVERDATA_AUTH, AUTH_PACKET_ID, r.password)); err != nil {
		return err
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
	reader := bufio.NewReader(r.client)

	for {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Connection closed: ", err)
			break
		}

		r.byteParser(b)
	}
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

			r.resExecuteChan <- r.responseBody
			r.responseBody = ""
			r.lastDataBuffer = make([]byte, 0)
		}

		if int32(len(r.lastDataBuffer)) == size {
			packet := utils.Decode(r.lastDataBuffer)
			if packet.Type == SERVERDATA_RESPONSE && packet.ID != AUTH_PACKET_ID && packet.ID != EMPTY_PACKET_ID {
				r.responseBody += packet.Body
			}

			if packet.Type == SERVERDATA_SERVER {
				r.resServerChan <- packet.Body
			}

			if packet.ID == AUTH_PACKET_ID && packet.Type == 2 {
				r.connected = true
				r.wg.Done()
			}

			r.lastDataBuffer = r.lastDataBuffer[size:]
		}
	}
}
