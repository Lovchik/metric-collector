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
	log.Println("Starting database migrations, dns: ", dns)

	if dns == "" {
		return errors.New("no database dns configured")
	}
	sourceURL := "file://./migrations"

	m, err := migrate.New(sourceURL, dns)
	if err != nil {
		log.Println("Failed to create migrator: ", err)
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Println("Migration failed: ", err)
		return err
	} else {
		log.Println("Migrations applied successfully or no changes needed.")
	}

	return nil
}
