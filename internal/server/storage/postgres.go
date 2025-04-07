package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/config"
	"metric-collector/internal/server/metric"
)

type PostgresStorage struct {
	Conn *pgx.Conn
}

func (p PostgresStorage) SetMetric(metric metric.Metrics) error {
	tx, err := p.Conn.Begin(context.Background())
	if err != nil {
		log.Error("Error starting transaction: ", err)
		return err
	}
	exec, err := tx.Exec(context.Background(), "delete from metrics where id = $1 ", metric.ID)
	log.Info("err", exec)
	if err != nil {
		log.Error(err)
		err := tx.Rollback(context.Background())
		if err != nil {
			log.Error(err)
			return err
		}
	}

	var id string
	err = tx.QueryRow(context.Background(), "INSERT INTO metrics (id, type, value, delta) VALUES ($1,$2,$3,$4) RETURNING id", metric.ID, metric.MType, metric.Value, metric.Delta).Scan(&id)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func (p PostgresStorage) GetMetricValueByName(name string) (metric.Metrics, bool) {
	var metrics metric.Metrics
	var count int64
	err := p.Conn.QueryRow(context.Background(), "select count(metrics) from metrics where id = $1", name).Scan(&count)
	if err != nil {
		return metrics, false
	}
	if count == 0 {
		return metrics, false
	}
	err = p.Conn.QueryRow(context.Background(), "select * from metrics where id = $1", name).Scan(&metrics)
	if err != nil {
		return metrics, false
	}
	return metrics, true
}

func (p PostgresStorage) GetAllMetrics() (map[string]metric.Metrics, error) {
	var metricMap = make(map[string]metric.Metrics)
	var metrics []metric.Metrics
	err := p.Conn.QueryRow(context.Background(), "select * from metrics").Scan(&metrics)
	if err != nil {
		log.Errorf("Error getting metrics: %v", err)
	}
	for _, m := range metrics {
		metricMap[m.ID] = m
	}
	return metricMap, nil
}

func (p PostgresStorage) UpdateMetric(metr metric.Metrics) (metric.Metrics, error) {
	switch metr.MType {
	case "counter":
		{
			var lastValue metric.Metrics
			err := p.Conn.QueryRow(context.Background(), "SELECT  FROM metrics WHERE id = $1", metr.ID).Scan(&lastValue)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					err := p.SetMetric(metr)
					if err != nil {
						return metric.Metrics{}, err
					}
					return metr, nil
				}
				return metric.Metrics{}, err
			}

			*metr.Delta += *lastValue.Delta

			err = p.SetMetric(metr)
			if err != nil {
				return metric.Metrics{}, err
			}
			return metr, nil
		}
	case "gauge":
		{
			err := p.SetMetric(metr)
			if err != nil {
				return metric.Metrics{}, err
			}
			return metr, nil
		}
	default:
		{
			return metric.Metrics{}, errors.New("Invalid metric type ")
		}
	}
}

func (p PostgresStorage) LoadMetricsInMemory(filename string) error {
	metrics, err := getMetricsFromFile(filename)
	if err != nil {
		return err
	}
	for _, metr := range metrics {
		err = p.SetMetric(metr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p PostgresStorage) SaveMemoryInfo(filename string) error {
	metrics, err := p.GetAllMetrics()
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

func HealthCheck() error {
	conn, err := pgx.Connect(context.Background(), config.GetConfig().DatabaseDNS)
	if err != nil {
		log.Error("Failed connection to database", err)
		return err
	}
	defer conn.Close(context.Background())
	return nil
}

func NewPgStorage() (*PostgresStorage, error) {
	conn, err := pgx.Connect(context.Background(), config.GetConfig().DatabaseDNS)
	if err != nil {
		log.Error("Unable to connect to database: ", err)
		return nil, err
	}
	defer conn.Close(context.Background())

	return &PostgresStorage{
		conn,
	}, nil
}
