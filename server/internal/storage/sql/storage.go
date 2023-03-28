package sql

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pgx4 "github.com/jackc/pgx/v4"
	"server/internal/storage"
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
	sql := `
		INSERT INTO events (id, address, command, hostname, worker_uuid, timestamp) VALUES 
		($1, $2, $3, $4, $5, $6)
	`

	_, err := s.conn.Exec(s.ctx, sql, e.ID.String(), e.Address, e.Command, e.Hostname, e.WorkerUuid.String(), e.Timestamp)

	return err
}

func (s *Storage) DeleteEvent(id uuid.UUID) error {
	sql := `
		DELETE FROM events where id=$1
	`

	_, err := s.conn.Exec(s.ctx, sql, id)

	return err
}

func (s *Storage) Find(workerUuid uuid.UUID) (*storage.Event, error) {
	var event storage.Event
	sql := `select id, address, command, hostname, worker_uuid, timestamp from events where worker_uuid = $1`

	err := s.conn.QueryRow(s.ctx, sql, workerUuid).Scan(
		&event.ID,
		&event.Address,
		&event.Command,
		&event.Hostname,
		&event.WorkerUuid,
		&event.Timestamp,
	)

	if err == nil {
		return &event, nil
	}

	if errors.Is(err, pgx4.ErrNoRows) {
		return nil, nil
	}

	return nil, fmt.Errorf("cant scan SQL result to struct %w", err)
}

func (s *Storage) FindAllEvents() ([]storage.Event, error) {
	var events []storage.Event

	sql := `
		SELECT id, address, command, hostname, worker_uuid, timestamp FROM events
	`

	rows, err := s.conn.Query(s.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var evt storage.Event
		if err = rows.Scan(&evt.ID,
			&evt.Address,
			&evt.Command,
			&evt.Hostname,
			&evt.WorkerUuid,
			&evt.Timestamp); err != nil {
			return nil, fmt.Errorf("cant convert result: %w", err)
		}

		events = append(events, evt)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
