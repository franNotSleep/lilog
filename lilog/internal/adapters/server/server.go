package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
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
			continue
		}

		rt, err := reqType(buf[:n])
		if err != nil {
			log.Println(err)
			continue
		}

		if rt == RTR {
			go a.handleRRQ(buf[:n], conn, clientAddr)
		} else if rt == RTS {
			go a.handleSRQ()
		}
	}
}

func (a Adapter) handle(buf []byte) {
	data := new(data)

	if err := json.Unmarshal(buf, data); err != nil {
		log.Printf("invalid json: %v\n", err)
		return
	}

	invoice := domain.NewInvoice(data.Time, data.Level, data.PID, data.Hostname, data.ResponseTime, data.Message, domain.InvoiceRequest(data.Request), domain.InvoiceResponse(data.Response))
	err := a.api.NewInvoice("ok", invoice)

	if err != nil {
		log.Printf("a.api.NewInvoice(): %v\n", err)
		return
	}
}

func (a *Adapter) handleRRQ(bytes []byte, conn net.PacketConn, clientAddr net.Addr) {
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
			conn.WriteTo(data, clientAddr)
		}
	}

	if rq.KeepListening != 0 && a.findListener(clientAddr) == -1 {
		a.listeners = append(a.listeners, clientAddr)
	}
}

func (a Adapter) handleSRQ() {

}
