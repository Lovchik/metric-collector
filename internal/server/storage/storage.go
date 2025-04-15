package storage

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"metric-collector/internal/retry"
	"metric-collector/internal/server/config"
	"metric-collector/internal/server/metric"
	"os"
)

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]metric.Metrics),
	}
}

type Storage interface {
	SetMetric(metric.Metrics) error
	GetMetricValueByName(name string) (metric.Metrics, bool)
	GetAllMetrics() (map[string]metric.Metrics, error)
	UpdateMetric(metric.Metrics) (metric.Metrics, error)
	LoadMetricsInMemory(string) error
	SaveMemoryInfo(string) error
	UpdateMetrics([]metric.Metrics) ([]metric.Metrics, error)
	HealthCheck() (err error)
}

func (m *MemStorage) HealthCheck() error {
	return nil
}

func (m *MemStorage) SetMetric(metric metric.Metrics) error {
	m.metrics[metric.ID] = metric
	return nil
}
func (m *MemStorage) GetMetricValueByName(name string) (metric.Metrics, bool) {
	v, ok := m.metrics[name]
	return v, ok
}

type MemStorage struct {
	metrics map[string]metric.Metrics
}

func (m *MemStorage) GetAllMetrics() (map[string]metric.Metrics, error) {
	return m.metrics, nil
}

func (m *MemStorage) UpdateMetric(metr metric.Metrics) (metric.Metrics, error) {
	if metr.MType == "counter" {
		lastValue, ok := m.GetMetricValueByName(metr.ID)
		if !ok {
			err := m.SetMetric(metr)
			if err != nil {
				return metric.Metrics{}, err
			}
			return metr, nil
		}

		*metr.Delta = *metr.Delta + *lastValue.Delta
		err := m.SetMetric(metr)
		if err != nil {
			return metric.Metrics{}, err
		}
		return metr, nil
	}

	if metr.MType == "gauge" {
		err := m.SetMetric(metr)
		if err != nil {
			return metric.Metrics{}, err
		}
		return metr, nil
	}

	if config.GetConfig().StoreInterval == 0 {
		err := UpdateMetricInFile(metr)
		log.Error(err)
		return metric.Metrics{}, err
	}

	return metric.Metrics{}, errors.New("unknown metric type")
}

func (m *MemStorage) SaveMemoryInfo(filename string) error {
	metrics, err := m.GetAllMetrics()
	if err != nil {
		log.Error(err)
		return err
	}
	all := metrics

	err = saveMapEntryToFile(filename, all)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil

}

func saveMapEntryToFile(filename string, data map[string]metric.Metrics) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		panic(err)
	}

	return nil
}

func getMetricsFromFile(filename string) ([]metric.Metrics, error) {
	file, err := retry.Retry(3, 1, func() (*os.File, error) {
		return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var metricMap map[string]metric.Metrics
	if err := decoder.Decode(&metricMap); err != nil {
		if errors.Is(err, io.EOF) {
			return []metric.Metrics{}, nil
		}
		return nil, err
	}

	metrics := make([]metric.Metrics, 0, len(metricMap))
	for _, m := range metricMap {
		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (m *MemStorage) LoadMetricsInMemory(filename string) error {
	metrics, err := getMetricsFromFile(filename)
	if err != nil {
		return err
	}
	for _, metr := range metrics {
		update, err := m.UpdateMetric(metr)
		log.Info(update)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemStorage) UpdateMetrics(metrics []metric.Metrics) ([]metric.Metrics, error) {
	if m.metrics == nil {
		m.metrics = make(map[string]metric.Metrics)
	}

	for _, newMetric := range metrics {
		existing, exists := m.metrics[newMetric.ID]

		switch newMetric.MType {
		case "counter":
			if exists && existing.Delta != nil && newMetric.Delta != nil {
				sum := *existing.Delta + *newMetric.Delta
				existing.Delta = &sum
				m.metrics[newMetric.ID] = existing
			} else {
				m.metrics[newMetric.ID] = newMetric
			}
		case "gauge":
			m.metrics[newMetric.ID] = newMetric
		default:
			return nil, errors.New("unsupported metric type ")
		}
	}
	if config.GetConfig().StoreInterval == 0 {
		for _, m := range metrics {
			err := UpdateMetricInFile(m)
			log.Error(err)
			return nil, err
		}
	}

	return metrics, nil
}

func UpdateMetricInFile(metr metric.Metrics) error {
	metrics, err := getMetricsFromFile(config.GetConfig().FileStoragePath)
	if err != nil {
		return err
	}

	metricMap := make(map[string]metric.Metrics)
	for _, m := range metrics {
		metricMap[m.ID] = m
	}

	metricMap[metr.ID] = metr

	updatedMetrics := make([]metric.Metrics, 0, len(metricMap))
	for _, m := range metricMap {
		updatedMetrics = append(updatedMetrics, m)
	}

	file, err := os.OpenFile(config.GetConfig().FileStoragePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(updatedMetrics); err != nil {
		return err
	}

	return nil
}
