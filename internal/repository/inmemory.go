package repository

import (
	"context"
	"fmt"

	"github.com/peertosir/metricoalert/internal/errs"
	"github.com/peertosir/metricoalert/internal/model"
)

type InMemMetricStorage struct {
	data map[string]*model.Metric
}

func NewInMemMetricStorage() *InMemMetricStorage {
	return &InMemMetricStorage{
		data: make(map[string]*model.Metric),
	}
}

func (ims *InMemMetricStorage) UpsertMetric(_ context.Context, metric *model.Metric) error {
	ims.data[metric.Name] = metric
	return nil
}
func (ims *InMemMetricStorage) GetMetricByName(_ context.Context, name string) (*model.Metric, error) {
	metric, ok := ims.data[name]
	if !ok {
		return nil, fmt.Errorf("%w. metrics with name=%q not found", errs.ErrMetricNotFound, name)
	}
	return metric, nil
}

func (ims *InMemMetricStorage) GetMetrics(_ context.Context) ([]model.Metric, error) {
	returnValue := make([]model.Metric, 0, len(ims.data))

	for _, m := range ims.data {
		returnValue = append(returnValue, *m)
	}

	return returnValue, nil
}
