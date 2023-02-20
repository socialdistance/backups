package app

import (
	"context"
	"log"
	"testing"
	"time"

	internallogger "server/internal/logger"
	internalcache "server/internal/storage/cache"
	internalstorage "server/internal/storage/memory"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	worker_uuid := uuid.New()

	logg, err := internallogger.NewLogger()
	if err != nil {
		log.Fatalf("Failed logger %s", err)
	}

	memmoryStorage := internalstorage.NewMemory()

	ctx := context.Background()

	cache := internalcache.NewCache(5*time.Minute, 10*time.Minute)

	testApp := NewApp(logg, memmoryStorage, cache)

	t.Run("CommandHandlerApp test", func(t *testing.T) {
		task, err := testApp.CommandHandlerApp(ctx, worker_uuid)
		if err != nil {
			t.FailNow()
			return
		}

		require.NotNil(t, task)
	})

}
