package main

import (
	"flag"
	"log"
)

var listenPort string

func parseFlags() {
	flag.StringVar(&listenPort, "a", "localhost:8080", "http listen port")
	flag.Parse()
	if len(flag.Args()) != 0 {
		log.Fatalf("got some unexpected args: %+v", flag.Args())
	}
}
