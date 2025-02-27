package metric

import (
	"metric-collector/cmd/server/storage"
	"strconv"
)

type Metric interface {
	Update() error
	GetName() string
	GetValue() any
}

func NewMetric(metricName, metricType, metricValue string) Metric {
	if metricType == "gauge" {
		value, _ := strconv.ParseFloat(metricValue, 64)
		return &Gauge{metricName, value}
	} else {
		value, _ := strconv.ParseInt(metricValue, 0, 64)
		return &Counter{metricName, value}
	}
}

func (g *Gauge) GetName() string {
	return g.Name
}
func (g *Gauge) GetValue() any {
	return g.Value
}

func (g *Gauge) Update() error {
	storage.Store.Set(g.Name, g.GetValue())
	return nil
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

func (c *Counter) Update() error {
	lastValue, ok := storage.Store.GetValueByName(c.Name)

	if !ok {
		return nil
	}
	lastFloat, ok := lastValue.(float64)

	value := lastFloat + float64(c.Value)
	storage.Store.Set(c.GetName(), value)
	return nil
}
