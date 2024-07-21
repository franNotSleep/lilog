package repl

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"

	"github.com/fatih/color"
	"github.com/frannotsleep/lilog/internal/application/core/domain"
	"github.com/frannotsleep/lilog/internal/application/ports"
)

func NewAdapter(api ports.APIPort) (Adapter, context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	return Adapter{api: api, ctx: ctx, cancel: cancel}, ctx
}

func (a Adapter) Run() {
	for {
		var op OP
		var pid int32
		displayOptions(&op, &pid)

		switch {
		case op == RONE:
			println("read one")
		case op == RALL:
			invoices, err := a.api.GetInvoices(pid)

			if err != nil {
				log.Printf("a.api.GetInvoices(): %v", err)
				continue
			}

			displayInvoices(invoices)
		case op == EXIT:
			a.cancel()
			<-a.ctx.Done()
			println("bye bye")
			os.Exit(0)
		default:
			println("invalid op code")
		}
	}
}

func displayOptions(op *OP, pid *int32) {
	options := []struct {
		msg  string
		code OP
	}{
		{msg: "Read All Logs", code: RALL},
		{msg: "Read One Logs", code: RONE},
		{msg: "Exit", code: EXIT},
	}

	fmt.Println("Options: ")
	for _, option := range options {
		fmt.Printf("%15s -> OP Code: %d\n", option.msg, option.code)
	}

	fmt.Scanf("%d %d\n", op, pid)
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
