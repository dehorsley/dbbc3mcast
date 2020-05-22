package dbbc3mcast

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/dehorsley/dbbc3mcast/versions"
	_ "github.com/dehorsley/dbbc3mcast/versions/all"
)

func cstr(str []byte) string {
	for n, b := range str {
		if b == 0 {
			return string(str[:n])
		}
	}
	return string(str)
}

const UDP_MAX_PACKET_SIZE = 64 * 1024

type DbbcMessage = versions.DbbcMessage

type dbbc3DDCMulticastListener struct {
	vals chan DbbcMessage
	done chan struct{}
}

func New(groupAddress string) (*dbbc3DDCMulticastListener, error) {
	done := make(chan struct{})
	vals := make(chan DbbcMessage)

	addr, err := net.ResolveUDPAddr("udp", groupAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)

	if err != nil {
		return nil, err
	}

	go func() {
		defer conn.Close()

		buf := make([]byte, UDP_MAX_PACKET_SIZE)

	Loop:
		for {
			select {
			case <-done:
				break Loop
			default:
				n, err := conn.Read(buf)
				if err != nil {
					// TODO backoff
					continue
				}
				packetVersion := cstr(buf[0:32])

				msg, ok := versions.Messages[packetVersion]
				if !ok {
					fmt.Println("unsupported version", packetVersion)
				}

				expectedSize := binary.Size(&msg)

				if n < expectedSize {
					// TODO: handle error?
					continue
				}

				reader := bytes.NewReader(buf)
				err = binary.Read(reader, binary.LittleEndian, &msg)
				if err != nil {
					continue
				}
				vals <- msg
			}
		}
	}()

	return &dbbc3DDCMulticastListener{
		done: done,
		vals: vals,
	}, nil
}

func (l *dbbc3DDCMulticastListener) Stop() {
	close(l.done)
}

func (l *dbbc3DDCMulticastListener) Values() chan DbbcMessage {
	return l.vals
}
