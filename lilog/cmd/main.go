package main

import (
	"time"

	"github.com/frannotsleep/lilog/internal/adapters/db"
	"github.com/frannotsleep/lilog/internal/adapters/server"
	"github.com/frannotsleep/lilog/internal/application/core/api"
)

func main() {
	memDBAdapter := db.NewMemKVSAdapter()
	api := api.NewApplication(memDBAdapter, 10 * time.Second)

	connConfig := server.ConnConfig{Address: "127.0.0.1:4119", AllowedClients: []string{"127.0.0.1:5697"}}
	serverAdapter := server.NewAdapter(api, connConfig)

	serverAdapter.ListenAndServeTLS("serverCert.pem", "serverKey.pem")
}
