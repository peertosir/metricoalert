package errs

import "errors"

var (
	ErrWrongMetricValueType = errors.New("wrong metric value type")
	ErrMetricNotFound       = errors.New("metric not found")
	ErrUnknownMetricType    = errors.New("unknown metric type")
	ErrMismatchMetricType   = errors.New("mismatch metric type for update")
)
