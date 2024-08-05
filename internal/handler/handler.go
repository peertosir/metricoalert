package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/peertosir/metricoalert/internal/errs"
)

type MetricsSvc interface {
	UpsertMetric(ctx context.Context, metricName, metricType, metricValue string) error
}

type MetricHandler struct {
	svc MetricsSvc
}

func NewMetricHandler(svc MetricsSvc) *MetricHandler {
	return &MetricHandler{
		svc: svc,
	}
}

func (h *MetricHandler) UpdateMetric(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "plain/text")

	reqURLParams := strings.TrimPrefix(req.URL.Path, UpdatePath)
	reqPathParts := strings.Split(strings.TrimSuffix(reqURLParams, "/"), "/")
	if len(reqPathParts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType, metricName, metricValue := reqPathParts[0], reqPathParts[1], reqPathParts[2]

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
