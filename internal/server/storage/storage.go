package storage

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/metric"
	"os"
	"strconv"
	"strings"
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

func (m *MemStorage) Update(metr metric.Metric) (metric.Metric, error) {
	if counter, ok := metr.(metric.Counter); ok {
		lastValue, ok := m.GetValueByName(counter.GetName())
		if !ok {
			value := counter.GetValue().(int64)
			m.Set(counter.GetName(), value)
			return counter, nil
		}
		lastInt, ok := lastValue.(int64)
		if !ok {
			return nil, errors.New("last value is not a float64")
		}

		value := lastInt + (counter.GetValue().(int64))
		m.Set(counter.GetName(), value)
		counter.Value = value
		return counter, nil
	}

	if gauge, ok := metr.(metric.Gauge); ok {
		m.Set(gauge.GetName(), gauge.GetValue().(float64))
		return gauge, nil
	}

	return nil, errors.New("unknown metric type")
}

func (m *MemStorage) SaveMemoryInfo(filename string) error {
	all := m.GetAll()

	for name, value := range all {
		err := saveMapEntryToFile(filename, name, value)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil

}

func saveMapEntryToFile(filename, key string, value interface{}) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	existingData := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			existingData[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	var valueStr string
	switch v := value.(type) {
	case int64:
		valueStr = strconv.FormatInt(v, 10)
	case float64:
		valueStr = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		valueStr = v
	default:
		return fmt.Errorf("unsupported value type")
	}

	existingData[key] = valueStr

	file.Seek(0, 0)
	file.Truncate(0)
	writer := bufio.NewWriter(file)
	for k, v := range existingData {
		fmt.Fprintf(writer, "%s=%s\n", k, v)
	}
	writer.Flush()

	return nil
}

func getMetricsFromFile(filename string) ([]metric.Metric, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metrics []metric.Metric
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := parts[0]
		valueStr := parts[1]

		if intValue, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			metrics = append(metrics, metric.Counter{Name: name, Value: intValue})
		} else if floatValue, err := strconv.ParseFloat(valueStr, 64); err == nil {
			metrics = append(metrics, metric.Gauge{Name: name, Value: floatValue})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (m *MemStorage) LoadMetricsInMemory(filename string) error {
	metrics, err := getMetricsFromFile(filename)
	if err != nil {
		return err
	}
	for _, metr := range metrics {
		update, err := m.Update(metr)
		log.Info(update)
		if err != nil {
			return err
		}
	}
	all := m.GetAll()
	log.Info(all)
	return nil

}
