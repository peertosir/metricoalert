package main

import (
	"time"

	"github.com/peertosir/metricoalert/internal/agent"
)

func main() {
	initData()
	mg := agent.NewMetricsGatherer(
		sendAddress, time.Duration(reportInterval)*time.Second,
		time.Duration(pollInterval)*time.Second,
	)
	mg.RunMetricsGatherer()
}
