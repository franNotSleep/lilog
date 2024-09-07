package main

import (
	"log"
	"os"
	"time"

	"github.com/frannotsleep/lilog/internal/adapters/backup"
	"github.com/frannotsleep/lilog/internal/adapters/db"
	"github.com/frannotsleep/lilog/internal/adapters/server"
	"github.com/frannotsleep/lilog/internal/application/core/api"
)

func main() {
	backupOut, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	memDBAdapter := db.NewMemKVSAdapter()
	backupAdapter := backup.NewBackupAdapter(5*time.Second, backupOut)
	api := api.NewApplication(memDBAdapter, backupAdapter)

	connConfig := server.ConnConfig{Address: "127.0.0.1:4119", AllowedClients: []string{"127.0.0.1:5697"}}
	serverAdapter := server.NewAdapter(api, connConfig)

	serverAdapter.ListenAndServeTLS("serverCert.pem", "serverKey.pem")
}
