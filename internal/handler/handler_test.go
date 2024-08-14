package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/peertosir/metricoalert/internal/model"
	"github.com/peertosir/metricoalert/internal/repository"
	"github.com/peertosir/metricoalert/internal/service"
	"github.com/stretchr/testify/assert"
)

func getMetricUpdateRequestWithCtx(metricName, metricType, metricValue string) *http.Request {
	var requestURL string
	if metricName != "" {
		requestURL = fmt.Sprintf("/update/%s/%s/%s", metricType, metricName, metricValue)
	} else {
		requestURL = fmt.Sprintf("/update/%s/%s", metricType, metricValue)
	}
	request := httptest.NewRequest(http.MethodPost, requestURL, nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("metricType", metricType)
	ctx.URLParams.Add("metricName", metricName)
	ctx.URLParams.Add("metricValue", metricValue)
	return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
}

func TestMetricHandlerAddMetric(t *testing.T) {
	tests := []struct {
		name           string
		metricType     string
		metricName     string
		metricValue    string
		requestURL     string
		wantStatusCode int
	}{
		{
			name:           "success gauge metric",
			metricType:     model.MetricTypeGauge,
			metricName:     "testgm1",
			metricValue:    "2.3",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "success counter metric",
			metricType:     model.MetricTypeCounter,
			metricName:     "testcm1",
			metricValue:    "2",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "wrong type value for gauge metric",
			metricType:     model.MetricTypeGauge,
			metricName:     "test_wrong_type_g",
			metricValue:    "invalidtype",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "wrong type value for counter metric",
			metricType:     model.MetricTypeCounter,
			metricName:     "test_wrong_type_c",
			metricValue:    "invalidtype",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "wrong type for metric",
			metricType:     "wonderfulmetric",
			metricName:     "wrong_type_for_metric",
			metricValue:    "2.1",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "no metric name provided",
			metricType:     model.MetricTypeCounter,
			metricName:     "",
			metricValue:    "2",
			wantStatusCode: http.StatusNotFound,
		},
	}

	repo := repository.NewInMemMetricStorage()
	svc := service.NewMetricService(repo)
	handler := NewMetricHandler(svc)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.UpdateMetric(w, getMetricUpdateRequestWithCtx(tt.metricName, tt.metricType, tt.metricValue))
			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.wantStatusCode, result.StatusCode)
		})
	}
}

func TestMetricHandlerUpdateMetric(t *testing.T) {
	expectedHeaderContentType := "plain/text"
	tests := []struct {
		name              string
		metricName        string
		addMetricType     string
		updateMetricType  string
		addMetricValue    string
		updateMetricValue string
		wantStatusCode    int
	}{
		{
			name:              "success gauge metric",
			addMetricType:     model.MetricTypeGauge,
			updateMetricType:  model.MetricTypeGauge,
			addMetricValue:    "2.2",
			updateMetricValue: "4.3",
			metricName:        "somename1",
			wantStatusCode:    http.StatusOK,
		},
		{
			name:              "success counter metric",
			addMetricType:     model.MetricTypeCounter,
			updateMetricType:  model.MetricTypeCounter,
			addMetricValue:    "2",
			updateMetricValue: "4",
			metricName:        "somename2",
			wantStatusCode:    http.StatusOK,
		},
		{
			name:              "wrong type for update metric",
			addMetricType:     model.MetricTypeGauge,
			updateMetricType:  model.MetricTypeCounter,
			addMetricValue:    "2.2",
			updateMetricValue: "4",
			metricName:        "somename3",
			wantStatusCode:    http.StatusBadRequest,
		},
	}

	repo := repository.NewInMemMetricStorage()
	svc := service.NewMetricService(repo)
	handler := NewMetricHandler(svc)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.UpdateMetric(w, getMetricUpdateRequestWithCtx(tt.metricName, tt.addMetricType, tt.addMetricValue))
			response := w.Result()
			defer response.Body.Close()
			assert.Equal(t, http.StatusOK, response.StatusCode)
			w = httptest.NewRecorder()
			handler.UpdateMetric(w, getMetricUpdateRequestWithCtx(tt.metricName, tt.updateMetricType, tt.updateMetricValue))
			response2 := w.Result()
			defer response2.Body.Close()
			gotStatusCode := response2.StatusCode
			assert.Equal(t, tt.wantStatusCode, gotStatusCode)
			assert.Equal(t, expectedHeaderContentType, response2.Header.Get("Content-Type"))
		})
	}
}
