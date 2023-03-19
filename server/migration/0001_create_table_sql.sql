-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    "id" uuid NOT NULL,
    "address" text COLLATE "pg_catalog"."default",
    "command" text COLLATE "pg_catalog"."default",
    "hostname" text COLLATE "pg_catalog"."default",
    "worker_uuid" uuid NOT NULL,
    "timestamp" date
);


-- +goose StatementEnd
