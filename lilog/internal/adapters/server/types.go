package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"
	"strings"

	"github.com/frannotsleep/lilog/internal/application/core/domain"
	"github.com/frannotsleep/lilog/internal/application/ports"
)

type Adapter struct {
	api        ports.APIPort
	connConfig ConnConfig
	ctx        context.Context
	listeners  []net.Addr
}

func (a Adapter) findListener(addr net.Addr) int {
	for i, listenerAddr := range a.listeners {
		if listenerAddr.String() == addr.String() {
			return i
		}
	}
	return -1
}

type ConnConfig struct {
	Address        string
	AllowedClients []string
}

type request struct {
	Method        string            `json:"method"`
	URL           string            `json:"url"`
	Query         map[string]string `json:"query"`
	Params        map[string]string `json:"params"`
	Headers       map[string]string `json:"headers"`
	RemoteAddress string            `json:"remoteAddress"`
	RemotePort    int32             `json:"remotePort"`
}

type response struct {
	StatusCode int32             `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
}

type data struct {
	Time         int64    `json:"time"`
	Level        uint8    `json:"level"`
	PID          int32    `json:"pid"`
	Hostname     string   `json:"hostname"`
	Request      request  `json:"req"`
	Response     response `json:"res"`
	ResponseTime int32    `json:"responseTime"`
	Message      string   `json:"msg"`
}

type ReqType int8

const (
	RTR ReqType = iota
	RTS
)

type opCode uint8

func (o opCode) Bytes() []byte {
	return []byte{uint8(o)}
}

const (
	OpRA opCode = iota + 1
	OpRO
	OpData
	OpAck
	OpErr
)

type SendReq struct {
	OpCode opCode
	Server string
	Data   domain.Invoice
}

func (s *SendReq) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code opCode
	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return errors.New("Invalid Send Request.")
	}

	if code != OpData {
		return errors.New("Invalid Send Request.")
	}
	s.OpCode = code

	server, err := r.ReadString(0)
	if err != nil {
		return errors.New("Invalid Send Request.")
	}

	s.Server = strings.TrimRight(server, "\x00")

	data := new(data)

	if err := json.Unmarshal(r.Bytes(), data); err != nil {
		log.Printf("invalid json: %v\n", err)
		return errors.New("Invalid Send Request.")
	}

	invoice := domain.NewInvoice(data.Time, data.Level, data.PID, data.Hostname, data.ResponseTime, data.Message, domain.InvoiceRequest(data.Request), domain.InvoiceResponse(data.Response))
	s.Data = invoice

	return nil
}

type ReadReq struct {
	OpCode        opCode
	Server        string
	From          uint64
	To            uint64
	KeepListening uint8
}

func (q ReadReq) MarshalBinary() ([]byte, error) {
	cap := 1 + len(q.Server) + 1 + 8 + 8 + 1
	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, q.OpCode)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(q.Server)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, q.From)
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, q.To)
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, q.KeepListening)

	return b.Bytes(), nil
}

func (q *ReadReq) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code opCode

	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OpRA && code != OpRO {
		return errors.New("Invalid Read Request.")
	}
	q.OpCode = code

	server, err := r.ReadString(0)
	if err != nil {
		return errors.New("Invalid Read Request.")
	}
	q.Server = strings.TrimRight(server, "\x00")

	var from uint64
	err = binary.Read(r, binary.BigEndian, &from)
	if err != nil {
		return errors.New("Invalid Read Request.")
	}
	q.From = from

	var to uint64
	err = binary.Read(r, binary.BigEndian, &to)
	if err != nil {
		return errors.New("Invalid Read Request.")
	}
	q.To = to

	var keep uint8
	err = binary.Read(r, binary.BigEndian, &keep)
	if err != nil {
		return errors.New("Invalid Read Request")
	}
	q.KeepListening = keep

	return nil
}

func reqType(b []byte) (ReqType, error) {
	r := bytes.NewReader(b)

	var code opCode
	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return -1, err
	}

	if code == OpRA || code == OpRO {
		b = append(code.Bytes(), b...)
		return RTR, nil
	}

	if code == OpData {
		b = append(code.Bytes(), b...)
		return RTS, nil
	}

	return -1, nil

}
