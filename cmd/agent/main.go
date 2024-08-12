package main

import (
	"time"

	"github.com/peertosir/metricoalert/internal/agent"
)

func main() {
	parseFlags()
	mg := agent.NewMetricsGatherer(
		sendAddress, time.Duration(reportInterval)*time.Second,
		time.Duration(pollInterval)*time.Second,
		1*time.Hour,
	)
	mg.RunMetricsGatherer()
}
