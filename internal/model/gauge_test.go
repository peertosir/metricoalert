package model

import (
	"testing"

	"github.com/peertosir/metricoalert/internal/errs"
	"github.com/stretchr/testify/assert"
)

func Test_CreateNewGaugeMetricPositive(t *testing.T) {
	metricName := "test"
	metricValue := "2.3"
	expectedMetricValue := 2.3
	newMetric, err := NewGaugeMetric(metricName, metricValue)
	assert.NoError(t, err)
	assert.Equal(t, expectedMetricValue, *newMetric.FValue)
	assert.Equal(t, metricName, newMetric.Name)
	assert.Equal(t, MetricTypeGauge, newMetric.Type)
}

func Test_CreateNewGaugeMetricErr(t *testing.T) {
	wantErr := errs.ErrWrongMetricValueType

	m, err := NewGaugeMetric("test", "invalid")
	assert.ErrorIs(t, err, wantErr)
	assert.Nil(t, m)
}
