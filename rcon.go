package rcon

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

var (
	lastDataBuffer = make([]byte, 0)
	resExecuteChan = make(chan string)
	resServerChan  = make(chan string)
	wg             sync.WaitGroup
	host           string
	port           string
	password       string
	responseBody   string
)

type Rcon struct {
	connected bool
	client    net.Conn
}

func Dial(host, port, password string) (*Rcon, error) {
	r := &Rcon{
		connected: false,
	}

	wg.Add(1)

	if err := r.connect(host, port); err != nil {
		fmt.Println("Connection error:", err)

		return nil, err
	}

	fmt.Println("Connection successful")

	if err := r.auth(password); err != nil {
		fmt.Println("Authorization error: ", err)

		return nil, err
	}

	fmt.Println("Authorization successful")

	go func() {
		r.byteReader()
	}()

	r.ping()
	wg.Wait()

	return r, nil
}

func (r *Rcon) Close() error {
	r.connected = false
	close(resExecuteChan)
	close(resServerChan)
	return r.client.Close()
}

func (r *Rcon) Execute(command string) string {
	r.client.Write(utils.Encode(SERVERDATA_COMMAND, EXECUTE_COMMAND_ID, command))
	r.client.Write(utils.Encode(SERVERDATA_COMMAND, EMPTY_PACKET_ID, ""))

	v, ok := <-resExecuteChan

	if ok {
		return v
	}

	return ""
}

func (r *Rcon) OnData(callback func(string)) {
	go func() {
		for data := range resServerChan {
			callback(data)
		}
	}()
}

func (r *Rcon) connect(host, port string) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), 5*time.Second)
	if err != nil {
		return err
	}

	r.client = conn

	return nil
}

func (r *Rcon) auth(password string) error {
	if _, err := r.client.Write(utils.Encode(SERVERDATA_AUTH, AUTH_PACKET_ID, password)); err != nil {
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
	lastDataBuffer = append(lastDataBuffer, b)

	if len(lastDataBuffer) >= 7 {
		size := int32(binary.LittleEndian.Uint32(lastDataBuffer[:4])) + 4

		if lastDataBuffer[0] == 0 &&
			lastDataBuffer[1] == 1 &&
			lastDataBuffer[2] == 0 &&
			lastDataBuffer[3] == 0 &&
			lastDataBuffer[4] == 0 &&
			lastDataBuffer[5] == 0 &&
			lastDataBuffer[6] == 0 {

			resExecuteChan <- responseBody
			responseBody = ""
			lastDataBuffer = make([]byte, 0)
		}

		if int32(len(lastDataBuffer)) == size {
			packet := utils.Decode(lastDataBuffer)
			if packet.Type == SERVERDATA_RESPONSE && packet.ID != AUTH_PACKET_ID && packet.ID != EMPTY_PACKET_ID {
				responseBody += packet.Body
			}

			if packet.Type == SERVERDATA_SERVER {
				resServerChan <- packet.Body
			}

			if packet.ID == AUTH_PACKET_ID && packet.Type == 2 {
				r.connected = true
				wg.Done()
			}

			lastDataBuffer = lastDataBuffer[size:]
		}
	}
}
