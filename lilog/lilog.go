package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/frannotsleep/lilog/options"
	"github.com/frannotsleep/lilog/types"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:4119")

	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Printf("bound to %q\n", addr)

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			log.Printf("while accepting connection: %v", err)
			return
		}

		go func(c *net.TCPConn) {
			defer c.Close()
			var opts = options.LigLogOptions{}
			var typ uint8
			err := binary.Read(c, binary.BigEndian, &typ)

			if err != nil {
				log.Printf("while reading first byte: %v", err)
				return
			}

			var payload types.Payload
			switch typ {
			case options.ServerNameType:
				payload = new(options.ServerName)
				opts.ServerName = payload.(*options.ServerName)
			default:
				log.Println(errors.New(fmt.Sprintf("unknown type: %d", typ)))
				return
			}

			_, err = payload.ReadFrom(c)

			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("opts: %+v\n", opts)
		}(conn)
	}
}
