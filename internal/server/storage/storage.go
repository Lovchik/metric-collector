package storage

import (
	"errors"
	"metric-collector/internal/server/metric"
)

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]any),
	}
}

type Storage interface {
	SetMetric(name string, value any)
	GetMetricValueByName(name string) (any, bool)
	GetAllMetrics() map[string]any
	UpdateMetric()
}

func (m *MemStorage) Set(name string, value any) {
	m.metrics[name] = value
}
func (m *MemStorage) GetValueByName(name string) (any, bool) {
	v, ok := m.metrics[name]
	return v, ok
}

type MemStorage struct {
	metrics map[string]any
}

func (m *MemStorage) GetAll() map[string]any {
	return m.metrics
}

func (m *MemStorage) Update(metr metric.Metric) error {
	if counter, ok := metr.(*metric.Counter); ok {
		lastValue, ok := m.GetValueByName(counter.GetName())
		if !ok {
			value := float64(counter.GetValue().(int64))
			m.Set(counter.GetName(), value)
			return nil
		}
		lastFloat, ok := lastValue.(float64)
		if !ok {
			return errors.New("last value is not a float64")
		}

		value := lastFloat + float64(counter.GetValue().(int64))
		m.Set(counter.GetName(), value)
		return nil
	}

	if gauge, ok := metr.(*metric.Gauge); ok {
		m.Set(gauge.GetName(), gauge.GetValue().(float64))
		return nil
	}

	return errors.New("unknown metric type")
}
