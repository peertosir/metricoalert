package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var (
	sendAddress    string
	reportInterval int
	pollInterval   int
)

func initData() {
	flag.StringVar(&sendAddress, "a", "localhost:8080", "http listen port")
	flag.IntVar(&reportInterval, "r", 10, "send metrics interval")
	flag.IntVar(&pollInterval, "p", 2, "metrics polling interval")
	flag.Parse()
	if len(flag.Args()) != 0 {
		log.Fatalf("got some unexpected args: %+v", flag.Args())
	}
	if srvAddr := os.Getenv("ADDRESS"); srvAddr != "" {
		sendAddress = srvAddr
	}

	if rInterval := os.Getenv("REPORT_INTERVAL"); rInterval != "" {
		value, err := strconv.Atoi(rInterval)
		if err != nil {
			log.Fatal("wrong value for REPORT_INTERVAL variable")
		}
		reportInterval = value
	}

	if pInterval := os.Getenv("pollInterval"); pInterval != "" {
		value, err := strconv.Atoi(pInterval)
		if err != nil {
			log.Fatal("wrong value for pollInterval variable")
		}
		pollInterval = value
	}
}
