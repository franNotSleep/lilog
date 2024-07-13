package main

import (
	"github.com/frannotsleep/lilog/internal/adapters/cmdServer"
	"github.com/frannotsleep/lilog/internal/adapters/db"
	"github.com/frannotsleep/lilog/internal/application/core/api"
)

func main() {
	memDBAdapter := db.NewMemKVSAdapter()
  app := api.NewApplication(memDBAdapter)
	server := cmdServer.NewAdapter(memDBAdapter, "127.0.0.1:4119", app)

	server.Run()
}
