package main

import (
	"flag"
	"log"
)

var (
	sendAddress    string
	reportInterval int
	pollInterval   int
)

func parseFlags() {
	flag.StringVar(&sendAddress, "a", "http://localhost:8080", "http listen port")
	flag.IntVar(&reportInterval, "r", 10, "send metrics interval")
	flag.IntVar(&pollInterval, "p", 2, "metrics polling interval")
	flag.Parse()
	if len(flag.Args()) != 0 {
		log.Fatalf("got some unexpected args: %+v", flag.Args())
	}
}
