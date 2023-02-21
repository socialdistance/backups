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
		t.Fatal("Failed to connect to DB server", err)
	}

	t.Run("test SQL", func(t *testing.T) {
		tx, err := storage.conn.BeginTx(ctx, pgx4.TxOptions{
			IsoLevel:       pgx4.Serializable,
			AccessMode:     pgx4.ReadWrite,
			DeferrableMode: pgx4.NotDeferrable,
		})
		if err != nil {
			t.Fatal("Failed to connect to DB server", err)
		}

		worker_UUID := uuid.New()

		timestamp, err := time.Parse("2006-01-02 15:04:05", "2022-03-14 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}

		event := internalstorage.NewEvent("hostname_test", "command_test", "description_test", timestamp, worker_UUID)

		err = storage.CreateEvent(*event)
		if err != nil {
			t.FailNow()
			return
		}

		saved, err := storage.FindAllEvents()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, saved, 1)
		require.Equal(t, event.Command, saved[0].Command)

		saveOne, err := storage.Find(event.Worker_UUID)
		if err != nil {
			t.FailNow()
			return
		}
		require.Equal(t, event.Worker_UUID, saveOne.Worker_UUID)

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
