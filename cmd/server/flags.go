package main

import (
	"flag"
	"log"
	"os"
)

var listenAddr string

func parseFlags() {
	flag.StringVar(&listenAddr, "a", "localhost:8080", "http listen port")
	flag.Parse()
	if len(flag.Args()) != 0 {
		log.Fatalf("got some unexpected args: %+v", flag.Args())
	}

	if srvAddr := os.Getenv("ADDRESS"); srvAddr != "" {
		listenAddr = srvAddr
	}
}
