package errs

import "errors"

var (
	ErrWrongMetricValueType = errors.New("wrong metric value type")
	ErrMetricNotFound       = errors.New("wrong metric value type")
	ErrUnknownMetricType    = errors.New("wrong metric value type")
	ErrMismatchMetricType   = errors.New("mismatch metric type for update")
	ErrEmptyMetricName      = errors.New("wrong metric value type")
)
