package storage

import (
	"fmt"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/storage"
	memorystorage "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/storage/memory"
	sqlstorage "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/storage/sql"
)

func ResolveStorage(cfg config.StorageConf) (storage.Storage, error) {
	if cfg.Type == config.STORAGE_MEMORY {
		return memorystorage.New(), nil
	}

	switch cfg.Type {
	case config.STORAGE_MEMORY:
		return memorystorage.New(), nil
	case config.STORAGE_SQL:
		return sqlstorage.New(cfg.DSN), nil
	default:
		return nil, fmt.Errorf("wrong storage type: %s", cfg.Type)
	}
}
