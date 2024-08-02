package app

import (
	"net/http"

	"github.com/peertosir/metricoalert/internal/handler"
	"github.com/peertosir/metricoalert/internal/repository"
	"github.com/peertosir/metricoalert/internal/service"
)

func RunApp() {
	mux := http.NewServeMux()
	inMemStorage := repository.NewInMemMetricStorage()
	svc := service.New(inMemStorage)
	mHandler := handler.NewMetricHandler(svc)
	registerHandlers(mux, mHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

func registerHandlers(
	mux *http.ServeMux, metricsHandler *handler.MetricHandler,
) {
	mux.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", metricsHandler.UpdateMetric)
}
