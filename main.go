package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const UDP_MAX_PACKET_SIZE = 64 * 1024

func main() {
	const ip = "224.0.0.255"
	const port = 25000

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Println(fmt.Errorf("unable resolving address: %s", err))
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Println(fmt.Errorf("unable to start UDP listener: %s", err))
		// This likely means the network isn't ready.  Loop until it comes up:

	conloop:
		for {
			select {
			case <-time.Tick(60 * time.Second):
				conn, err = net.ListenMulticastUDP("udp", nil, addr)
				if err != nil {
					log.Println(fmt.Errorf("unable to start UDP listener: %s", err))
					continue
				}
				break conloop
			}
		}
	}
	defer conn.Close()

	log.Printf("DBBC multicast listening to %s\n", addr.String())

	buf := make([]byte, UDP_MAX_PACKET_SIZE)
	pack := Dbbc3DdcMulticast{}

	expectedSize := binary.Size(pack)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(fmt.Errorf("unable to read from conn: %s", err))
			continue
		}

		if n < expectedSize {
			log.Println(fmt.Errorf("packet wrong size, recvd %d, expected %d", n, expectedSize))
		}

		reader := bytes.NewReader(buf)
		err = binary.Read(reader, binary.LittleEndian, &pack)
		if err != nil {
			log.Println(fmt.Errorf("unable to unpack DBBC packet: %s", err))
			continue
		}

		json, err := json.MarshalIndent(&pack, "", "  ")
		if err != nil {
			log.Println("error marshaling packet:", err)
		}

		fmt.Println(string(json))
	}
}
