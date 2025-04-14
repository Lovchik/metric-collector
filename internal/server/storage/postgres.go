package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/metric"
	"time"
)

func (p PostgresStorage) SetMetric(metric metric.Metrics) error {
	query := `
		INSERT INTO metrics (id, type, value, delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE 
		SET type = EXCLUDED.type,
		    value = EXCLUDED.value,
		    delta = EXCLUDED.delta;
	`

	exec, err := p.Conn.Exec(context.Background(), query, metric.ID, metric.MType, metric.Value, metric.Delta)
	log.Info("exec ", exec)
	if err != nil {
		log.Error("Insert or update error: ", err)
		return err
	}

	return nil
}

func (p PostgresStorage) GetMetricValueByName(name string) (metric.Metrics, bool) {
	var metrics metric.Metrics
	err := p.Conn.QueryRow(context.Background(), "select * from metrics where id = $1", name).Scan(&metrics.ID, &metrics.MType, &metrics.Value, &metrics.Delta)
	if err != nil {
		log.Error(err)
		return metrics, false
	}
	return metrics, true
}

func (p PostgresStorage) GetAllMetrics() (map[string]metric.Metrics, error) {
	var metricMap = make(map[string]metric.Metrics)
	rows, err := p.Conn.Query(context.Background(), "SELECT id, type, value, delta FROM metrics")
	if err != nil {
		log.Error("Error querying metrics: ", err)
		return metricMap, err
	}
	defer rows.Close()

	for rows.Next() {
		var m metric.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			log.Error("Error scanning row: ", err)
			return metricMap, err
		}

		metricMap[m.ID] = m
	}

	if err = rows.Err(); err != nil {
		log.Error("Error iterating over rows: ", err)
		return metricMap, err
	}

	return metricMap, nil
}

func (p PostgresStorage) UpdateMetric(metr metric.Metrics) (metric.Metrics, error) {
	switch metr.MType {
	case "counter":
		{
			var lastValue metric.Metrics

			err := p.Conn.QueryRow(context.Background(), "SELECT * FROM metrics WHERE id = $1", metr.ID).Scan(&lastValue.ID, &lastValue.MType, &lastValue.Value, &lastValue.Delta)
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
			return metric.Metrics{}, errors.New("invalid metric type ")
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

func (p PostgresStorage) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.Conn.Ping(ctx); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (p PostgresStorage) UpdateMetrics(metrics []metric.Metrics) ([]metric.Metrics, error) {
	tx, err := p.Conn.Begin(context.Background())
	if err != nil {
		log.Error("Error updating metrics: ", err)
		return nil, err
	}
	defer tx.Rollback(context.Background())

	for _, m := range metrics {
		switch m.MType {
		case "counter":
			_, err := tx.Exec(context.Background(), `
				INSERT INTO metrics (id, type, delta)
				VALUES ($1, $2, $3)
				ON CONFLICT (id) DO UPDATE 
				SET delta = metrics.delta + EXCLUDED.delta
			`, m.ID, m.MType, m.Delta)
			if err != nil {
				log.Error("Error updating metrics: ", err)
				return nil, err
			}
		case "gauge":
			_, err := tx.Exec(context.Background(), `
				INSERT INTO metrics (id, type, value)
				VALUES ($1, $2, $3)
				ON CONFLICT (id) DO UPDATE 
				SET value = EXCLUDED.value
			`, m.ID, m.MType, m.Value)
			if err != nil {
				log.Error("Error updating metrics: ", err)
				return nil, err
			}
		default:
			return nil, errors.New("unsupported metric type")
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Error("Error updating metrics: ", err)
		return nil, err
	}

	return metrics, nil
}

type PostgresStorage struct {
	Conn *pgxpool.Pool
}

func NewPgStorage(ctx context.Context, dataBaseDSN string) (*PostgresStorage, error) {
	const op = "internal.repo.storage.postgres.NewPgStorage"

	pool, err := pgxpool.New(ctx, dataBaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool, %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS metrics(
                                    id VARCHAR (50) UNIQUE NOT NULL,
                                    type VARCHAR (50)  NOT NULL,
                                    value double precision ,
                                    delta bigint
);
	`

	_, err = pool.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &PostgresStorage{Conn: pool}, nil
}
