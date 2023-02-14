package storage

import (
	"time"

	"github.com/gofrs/uuid"
)

// struct information from workers
type Information struct {
	ID          uuid.UUID
	Hostname    string
	Command     string
	Description string
	Time        time.Time
	Worker      string
}
