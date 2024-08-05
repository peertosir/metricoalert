package repository

import (
	"context"
	"strconv"
	"testing"

	"github.com/peertosir/metricoalert/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestInMemMetricStorage_GetMetricByName(t *testing.T) {
	expectedName := "test name"
	var expectedValue int64 = 3
	expectedType := model.MetricTypeCounter

	repo := NewInMemMetricStorage()
	existingMetric, _ := model.NewCounterMetric(nil, expectedName, strconv.Itoa(int(expectedValue)))
	repo.data = map[string]*model.Metric{
		expectedName: existingMetric,
	}

	gotMetric, err := repo.GetMetricByName(context.Background(), expectedName)
	assert.NoError(t, err)
	assert.Equal(t, expectedName, gotMetric.Name)
	assert.Equal(t, expectedType, gotMetric.Type)
	assert.Equal(t, expectedValue, *gotMetric.IValue)
	assert.Nil(t, gotMetric.FValue, nil)
}

func TestInMemMetricStorage_UpsertMetric(t *testing.T) {
	repo := NewInMemMetricStorage()
	newMetric, _ := model.NewCounterMetric(nil, "test", "1")
	err := repo.UpsertMetric(context.Background(), newMetric)
	assert.NoError(t, err)
}
