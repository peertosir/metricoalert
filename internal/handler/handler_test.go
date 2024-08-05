package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peertosir/metricoalert/internal/model"
	"github.com/peertosir/metricoalert/internal/repository"
	"github.com/peertosir/metricoalert/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestMetricHandlerAddMetric(t *testing.T) {
	tests := []struct {
		name           string
		requestURL     string
		wantStatusCode int
	}{
		{
			name:           "success gauge metric",
			requestURL:     fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeGauge, "testgm1", "2.3"),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "success counter metric",
			requestURL:     fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeCounter, "testcm1", "2"),
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "wrong type value for gauge metric",
			requestURL:     fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeGauge, "test_wrong_type_g", "invalidtype"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "wrong type value for counter metric",
			requestURL:     fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeCounter, "test_wrong_type_c", "invalidtype"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "wrong type for metric",
			requestURL:     fmt.Sprintf("/update/%s/%s/%s", "wonderfulmetric", "wrong_type_for_metric", "2.1"),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "no metric name provided",
			requestURL:     fmt.Sprintf("/update/%s/%s", "wonderfulmetric", "2.1"),
			wantStatusCode: http.StatusNotFound,
		},
	}

	repo := repository.NewInMemMetricStorage()
	svc := service.NewMetricService(repo)
	handler := NewMetricHandler(svc)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.requestURL, nil)
			w := httptest.NewRecorder()
			handler.UpdateMetric(w, request)
			gotStatusCode := w.Code
			assert.Equal(t, tt.wantStatusCode, gotStatusCode)
		})
	}
}

func TestMetricHandlerUpdateMetric(t *testing.T) {
	expectedHeaderContentType := "plain/text"
	tests := []struct {
		name             string
		addRequestURL    string
		updateRequestURL string
		wantStatusCode   int
	}{
		{
			name:             "success gauge metric",
			addRequestURL:    fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeGauge, "testgm1", "2.3"),
			updateRequestURL: fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeGauge, "testgm1", "4.3"),
			wantStatusCode:   http.StatusOK,
		},
		{
			name:             "success counter metric",
			addRequestURL:    fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeGauge, "testcm1", "2"),
			updateRequestURL: fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeGauge, "testcm1", "4"),
			wantStatusCode:   http.StatusOK,
		},
		{
			name:             "wrong type for update metric",
			addRequestURL:    fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeCounter, "testcm1", "2"),
			updateRequestURL: fmt.Sprintf("/update/%s/%s/%s", model.MetricTypeGauge, "testcm1", "2.3"),
			wantStatusCode:   http.StatusBadRequest,
		},
	}

	repo := repository.NewInMemMetricStorage()
	svc := service.NewMetricService(repo)
	handler := NewMetricHandler(svc)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addRequest := httptest.NewRequest(http.MethodPost, tt.addRequestURL, nil)
			w := httptest.NewRecorder()
			handler.UpdateMetric(w, addRequest)
			updateRequest := httptest.NewRequest(http.MethodPost, tt.updateRequestURL, nil)
			handler.UpdateMetric(w, updateRequest)
			defer w.Result().Body.Close()
			gotStatusCode := w.Code
			assert.Equal(t, tt.wantStatusCode, gotStatusCode)
			assert.Equal(t, expectedHeaderContentType, w.Result().Header.Get("Content-Type"))
		})
	}
}
