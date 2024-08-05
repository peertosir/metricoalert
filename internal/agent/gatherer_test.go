package agent

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGatherMetrics(t *testing.T) {
	wantPollCounter := 1
	wantLenMCounter := 1
	wantLenMGauge := 28
	mg := NewMetricsGatherer("", 2*time.Second, 1*time.Second, 2100*time.Millisecond)
	mg.gatherMetrics()

	assert.Equal(t, wantPollCounter, mg.pollCounter)
	assert.Len(t, mg.dataGauge, wantLenMGauge)
	assert.Len(t, mg.dataCounter, wantLenMCounter)
}

func TestSendMetrics(t *testing.T) {
	wantHandlerHit := 2
	gotHandlerHit := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHandlerHit += 1
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	mg := NewMetricsGatherer(ts.URL, 2*time.Second, 1*time.Second, 2100*time.Millisecond)
	mg.dataCounter["MockMetricC"] = "someC"
	mg.dataGauge["MockMetricG"] = "someG"
	mg.sendMetricsData(context.Background())

	assert.Equal(t, wantHandlerHit, gotHandlerHit)
}
