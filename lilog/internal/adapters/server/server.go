package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"errors"
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

	return nil
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
			defer conn.Close()

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
					fmt.Printf("Read: %v", err)
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
					go a.handleRRQ(buf[:n], clientAddr)
				} else if rt == RTS {
					go a.handleSRQ(buf[:n], conn)
				}
			}
		}()
	}
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

	conn, err := net.Dial("udp", clientAddr.String())
	if err != nil {
		log.Println(err)
		return
	}

	defer func() { _ = conn.Close() }()

	if rq.OpCode == OpRA {
		invoices, err := a.api.GetInvoices(rq.Server)
		if err != nil {
			return
		}

		for _, invoice := range invoices {
			for i := a.connConfig.Retries; i > 0; i-- {
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

				buf := make([]byte, 4096)
				_ = conn.SetReadDeadline(time.Now().Add(a.connConfig.Timeout))
				_, err = conn.Read(buf)
				if err != nil {
					if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
						continue
					}
					log.Println(err)
					return
				}

				typ, err := reqType(buf)
				if err != nil {
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

				break
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
