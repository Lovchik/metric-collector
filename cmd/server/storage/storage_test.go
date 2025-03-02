package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage(t *testing.T) {
	NewMemStorage()

	// Проверка пустого хранилища
	assert.Equal(t, 0, len(Store.GetAll()))

	// Добавление значений
	Store.Set("cpu", 90.5)
	Store.Set("memory", 2048)
	Store.Set("disk", "500GB")

	// Проверка получения отдельных значений
	cpu, exists := Store.GetValueByName("cpu")
	assert.True(t, exists)
	assert.Equal(t, 90.5, cpu)

	memory, exists := Store.GetValueByName("memory")
	assert.True(t, exists)
	assert.Equal(t, 2048, memory)

	disk, exists := Store.GetValueByName("disk")
	assert.True(t, exists)
	assert.Equal(t, "500GB", disk)

	// Проверка получения всех значений
	allMetrics := Store.GetAll()
	assert.Equal(t, 3, len(allMetrics))
	assert.Equal(t, 90.5, allMetrics["cpu"])
	assert.Equal(t, 2048, allMetrics["memory"])
	assert.Equal(t, "500GB", allMetrics["disk"])

	// Проверка отсутствующего значения
	_, exists = Store.GetValueByName("gpu")
	assert.False(t, exists)
}
