package main

import (
	"flag"
	"fmt"
)

const (
	READ = "READ"
)

func main() {
	port := flag.Int("port", 4119, "port")
	host := flag.String("host", "127.0.0.1", "host")
	flag.Parse()

	fmt.Printf("Dial: %s:%d\n", *host, *port)
}
