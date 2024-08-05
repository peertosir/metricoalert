package model

import (
	"testing"

	"github.com/peertosir/metricoalert/internal/errs"
	"github.com/peertosir/metricoalert/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func Test_CreateNewCounterMetricPositive(t *testing.T) {
	tests := []struct {
		name          string
		metricName    string
		prevMetric    *Metric
		metricValue   string
		expectedValue int64
		expectedType  string
	}{
		{
			name:          "valid counter without previous",
			metricName:    "test",
			prevMetric:    nil,
			metricValue:   "2",
			expectedValue: 2,
		},
		{
			name:        "valid counter with previous",
			metricName:  "test",
			metricValue: "2",
			prevMetric: &Metric{
				IValue: utils.Ptr[int64](2),
			},
			expectedValue: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newMetric, err := NewCounterMetric(tt.prevMetric, tt.metricName, tt.metricValue)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedValue, *newMetric.IValue)
			assert.Equal(t, tt.metricName, newMetric.Name)
			assert.Equal(t, MetricTypeCounter, newMetric.Type)
		})
	}
}

func Test_CreateNewCounterMetricErr(t *testing.T) {
	wantErr := errs.ErrWrongMetricValueType

	m, err := NewCounterMetric(nil, "test", "invalid")
	assert.ErrorIs(t, err, wantErr)
	assert.Nil(t, m)
}
