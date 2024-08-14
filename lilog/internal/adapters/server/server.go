package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	"github.com/frannotsleep/lilog/internal/application/ports"
)

func NewAdapter(api ports.APIPort, connConfig ConnConfig) Adapter {
	return Adapter{api: api, connConfig: connConfig, listeners: []net.Addr{}}
}

func (a Adapter) ListenAndServe() error {
	conn, err := net.ListenPacket("udp", a.connConfig.Address)
	if err != nil {
		return err
	}

	defer func() { _ = conn.Close() }()
	log.Printf("Listening on %s...\n", conn.LocalAddr())

	return a.serve(conn)
}

func (a Adapter) serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("nil connection")
	}

	for {
		buf := make([]byte, 4096)
		n, clientAddr, err := conn.ReadFrom(buf)
		if err != nil {
			continue
		}

		var code opCode

		r := bytes.NewBuffer(buf)
		err = binary.Read(r, binary.BigEndian, &code)
		if err != nil {
			log.Println(err)
			continue
		}

		rt, err := reqType(buf[:n])
		if err != nil {
			log.Println(err)
			continue
		}

		if rt == RTR {
			go a.handleRRQ(buf[:n], clientAddr)
		} else if rt == RTS {
			go a.handleSRQ(buf[:n], conn)
		}
	}
}

func (a *Adapter) handleRRQ(bytes []byte, clientAddr net.Addr) {
	rq := ReadReq{}
	rq.UnmarshalBinary(bytes)

	log.Println(clientAddr.String())
	conn, err := net.Dial("udp", clientAddr.String())
	if err != nil {
		log.Println(err)
		return
	}

	defer func() { _ = conn.Close() }()
	_ = conn.SetReadDeadline(time.Now().Add(a.connConfig.Timeout))

	if rq.OpCode == OpRA {
		invoices, err := a.api.GetInvoices(rq.Server)
		if err != nil {
			return
		}

		for _, invoice := range invoices {
			data, err := json.Marshal(invoice)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = conn.Write(data)
			if err != nil {
				// TODO: apply retries feature
				log.Println(err)
				return
			}

			buf := make([]byte, 4096)
			_, err = conn.Read(buf)
			if err != nil {
				// TODO: apply retries feature
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					// TODO: retries
					log.Println(err)
				}
				log.Println(err)
				return
			}

			typ, err := reqType(buf)
			if err != nil {
				// TODO: apply retries feature
				log.Println(err)
				return
			}

			if typ != RTA {
				errData := Err{}
				errData.Message = "Invalid ACK."
				errData.OpCode = OpErr
				errData.Error = ErrIllegalOp
				b, err := errData.MarshalBinary()
				if err != nil {
					log.Println(err)
					return
				}
				_, err = conn.Write(b)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}

	if rq.KeepListening != 0 && a.findListener(clientAddr) == -1 {
		log.Printf("\033[32m%s\033[0m is waiting for invoices...😗\n", clientAddr)
		a.listeners = append(a.listeners, clientAddr)
	}
}

func (a Adapter) handleSRQ(bytes []byte, conn net.PacketConn) {
	sq := SendReq{}
	err := sq.UnmarshalBinary(bytes)
	if err != nil {
		log.Println(err)
		return
	}

	err = a.api.NewInvoice(sq.Server, sq.Data)
	if err != nil {
		log.Println(err)
		return
	}

	for _, clientAddr := range a.listeners {
		data, err := json.Marshal(sq.Data)
		if err != nil {
			log.Println(err)
			return
		}

		conn.WriteTo(data, clientAddr)
	}
}
