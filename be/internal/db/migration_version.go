package db

import (
	"context"
)

func GetLatestMigrationVersion() (*int, error) {
	var version int

	qStr := `
		SELECT version
		FROM schema_migrations`

	pool, err := Init()
	if err != nil {
		return nil, err
	}

	if err := pool.QueryRow(context.Background(), qStr).Scan(&version); err != nil {
		return nil, err
	}

	version--

	return &version, nil
}
