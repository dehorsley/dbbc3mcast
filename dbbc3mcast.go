package dbbc3mcast

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

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

var ErrUnknownVersion = errors.New("unknown DBBC version")

type DbbcMessage = versions.DbbcMessage

type Dbbc3DDCMulticastListener struct {
	vals chan DbbcMessage
	done chan struct{}
	errs chan error
}

func New(groupAddress string) (*Dbbc3DDCMulticastListener, error) {
	done := make(chan struct{})
	vals := make(chan DbbcMessage)
	errs := make(chan error)

	addr, err := net.ResolveUDPAddr("udp", groupAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)

	if err != nil {
		return nil, err
	}

	data := make(chan []byte)
	dataErrs := make(chan error)
	go func() {
		defer close(data)
		defer close(dataErrs)
		for {
			buf := make([]byte, UDP_MAX_PACKET_SIZE)
			_, err := conn.Read(buf)
			if err != nil && err != io.EOF {
				dataErrs <- err
				return
			}
			if err == io.EOF {
				return
			}
			data <- buf
		}
	}()

	go func() {
		defer conn.Close()
		defer close(vals)
		defer close(errs)

	Loop:
		for {
			select {
			case <-done:
				break Loop
			case err := <-dataErrs:
				if err == nil {
					dataErrs = nil
				}
				errs <- err
			case buf := <-data:
				if buf == nil {
					return
				}
				packetVersion := cstr(buf[0:32])

				fields := strings.Split(packetVersion, ",")
				if len(fields) != 3 {
					errs <- fmt.Errorf("%w: %s", ErrUnknownVersion, packetVersion)
					continue
				}

				msg, ok := versions.Messages[strings.Join(fields[0:2], ",")]
				if !ok {
					errs <- fmt.Errorf("%w: %s", ErrUnknownVersion, packetVersion)
					continue
				}

				err = msg.UnmarshalBinary(buf)
				if err != nil {
					errs <- fmt.Errorf("unpacking msg: %w", err)
					continue
				}
				vals <- msg
			}
		}
	}()

	return &Dbbc3DDCMulticastListener{
		done: done,
		vals: vals,
		errs: errs,
	}, nil
}

func (l *Dbbc3DDCMulticastListener) Stop() {
	close(l.done)
}

func (l *Dbbc3DDCMulticastListener) Next() (DbbcMessage, error) {
	select {
	case <-l.done:
		return nil, io.EOF
	case err := <-l.errs:
		return nil, err
	case msg := <-l.vals:
		return msg, nil
	}
}

func (l *Dbbc3DDCMulticastListener) Chans() (chan DbbcMessage, chan error) {
	return l.vals, l.errs
}
