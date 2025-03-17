package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage(t *testing.T) {
	storage := NewMemStorage()

	assert.Equal(t, 0, len(storage.GetAll()))

	storage.Set("cpu", 90.5)
	storage.Set("memory", 2048)
	storage.Set("disk", "500GB")

	cpu, exists := storage.GetValueByName("cpu")
	assert.True(t, exists)
	assert.Equal(t, 90.5, cpu)

	memory, exists := storage.GetValueByName("memory")
	assert.True(t, exists)
	assert.Equal(t, 2048, memory)

	disk, exists := storage.GetValueByName("disk")
	assert.True(t, exists)
	assert.Equal(t, "500GB", disk)

	allMetrics := storage.GetAll()
	assert.Equal(t, 3, len(allMetrics))
	assert.Equal(t, 90.5, allMetrics["cpu"])
	assert.Equal(t, 2048, allMetrics["memory"])
	assert.Equal(t, "500GB", allMetrics["disk"])

	_, exists = storage.GetValueByName("gpu")
	assert.False(t, exists)
}
