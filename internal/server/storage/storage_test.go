package storage

import (
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/metric"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage(t *testing.T) {
	storage := NewMemStorage()

	metrics, err := storage.GetAllMetrics()
	if err != nil {
		log.Error(err)
	}
	assert.Equal(t, 0, len(metrics))
	firstValue := 90.5
	secondValue := int64(2048)
	err = storage.SetMetric(metric.Metrics{ID: "cpu", MType: "gauge", Value: &firstValue})
	if err != nil {
		log.Error(err)
	}
	err = storage.SetMetric(metric.Metrics{ID: "memory", MType: "counter", Delta: &secondValue})
	if err != nil {
		log.Error(err)
	}
	cpu, exists := storage.GetMetricValueByName("cpu")
	assert.True(t, exists)
	assert.Equal(t, 90.5, *cpu.Value)

	memory, exists := storage.GetMetricValueByName("memory")
	assert.True(t, exists)
	assert.Equal(t, int64(2048), *memory.Delta)

	allMetrics, err := storage.GetAllMetrics()
	if err != nil {
		log.Error(err)
	}
	assert.Equal(t, 2, len(allMetrics))
	assert.Equal(t, 90.5, *allMetrics["cpu"].Value)
	assert.Equal(t, int64(2048), *allMetrics["memory"].Delta)

	_, exists = storage.GetMetricValueByName("gpu")
	assert.False(t, exists)
}
