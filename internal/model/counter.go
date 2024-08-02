package model

import (
	"fmt"
	"strconv"

	"github.com/peertosir/metricoalert/internal/errs"
)

func NewCounterMetric(prevMetric *Metric, metricName, value string) (*Metric, error) {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"%w. counter metric value=%q, but must be 'int64'",
			errs.ErrWrongMetricValueType, value,
		)
	}

	if prevMetric != nil {
		val += *prevMetric.IValue
	}

	return &Metric{
		Name:   metricName,
		Type:   MetricTypeCounter,
		IValue: &val,
	}, nil
}
