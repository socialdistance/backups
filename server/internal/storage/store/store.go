package store

import (
	"context"
	"log"
	internalapp "server/internal/app"
	internalconfig "server/internal/config"
	internalmemory "server/internal/storage/memory"
	internalsql "server/internal/storage/sql"
)

func CreateStorage(ctx context.Context, config internalconfig.Config) internalapp.Storage {
	var store internalapp.Storage

	switch config.Storage.Type {
	case "in-memory":
		store = internalmemory.NewMemory()

	case "sql":
		sqlStore := internalsql.NewConnect(ctx, config.Storage.Dsn)
		if err := sqlStore.Connect(ctx); err != nil {
			log.Fatalf("Unable to connect database %s", err)
		}
		store = sqlStore

	default:
		log.Fatalf("Dont know type storage: %s", config.Storage.Type)
	}

	return store
}
