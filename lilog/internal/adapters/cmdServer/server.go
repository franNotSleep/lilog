package cmdServer

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/frannotsleep/lilog/internal/application/core/domain"
	"github.com/frannotsleep/lilog/internal/application/ports"
)

type ConnOptions struct {
	address string
}

type Adapter struct {
	db          ports.DBPort
	app         ports.APIPort
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
	Level        uint8    `json:"level"`
	PID          int32    `json:"pid"`
	Hostname     string   `json:"hostname"`
	Request      request  `json:"req"`
	Response     response `json:"res"`
	ResponseTime int32    `json:"responseTime"`
	Message      string   `json:"msg"`
}

func NewAdapter(db ports.DBPort, address string, app ports.APIPort) Adapter {
	connOpts := ConnOptions{address: address}
	return Adapter{db: db, connOptions: connOpts, app: app}
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
			fmt.Fprint(os.Stderr, err)
			continue
		}

		data := data{}

		if err := json.Unmarshal(buf[:n], &data); err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}

		invoice := domain.NewInvoice(data.Time, data.Level, data.PID, data.Hostname, data.ResponseTime, data.Message, domain.InvoiceRequest(data.Request), domain.InvoiceResponse(data.Response))

		err = a.app.NewInvoice(invoice)

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}

		displayInvoice(invoice)

		//	invoices, err := a.app.GetInvoices(invoice.PID)

		//		if err != nil {
		//			fmt.Fprint(os.Stderr, err)
		//			continue
		//		}
	}
}

func displayInvoices(invs []domain.Invoice) {
	for _, inv := range invs {
		displayInvoice(inv)
	}
}

func displayInvoice(inv domain.Invoice) {
	mapLevelToStr := map[uint8]string{
		10: color.New(color.BgBlue).Sprint("DEBUG"),
		20: color.New(color.BgGreen).Sprint("INFO"),
		30: color.New(color.BgHiGreen).Sprint("NOTICE"),
		40: color.New(color.BgYellow).Sprint("WARN"),
		50: color.New(color.BgRed).Sprint("ERROR"),
		60: color.New(color.BgHiRed).Sprint("CRITIC"),
		70: color.New(color.BgMagenta).Sprint("ALERT"),
		80: color.New(color.BgHiMagenta).Sprint("EMERG"),
	}

	localTime := time.UnixMilli(inv.Time).Local().Format(time.ANSIC)
	level := mapLevelToStr[inv.Level]
	fmt.Printf("\n[[ [%v] %s (%d): %s ]]\n", localTime, level, inv.PID, inv.Message)
	fmt.Printf("|%-14s|%-10s|%-20s|\n", "Response Time", "Method", "URL")
  fmt.Printf("---------------|----------|--------------------|\n")
	fmt.Printf("|%-14d|%-10s|%-20s|\n", inv.ResponseTime, inv.InvoiceRequest.Method, inv.InvoiceRequest.URL)
	color.Green("Response Time: %d (ms)", inv.ResponseTime)
}
