package main

import (
	"log"
	"net"

	"github.com/frannotsleep/lilog/internal/adapters/db"
	"github.com/frannotsleep/lilog/internal/application/core/api"
)

func main() {
	server, err := net.ListenPacket("udp", "127.0.0.1:4119")

	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
  log.Printf("bound to %q", server.LocalAddr())

	memDBAdapter := db.NewMemKVSAdapter()
	api.NewApplication(memDBAdapter)

	buf := make([]byte, 1024)
	for {
		n, _, err := server.ReadFrom(buf)

		if err != nil {
			log.Println(err)
			return
		}
    log.Println(string(buf[:n]))
	}
}
