package service

import (
	"context"
	"errors"

	"github.com/peertosir/metricoalert/internal/errs"
	"github.com/peertosir/metricoalert/internal/model"
)

type MetricsRepository interface {
	UpsertMetric(ctx context.Context, metric *model.Metric) error
	GetMetricByName(ctx context.Context, name string) (*model.Metric, error)
}

type MetricService struct {
	repo MetricsRepository
}

func NewMetricService(repo MetricsRepository) *MetricService {
	return &MetricService{
		repo: repo,
	}
}

func (ms *MetricService) UpsertMetric(ctx context.Context, metricName, metricType, metricValue string) error {
	var m *model.Metric

	prevMetric, err := ms.repo.GetMetricByName(ctx, metricName)
	if err != nil && !errors.Is(err, errs.ErrMetricNotFound) {
		return err
	}

	if model.AvailableMetricTypes.Contains(metricType) &&
		prevMetric != nil &&
		prevMetric.Type != metricType {
		return errs.ErrMismatchMetricType
	}

	switch metricType {
	case model.MetricTypeCounter:
		m, err = model.NewCounterMetric(prevMetric, metricName, metricValue)
		if err != nil {
			return err
		}

	case model.MetricTypeGauge:
		m, err = model.NewGaugeMetric(metricName, metricValue)
		if err != nil {
			return err
		}

	default:
		return errs.ErrUnknownMetricType
	}

	return ms.repo.UpsertMetric(ctx, m)
}
