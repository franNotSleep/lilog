package main

import (
	"github.com/frannotsleep/lilog/internal/adapters/db"
	"github.com/frannotsleep/lilog/internal/adapters/server"
	"github.com/frannotsleep/lilog/internal/application/core/api"
	"github.com/frannotsleep/lilog/internal/application/core/domain"
)

func main() {
	invoice := domain.Invoice{
		Time:     1722088153891,
		Level:    30,
		PID:      108791,
		Hostname: "frannotsleep-on-ubuntu",
		InvoiceRequest: domain.InvoiceRequest{
			Method: "GET",
			URL:    "/ping/warn/100",
			Query:  map[string]string{},
			Params: map[string]string{},
			Headers: map[string]string{
				"host":       "localhost:3032",
				"user-agent": "curl/7.81.0",
				"accept":     "*/*",
			},
			RemoteAddress: "::ffff:127.0.0.1",
			RemotePort:    41760,
		},
		InvoiceResponse: domain.InvoiceResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"x-powered-by":   "Express",
				"content-type":   "application/json; charset=utf-8",
				"content-length": "6",
				"etag":           "W/\"6-uBnwlsJiQ3kuZAgKKWB4aV5ugdE\"",
			},
		},
		ResponseTime: 230,
		Message:      "request completed",
	}
	memDBAdapter := db.NewMemKVSAdapter()
	api := api.NewApplication(memDBAdapter)
	api.NewInvoice("web api", invoice)

	connConfig := server.ConnConfig{Address: "127.0.0.1:4119", AllowedClients: []string{"127.0.0.1:5697"}}
	serverAdapter := server.NewAdapter(api, connConfig)

	serverAdapter.ListenAndServe()
}
