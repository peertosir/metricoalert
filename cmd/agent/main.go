package main

import (
	"time"

	"github.com/peertosir/metricoalert/internal/agent"
)

func main() {
	mg := agent.NewMetricsGatherer("http://localhost:8080", 10*time.Second, 2*time.Second, 1*time.Minute)
	mg.RunMetricsGatherer()
}
