package metric

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Metric interface {
	GetName() string
	GetValue() any
}

func NewMetric(metricName, metricType, metricValue string) (Metric, error) {
	if metricType == "gauge" {
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			log.Error(err)
			return nil, err

		}
		return &Gauge{metricName, value}, nil
	} else {
		value, err := strconv.ParseInt(metricValue, 0, 64)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return &Counter{metricName, value}, nil
	}
}

func (g *Gauge) GetName() string {
	return g.Name
}
func (g *Gauge) GetValue() any {
	return g.Value
}

type Gauge struct {
	Name  string
	Value float64
}

func (c *Counter) GetValue() any {
	return c.Value
}
func (c *Counter) GetName() string {
	return c.Name
}

type Counter struct {
	Name  string
	Value int64
}
