package dbbc3mcast

import (
	"bytes"
	"encoding/binary"
	"net"
)

const UDP_MAX_PACKET_SIZE = 64 * 1024

type dbbc3DDCMulticastListener struct {
	vals chan Dbbc3DdcMulticast
	done chan struct{}
}

func New(groupAddress string) (*dbbc3DDCMulticastListener, error) {
	done := make(chan struct{})
	vals := make(chan Dbbc3DdcMulticast)

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
		pack := Dbbc3DdcMulticast{}

		expectedSize := binary.Size(pack)

		for {
			n, err := conn.Read(buf)
			if err != nil {
				// TODO backoff
				continue
			}

			if n < expectedSize {
				continue
			}

			reader := bytes.NewReader(buf)
			err = binary.Read(reader, binary.LittleEndian, &pack)
			if err != nil {
				continue
			}
			vals <- pack
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

func (l *dbbc3DDCMulticastListener) Values() chan Dbbc3DdcMulticast {
	return l.vals
}
