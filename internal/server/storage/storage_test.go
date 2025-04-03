package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage(t *testing.T) {
	storage := NewMemStorage()

	assert.Equal(t, 0, len(storage.GetAllMetrics()))

	storage.SetMetric("cpu", 90.5)
	storage.SetMetric("memory", 2048)
	storage.SetMetric("disk", "500GB")

	cpu, exists := storage.GetMetricValueByName("cpu")
	assert.True(t, exists)
	assert.Equal(t, 90.5, cpu)

	memory, exists := storage.GetMetricValueByName("memory")
	assert.True(t, exists)
	assert.Equal(t, 2048, memory)

	disk, exists := storage.GetMetricValueByName("disk")
	assert.True(t, exists)
	assert.Equal(t, "500GB", disk)

	allMetrics := storage.GetAllMetrics()
	assert.Equal(t, 3, len(allMetrics))
	assert.Equal(t, 90.5, allMetrics["cpu"])
	assert.Equal(t, 2048, allMetrics["memory"])
	assert.Equal(t, "500GB", allMetrics["disk"])

	_, exists = storage.GetMetricValueByName("gpu")
	assert.False(t, exists)
}
