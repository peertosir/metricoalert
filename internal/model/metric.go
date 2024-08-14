package model

import (
	"github.com/peertosir/metricoalert/pkg/utils/datastructs"
)

const (
	MetricTypeCounter = "counter"
	MetricTypeGauge   = "gauge"
)

var AvailableMetricTypes = datastructs.NewHashSetWithValues(MetricTypeCounter, MetricTypeGauge)

type Metric struct {
	Name   string
	Type   string
	IValue *int64
	FValue *float64
}
