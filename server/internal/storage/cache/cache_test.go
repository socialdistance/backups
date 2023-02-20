package storage

import (
	internalmemory "server/internal/storage"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cache := NewCache(10*time.Minute, 1*time.Hour)

	worker_UUID := uuid.New()

	timestamp, err := time.Parse("2006-01-02 15:04:05", "2022-03-14 12:00:00")
	if err != nil {
		t.FailNow()
		return
	}

	event := internalmemory.NewEvent(
		"hostname_test", "command_test", "description_test", timestamp, worker_UUID)

	t.Run("cache test", func(t *testing.T) {
		// test empty key
		value, found := cache.Get(uuid.New())
		require.Nil(t, value)
		require.False(t, found)

		cache.Set(event.ID, *event, 1*time.Minute)

		value, found = cache.Get(event.ID)
		require.Equal(t, value.Command, event.Command)
		require.True(t, found)

		error := cache.Delete(event.ID)
		require.Nil(t, error)

		// trying get value from delete cache
		value, found = cache.Get(event.ID)
		require.Nil(t, value)
		require.False(t, found)

		// trying delete not existing key
		error = cache.Delete(uuid.New())
		require.EqualError(t, error, "Key not found")
	})
}
