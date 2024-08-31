package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/peertosir/metricoalert/internal/handler"
	"github.com/peertosir/metricoalert/internal/logger"
	"github.com/peertosir/metricoalert/internal/middleware"
	"github.com/peertosir/metricoalert/internal/repository"
	"github.com/peertosir/metricoalert/internal/service"
)

func RunApp(port string, logLevel string) {
	logger.InitLogger(logLevel)
	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	inMemStorage := repository.NewInMemMetricStorage()
	svc := service.NewMetricService(inMemStorage)
	mHandler := handler.NewMetricHandler(svc)
	registerHandlers(r, mHandler)

	if err := http.ListenAndServe(port, r); err != nil {
		panic(err)
	}
}

func registerHandlers(
	mux *chi.Mux, metricsHandler *handler.MetricHandler,
) {
	mux.Post(handler.UpdatePath, metricsHandler.UpdateMetric)
	mux.Get(handler.ValuePath, metricsHandler.GetMetric)
	mux.Get(handler.HomePath, metricsHandler.GetAllMetricsHTML)
}
