package sql

import (
	"context"
	"testing"
	"time"

	internalstorage "server/internal/storage"

	pgx4 "github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestStorageSql(t *testing.T) {
	ctx := context.Background()
	storage := NewConnect(ctx, "postgres://postgres:postgres@localhost:54321/backups?sslmode=disable")
	if err := storage.Connect(ctx); err != nil {
		t.Fatal("Failed to connect to DB client", err)
	}

	t.Run("test SQL", func(t *testing.T) {
		tx, err := storage.conn.BeginTx(ctx, pgx4.TxOptions{
			IsoLevel:       pgx4.Serializable,
			AccessMode:     pgx4.ReadWrite,
			DeferrableMode: pgx4.NotDeferrable,
		})
		if err != nil {
			t.Fatal("Failed to connect to DB client", err)
		}

		workerUuid := uuid.New()

		timestamp := time.Now()

		event := internalstorage.NewEvent("hostname_test", "command_test", "description_test", timestamp, workerUuid)

		err = storage.CreateEvent(*event)
		if err != nil {
			t.FailNow()
			return
		}
		event.Command = "test"

		err = storage.UpdateCommand(*event)
		if err != nil {
			t.FailNow()
			return
		}
		require.Nil(t, err)

		saved, err := storage.FindAllEvents()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 1)
		require.Equal(t, event.Command, saved[0].Command)

		saveOne, err := storage.Find(event.WorkerUuid)
		if err != nil {
			t.FailNow()
			return
		}
		require.Equal(t, event.WorkerUuid, saveOne.WorkerUuid)

		err = storage.DeleteEvent(event.ID)
		if err != nil {
			t.FailNow()
			return
		}
		require.Nil(t, err)

		err = tx.Rollback(ctx)
		if err != nil {
			t.Fatal("Failed to rollback tx", err)
		}
	})
}
