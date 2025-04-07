package migrations

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"metric-collector/internal/server/config"
)

type MigrateRunner interface {
	Up() error
}

var newMigrator = func(sourceURL, databaseURL string) (MigrateRunner, error) {
	return migrate.New(sourceURL, databaseURL)
}

func StartMigrations() error {
	dns := config.GetConfig().DatabaseDNS
	log.Info("Starting database migrations , dns: ", dns)

	if dns == "" {
		return errors.New("no database dns configured")
	}

	m, err := newMigrator("file://./migrations", dns)
	if err != nil {
		log.Error("Failed to create migrator: ", err)
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error("Migration failed: ", err)
		return err
	} else {
		log.Info("Migrations applied successfully or no changes needed.")
	}
	return nil
}
