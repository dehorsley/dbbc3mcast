package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dehorsley/dbbc3mcast"
)

func main() {
	const address = "224.0.0.255:25000"

	l, err := dbbc3mcast.New(address)

	if err != nil {
		log.Fatal(err)
	}

	for {
		p, err := l.Next()
		if err != nil {
			break
		}
		json, err := json.MarshalIndent(&p, "", "  ")
		if err != nil {
			log.Println("error marshaling packet:", err)
		}

		fmt.Println(string(json))
	}

}
