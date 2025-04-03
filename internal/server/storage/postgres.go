package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/config"
)

func HealthCheck() error {
	conn, err := pgx.Connect(context.Background(), config.GetConfig().DatabaseDNS)
	if err != nil {
		log.Error("Failed connection to database", err)
		return err
	}
	defer conn.Close(context.Background())
	return nil
}
