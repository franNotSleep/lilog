package main

import (
	"github.com/frannotsleep/lilog/internal/adapters/db"
	"github.com/frannotsleep/lilog/internal/adapters/repl"
	"github.com/frannotsleep/lilog/internal/adapters/server"
	"github.com/frannotsleep/lilog/internal/application/core/api"
)

func main() {
	memDBAdapter := db.NewMemKVSAdapter()
	api := api.NewApplication(memDBAdapter)

  replAdapter, ctx := repl.NewAdapter(api)

	connConfig := server.ConnConfig{Address: "127.0.0.1:4119", AllowedClients: []string{"127.0.0.1:5697"}}
	serverAdapter := server.NewAdapter(api, connConfig, ctx)

	//	server := cmdServer.NewAdapter("127.0.0.1:4119", api)

	// server.Run()
	go serverAdapter.ListenAndServe()
  replAdapter.Run()
}
