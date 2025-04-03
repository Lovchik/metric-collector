package metric

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Metric interface {
	GetName() string
	GetValue() any
}

func NewMetricFromJSON(metrics Metrics) (Metric, error) {
	if metrics.MType == "gauge" {
		if metrics.Value == nil {
			err := errors.New("empty metric gauge value")
			log.Error(err)
			return nil, err

		}
		return Gauge{metrics.ID, *metrics.Value}, nil
	} else {

		if metrics.Delta == nil {
			err := errors.New("metric delta is nil")
			log.Error(err)
			return nil, err
		}
		return Counter{metrics.ID, *metrics.Delta}, nil
	}
}

func (g Gauge) GetName() string {
	return g.Name
}
func (g Gauge) GetValue() any {
	return g.Value
}

type Gauge struct {
	Name  string
	Value float64
}

func (c Counter) GetValue() any {
	return c.Value
}
func (c Counter) GetName() string {
	return c.Name
}

type Counter struct {
	Name  string
	Value int64
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty" `
	Value *float64 `json:"value,omitempty" `
}

func NewMetric(metricName, metricType, metricValue string) (Metric, error) {

	switch metricType {
	case "gauge":
		{
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				log.Error(err)
				return nil, err

			}
			return Gauge{metricName, value}, nil
		}
	case "counter":
		{
			value, err := strconv.ParseInt(metricValue, 0, 64)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			return Counter{metricName, value}, nil
		}
	default:
		{
			return nil, errors.New("invalid metric type")
		}
	}
}
