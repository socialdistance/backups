package sql

import "context"

import (
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
