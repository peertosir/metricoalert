package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/peertosir/metricoalert/internal/errs"
	"github.com/peertosir/metricoalert/internal/model"
)

type MetricsRepository interface {
	UpsertMetric(ctx context.Context, metric *model.Metric) error
	GetMetricByName(ctx context.Context, name string) (*model.Metric, error)
	GetMetrics(ctx context.Context) ([]model.Metric, error)
}

type MetricService struct {
	repo MetricsRepository
}

func NewMetricService(repo MetricsRepository) *MetricService {
	return &MetricService{
		repo: repo,
	}
}

func (ms *MetricService) GetMetrics(ctx context.Context) ([]model.Metric, error) {
	return ms.repo.GetMetrics(ctx)
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

func (ms *MetricService) GetMetric(ctx context.Context, metricName, metricType string) (string, error) {
	metric, err := ms.repo.GetMetricByName(ctx, metricName)
	if err != nil {
		return "", err
	}

	if metric.Type != metricType {
		return "", errs.ErrMetricNotFound
	}

	switch metricType {
	case model.MetricTypeCounter:
		return strconv.FormatInt(*metric.IValue, 10), nil

	case model.MetricTypeGauge:
		return fmt.Sprintf("%f", *metric.FValue), nil
	default:
		return "", errs.ErrUnknownMetricType
	}
}
