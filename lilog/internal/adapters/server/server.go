package server

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/frannotsleep/lilog/internal/application/ports"
)

func NewAdapter(api ports.APIPort, connConfig ConnConfig) Adapter {
	if connConfig.Retries == 0 {
		connConfig.Retries = 3
	}

	if connConfig.Timeout == 0 {
		connConfig.Timeout = 6 * time.Second
	}

	return Adapter{api: api, connConfig: connConfig, listeners: []net.Addr{}}
}

func (a Adapter) ListenAndServeTLS(certFn, keyFn string) error {
	l, err := net.Listen("tcp", a.connConfig.Address)
	if err != nil {
		return fmt.Errorf("binding to tcp %s: %w", a.connConfig.Address, err)
	}

	log.Printf("Listening on %s...\n", l.Addr().String())
	return a.ServeTLS(l, certFn, keyFn)
}

func (a Adapter) ServeTLS(l net.Listener, certFn, keyFn string) error {
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	cert, err := tls.LoadX509KeyPair(certFn, keyFn)
	if err != nil {
		return fmt.Errorf("loading key pair: %v", err)
	}

	tlsConfig.Certificates = append(tlsConfig.Certificates, cert)

	tlsListener := tls.NewListener(l, tlsConfig)

	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			return fmt.Errorf("accept: %v", err)
		}

		go func() {
			defer func() {
				log.Printf("\033[32m%s\033[0m \033[34mbye...\033[0m👋\n", conn.RemoteAddr())
				_ = conn.Close()
			}()

			for {
				if a.connConfig.Timeout > 0 {
					err := conn.SetDeadline(time.Now().Add(a.connConfig.Timeout))
					if err != nil {
						fmt.Printf("SetDeadline: %v", err)
						return
					}
				}

				buf := make([]byte, 4096)
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Printf("Read: %v\n", err)
					return
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
					go a.handleRRQ(buf[:n], conn)
				} else if rt == RTS {
					go a.handleSRQ(buf[:n], conn)
				}
			}
		}()
	}
}

func (a *Adapter) handleRRQ(bytes []byte, conn net.Conn) {
	rq := ReadReq{}
	rq.UnmarshalBinary(bytes)

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
				log.Println(err)
				return
			}
		}
	}

	log.Printf("\033[32m%s\033[0m is \033[33mwaiting\033[0m for minvoices...😗\n", conn.RemoteAddr())
}

func (a Adapter) handleSRQ(bytes []byte, conn net.Conn) {
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

	data, err := json.Marshal(sq.Data)
	if err != nil {
		log.Println(err)
		return
	}

	conn.Write(data)
}
