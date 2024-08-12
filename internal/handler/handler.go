package handler

import (
	"context"
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/peertosir/metricoalert/internal/errs"
	"github.com/peertosir/metricoalert/internal/model"
)

type MetricsSvc interface {
	UpsertMetric(ctx context.Context, metricName, metricType, metricValue string) error
	GetMetric(ctx context.Context, metricName, metricType string) (string, error)
	GetMetrics(ctx context.Context) ([]model.Metric, error)
}

type MetricHandler struct {
	svc MetricsSvc
}

func NewMetricHandler(svc MetricsSvc) *MetricHandler {
	return &MetricHandler{
		svc: svc,
	}
}

func (h *MetricHandler) GetAllMetricsHTML(w http.ResponseWriter, req *http.Request) {
	metrics, err := h.svc.GetMetrics(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *MetricHandler) GetMetric(w http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	value, err := h.svc.GetMetric(req.Context(), metricName, metricType)
	if errors.Is(err, errs.ErrMetricNotFound) || errors.Is(err, errs.ErrUnknownMetricType) {
		w.WriteHeader(http.StatusNotFound)

		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))
}

func (h *MetricHandler) UpdateMetric(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "plain/text")

	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := h.svc.UpsertMetric(req.Context(), metricName, metricType, metricValue)
	if err != nil {
		if errors.Is(err, errs.ErrUnknownMetricType) ||
			errors.Is(err, errs.ErrWrongMetricValueType) ||
			errors.Is(err, errs.ErrMismatchMetricType) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
