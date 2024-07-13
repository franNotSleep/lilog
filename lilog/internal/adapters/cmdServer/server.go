package cmdServer

import (
	"encoding/json"
	"log"
	"net"

	"github.com/frannotsleep/lilog/internal/application/ports"
)

type ConnOptions struct {
	address string
}

type Adapter struct {
	db          ports.DBPort
	connOptions ConnOptions
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
	Level        int32    `json:"level"`
	PID          int32    `json:"pid"`
	Hostname     string   `json:"hostname"`
	Request      request  `json:"req"`
	Response     response `json:"res"`
	ResponseTime int32    `json:"responseTime"`
	Message      string   `json:"message"`
}

func NewAdapter(db ports.DBPort, address string) Adapter {
	connOpts := ConnOptions{address: address}
	return Adapter{db: db, connOptions: connOpts}
}

func (a Adapter) Run() {
	server, err := net.ListenPacket("udp", a.connOptions.address)

	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()
	log.Printf("bound to %q", server.LocalAddr())

	buf := make([]byte, 1024)
	for {
		n, _, err := server.ReadFrom(buf)

		if err != nil {
			log.Println(err)
			return
		}

		data := data{}
		if err := json.Unmarshal(buf[:n], &data); err != nil {
			log.Println(err)
			return
		}

		log.Println(string(buf[:n]))
		log.Printf("%+v\n", data)
	}
}
