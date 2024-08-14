package model

import (
	"fmt"
	"strconv"

	"github.com/peertosir/metricoalert/internal/errs"
)

func NewGaugeMetric(metricName, value string) (*Metric, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("%w. gauge metric value=%q, but must be 'float64'", errs.ErrWrongMetricValueType, value)
	}
	return &Metric{
		Name:   metricName,
		Type:   MetricTypeGauge,
		FValue: &val,
	}, nil
}
