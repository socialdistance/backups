package storage

import (
	internalmemory "server/internal/storage"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	storage := NewMemory()

	t.Run("storage test", func(t *testing.T) {
		worker_UUID := uuid.New()

		timestamp, err := time.Parse("2006-01-02 15:04:05", "2022-03-14 12:00:00")
		if err != nil {
			t.FailNow()
			return
		}

		event := internalmemory.NewEvent(
			"hostname_test", "command_test", "description_test", timestamp, worker_UUID)

		_, err = storage.CreateEvent(*event)
		if err != nil {
			t.FailNow()
			return
		}

		findEvent, err := storage.Find(event.ID)
		if err != nil {
			t.FailNow()
			return
		}
		require.Equal(t, findEvent.Description, event.Description)

		events, err := storage.FindAllEvents()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, events, 1)
		require.Equal(t, *event, events[0])

		err = storage.DeleteEvent(event.ID)
		if err != nil {
			t.FailNow()
			return
		}

		events, err = storage.FindAllEvents()
		if err != nil {
			t.FailNow()
			return
		}
		require.Len(t, events, 0)
	})
}
