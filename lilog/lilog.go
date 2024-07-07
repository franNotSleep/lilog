package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	ServerNameType              uint8 = iota + 1
	MaxLigLogOptionsPayloadSize       = 1 << 20 // 1 MB
)

var ErrMaxLigLogOptionsPayloadSize = errors.New("maximum payload size exceeded.")

type ServerName string

type LiLogOptions struct {
	ServerName ServerName
}

type Payload interface {
	fmt.Stringer
	// io.WriterTo
	io.ReaderFrom
	Bytes() []byte
}

func (n ServerName) Bytes() []byte  { return []byte(n) }
func (n ServerName) String() string { return string(n) }

func (sn *ServerName) ReadFrom(r io.Reader) (int64, error) {
	var n int64 = 1

	var size uint32
	err := binary.Read(r, binary.BigEndian, &size)

	if err != nil {
		return n, err
	}

	n += 4
	if size > MaxLigLogOptionsPayloadSize {
    log.Printf("invalid payload size: %d", size)
		return n, ErrMaxLigLogOptionsPayloadSize
	}

	buf := make([]byte, size)
	o, err := r.Read(buf)

	if err != nil {
		return n, err
	}

	*sn = ServerName(buf)
	return n + int64(o), nil
}

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
			var typ uint8
			err := binary.Read(c, binary.BigEndian, &typ)

			if err != nil {
				log.Printf("while reading first byte: %v", err)
				return
			}

			var payload Payload
			switch typ {
			case ServerNameType:
				payload = new(ServerName)
			default:
				log.Println(errors.New(fmt.Sprintf("unknown type: %d", typ)))
				return
			}

			_, err = payload.ReadFrom(c)

			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("payload: %v\n", payload)
		}(conn)
	}
}
