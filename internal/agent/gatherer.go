package agent

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/peertosir/metricoalert/internal/model"
)

const (
	PlainTextContentType = "text/plain"
)

type MetricsGatherer struct {
	l                sync.RWMutex
	dataCounter      map[string]string
	dataGauge        map[string]string
	sendTickerTime   time.Duration
	gatherTickerTime time.Duration
	overallTimeout   time.Duration
	pollCounter      int
	metricServerURL  string
}

func NewMetricsGatherer(
	metricServerURL string,
	sendTickerTime, gatherTickerTime, overallTimeout time.Duration,
) *MetricsGatherer {
	return &MetricsGatherer{
		dataCounter:      make(map[string]string),
		dataGauge:        make(map[string]string),
		pollCounter:      0,
		sendTickerTime:   sendTickerTime,
		gatherTickerTime: gatherTickerTime,
		overallTimeout:   overallTimeout,
		metricServerURL:  metricServerURL,
	}
}

func (mg *MetricsGatherer) gatherMetrics() {
	memData := &runtime.MemStats{}
	runtime.ReadMemStats(memData)
	mg.l.Lock()
	defer mg.l.Unlock()
	mg.pollCounter += 1
	mg.dataGauge["Alloc"] = fmt.Sprintf("%v", memData.Alloc)
	mg.dataGauge["BuckHashSys"] = fmt.Sprintf("%v", memData.BuckHashSys)
	mg.dataGauge["Frees"] = fmt.Sprintf("%v", memData.Frees)
	mg.dataGauge["GCCPUFraction"] = fmt.Sprintf("%v", memData.GCCPUFraction)
	mg.dataGauge["GCSys"] = fmt.Sprintf("%v", memData.GCSys)
	mg.dataGauge["HeapAlloc"] = fmt.Sprintf("%v", memData.HeapAlloc)
	mg.dataGauge["HeapIdle"] = fmt.Sprintf("%v", memData.HeapIdle)
	mg.dataGauge["HeapInuse"] = fmt.Sprintf("%v", memData.HeapInuse)
	mg.dataGauge["HeapObjects"] = fmt.Sprintf("%v", memData.HeapObjects)
	mg.dataGauge["HeapReleased"] = fmt.Sprintf("%v", memData.HeapReleased)
	mg.dataGauge["HeapSys"] = fmt.Sprintf("%v", memData.HeapSys)
	mg.dataGauge["LastGC"] = fmt.Sprintf("%v", memData.LastGC)
	mg.dataGauge["Lookups"] = fmt.Sprintf("%v", memData.Lookups)
	mg.dataGauge["MCacheInuse"] = fmt.Sprintf("%v", memData.MCacheInuse)
	mg.dataGauge["MCacheSys"] = fmt.Sprintf("%v", memData.MCacheSys)
	mg.dataGauge["MCacheSys"] = fmt.Sprintf("%v", memData.MSpanInuse)
	mg.dataGauge["MSpanInuse"] = fmt.Sprintf("%v", memData.MSpanInuse)
	mg.dataGauge["MSpanSys"] = fmt.Sprintf("%v", memData.MSpanSys)
	mg.dataGauge["Mallocs"] = fmt.Sprintf("%v", memData.Mallocs)
	mg.dataGauge["NextGC"] = fmt.Sprintf("%v", memData.NextGC)
	mg.dataGauge["NumForcedGC"] = fmt.Sprintf("%v", memData.NumForcedGC)
	mg.dataGauge["NumGC"] = fmt.Sprintf("%v", memData.NumGC)
	mg.dataGauge["OtherSys"] = fmt.Sprintf("%v", memData.OtherSys)
	mg.dataGauge["PauseTotalNs"] = fmt.Sprintf("%v", memData.PauseTotalNs)
	mg.dataGauge["StackInuse"] = fmt.Sprintf("%v", memData.StackInuse)
	mg.dataGauge["StackSys"] = fmt.Sprintf("%v", memData.StackSys)
	mg.dataGauge["Sys"] = fmt.Sprintf("%v", memData.Sys)
	mg.dataGauge["TotalAlloc"] = fmt.Sprintf("%v", memData.TotalAlloc)
	mg.dataGauge["RandomValue"] = fmt.Sprintf("%v", rand.Intn(123124122))
	mg.dataCounter["PollCount"] = fmt.Sprintf("%d", mg.pollCounter)
	fmt.Printf("Gauge metrics: %+v\n", mg.dataGauge)
	fmt.Printf("Counter metrics: %+v\n", mg.dataCounter)
}

func (mg *MetricsGatherer) RunMetricsGatherer() {
	sendTicker := time.NewTicker(mg.sendTickerTime)
	gatherTicker := time.NewTicker(mg.gatherTickerTime)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), mg.overallTimeout)
	defer cancel()
	for {
		select {
		case <-gatherTicker.C:
			go mg.gatherMetrics()
		case <-sendTicker.C:
			go mg.sendMetricsData(timeoutCtx)
		case <-timeoutCtx.Done():
			mg.sendMetricsData(context.Background())
			return
		}
	}
}

func (mg *MetricsGatherer) sendMetricsData(ctx context.Context) {
	fmt.Println("sending metrics data to server")
	mg.l.RLock()
	defer mg.l.RUnlock()
	for name, value := range mg.dataGauge {
		fmt.Printf("[%s] sending %s=%s to server\n", model.MetricTypeGauge, name, value)
		mg.sendMetricsRequest(ctx, model.MetricTypeGauge, name, value)
	}

	for name, value := range mg.dataCounter {
		fmt.Printf("[%s] sending %s=%s to server\n", model.MetricTypeCounter, name, value)
		mg.sendMetricsRequest(ctx, model.MetricTypeCounter, name, value)
	}
}

func (mg *MetricsGatherer) sendMetricsRequest(ctx context.Context, mType, mName, mValue string) {
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost,
		fmt.Sprintf("%s/update/%s/%s/%s", mg.metricServerURL, mType, mName, mValue),
		nil,
	)
	if err != nil {
		log.Fatal("cannot create request to server")
	}
	req.Header.Set("Content-Type", "plain/text")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("cannot send metrics to server")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("sending metrics to server ended with unexpected status code: %d\n", resp.StatusCode)
	}
}
