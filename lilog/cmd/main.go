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
	backupOut, err := os.OpenFile("logs.json", os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	sqliteAdapter, err := db.NewSqliteAdapter("foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteAdapter.Close()

	backupAdapter := backup.NewBackupAdapter(8 * time.Hour, backupOut)
	api := api.NewApplication(sqliteAdapter, backupAdapter)

	connConfig := server.ConnConfig{Address: "127.0.0.1:4119", AllowedClients: []string{"127.0.0.1:5697"}}
	serverAdapter := server.NewAdapter(api, connConfig)

	if err := serverAdapter.ListenAndServeTLS("serverCert.pem", "serverKey.pem"); err != nil {
		log.Fatal(err)
	}
}
