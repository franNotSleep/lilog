package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/frannotsleep/lilog/internal/application/ports"
)

type Adapter struct {
	api        ports.APIPort
	connConfig ConnConfig
	ctx        context.Context
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

type opCode uint16

const (
	OpRA opCode = iota + 1
	OpRO
	OpData
	OpAck
	OpErr
)

type ReadReq struct {
	OpCode opCode
	Server string
	From   int64
	To     int64
}

func (q *ReadReq) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code opCode

	err := binary.Read(r, binary.BigEndian, &code)
	if err != nil {
		return err
	}

	if code != OpRA || code != OpRO {
		return errors.New("Invalid Read Request.")
	}
	q.OpCode = code

	server, err := r.ReadString(0)
	if err != nil {
		return errors.New("Invalid Read Request.")
	}
	q.Server = strings.TrimRight(server, "\x00")

	var from int64
	err = binary.Read(r, binary.BigEndian, &from)
	if err != nil {
		return errors.New("Invalid Read Request.")
	}
	q.From = from

	var to int64
	err = binary.Read(r, binary.BigEndian, &to)
	if err != nil {
		return errors.New("Invalid Read Request.")
	}
	q.To = to

	return nil
}
