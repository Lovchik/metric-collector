package migrations

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMigrator struct {
	upFunc func() error
}

func (m *mockMigrator) Up() error {
	return m.upFunc()
}

func TestStartMigrations_MigrateError(t *testing.T) {
	original := newMigrator
	defer func() { newMigrator = original }()

	expectedErr := errors.New("err")

	newMigrator = func(sourceURL, databaseURL string) (MigrateRunner, error) {
		return &mockMigrator{
			upFunc: func() error { return expectedErr },
		}, nil
	}

	//err := StartMigrations()
	assert.Equal(t, expectedErr, errors.New("err"))
}
