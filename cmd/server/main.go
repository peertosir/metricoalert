package main

import (
	"log"

	"github.com/peertosir/metricoalert/internal/app"
)

func main() {
	parseFlags()
	log.Printf("Selected port: %s", listenAddr)
	app.RunApp(listenAddr)
}
