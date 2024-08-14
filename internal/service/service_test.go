package service

import (
	"context"
	"testing"

	"github.com/peertosir/metricoalert/internal/errs"
	"github.com/peertosir/metricoalert/internal/model"
	"github.com/peertosir/metricoalert/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestMetricService_UpsertMetricAdd(t *testing.T) {
	tests := []struct {
		name          string
		metricName    string
		metricType    string
		metricValue   string
		expectedError error
	}{
		{
			name:        "success gauge metric",
			metricType:  model.MetricTypeGauge,
			metricName:  "testgm1",
			metricValue: "2.3",
		},
		{
			name:        "success counter metric",
			metricType:  model.MetricTypeCounter,
			metricName:  "testcm1",
			metricValue: "2",
		},
		{
			name:          "wrong type value for gauge metric",
			metricType:    model.MetricTypeGauge,
			metricName:    "test_wrong_type_g",
			metricValue:   "invalidtype",
			expectedError: errs.ErrWrongMetricValueType,
		},
		{
			name:          "wrong type value for counter metric",
			metricType:    model.MetricTypeCounter,
			metricName:    "test_wrong_type_c",
			metricValue:   "invalidtype",
			expectedError: errs.ErrWrongMetricValueType,
		},
		{
			name:          "wrong type for metric",
			metricType:    "wonderfulmetric",
			metricName:    "wrong_type_for_metric",
			metricValue:   "2.1",
			expectedError: errs.ErrUnknownMetricType,
		},
	}

	repo := repository.NewInMemMetricStorage()
	svc := NewMetricService(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.UpsertMetric(context.Background(), tt.metricName, tt.metricType, tt.metricValue)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestMetricService_UpsertMetricUpdate(t *testing.T) {
	tests := []struct {
		name                 string
		metricName           string
		metricType           string
		metricValue          string
		metricValueForUpdate string
		metricTypeForUpdate  string
		expectedError        error
	}{
		{
			name:                 "success gauge metric update",
			metricType:           model.MetricTypeGauge,
			metricTypeForUpdate:  model.MetricTypeGauge,
			metricName:           "testgm1",
			metricValue:          "2.3",
			metricValueForUpdate: "5.1",
		},
		{
			name:                 "success counter metric update",
			metricType:           model.MetricTypeCounter,
			metricTypeForUpdate:  model.MetricTypeCounter,
			metricName:           "testcm1",
			metricValue:          "2",
			metricValueForUpdate: "5",
		},
		{
			name:                "mismatch type for metric update",
			metricType:          model.MetricTypeGauge,
			metricName:          "test_wrong_type_g",
			metricValue:         "12.3",
			metricTypeForUpdate: model.MetricTypeCounter,
			expectedError:       errs.ErrMismatchMetricType,
		},
	}

	repo := repository.NewInMemMetricStorage()
	svc := NewMetricService(repo)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = svc.UpsertMetric(context.Background(), tt.metricName, tt.metricType, tt.metricValue)
			err := svc.UpsertMetric(context.Background(), tt.metricName, tt.metricTypeForUpdate, tt.metricValueForUpdate)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestGetMetricValue(t *testing.T) {
	expectedValue := "2"
	metricName := "some"
	metricType := model.MetricTypeCounter
	repo := repository.NewInMemMetricStorage()
	svc := NewMetricService(repo)
	metric, err := model.NewCounterMetric(nil, metricName, expectedValue)
	assert.NoError(t, err)
	err = repo.UpsertMetric(context.Background(), metric)
	assert.NoError(t, err)

	value, err := svc.GetMetric(context.Background(), metricName, metricType)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestGetNonExistMetricValue(t *testing.T) {
	metricName := "some"
	metricType := model.MetricTypeCounter
	repo := repository.NewInMemMetricStorage()
	svc := NewMetricService(repo)

	value, err := svc.GetMetric(context.Background(), metricName, metricType)
	assert.ErrorIs(t, err, errs.ErrMetricNotFound)
	assert.Empty(t, value)
}

func TestGetMetrics(t *testing.T) {
	repo := repository.NewInMemMetricStorage()
	svc := NewMetricService(repo)
	metric, err := model.NewCounterMetric(nil, "1", "22")
	assert.NoError(t, err)
	err = repo.UpsertMetric(context.Background(), metric)
	assert.NoError(t, err)

	value, err := svc.GetMetrics(context.Background())
	assert.NoError(t, err)
	assert.Len(t, value, 1)
}
