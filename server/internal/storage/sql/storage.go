package sql

import (
	"context"
	"server/internal/storage"

	"github.com/gofrs/uuid"
	pgx4 "github.com/jackc/pgx/v4"
)

type Storage struct {
	ctx  context.Context
	conn *pgx4.Conn
	url  string
}

func NewConnect(ctx context.Context, url string) *Storage {
	return &Storage{
		ctx: ctx,
		url: url,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := pgx4.Connect(ctx, s.url)
	if err != nil {
		return err
	}

	s.conn = conn

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) CreateEvent(e storage.Event) error {
	return nil
}

func (s *Storage) DeleteEvent(id uuid.UUID) error {
	return nil
}

func (s *Storage) Find(id uuid.UUID) (*storage.Event, error) {
	return nil, nil
}

func (s *Storage) FindAllEvents() ([]storage.Event, error) {
	return nil, nil
}
