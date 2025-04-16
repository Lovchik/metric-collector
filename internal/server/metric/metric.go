package metric

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty" `
	Value *float64 `json:"value,omitempty" `
}

func NewMetric(metricName, metricType, metricValue string) (Metrics, error) {

	switch metricType {
	case "gauge":
		{
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				log.Error(err)
				return Metrics{}, err

			}
			return Metrics{metricName, "gauge", nil, &value}, nil
		}
	case "counter":
		{
			value, err := strconv.ParseInt(metricValue, 0, 64)
			if err != nil {
				log.Error(err)
				return Metrics{}, err
			}
			return Metrics{metricName, "counter", &value, nil}, nil
		}
	default:
		{
			return Metrics{}, errors.New("invalid metric type")
		}
	}
}
