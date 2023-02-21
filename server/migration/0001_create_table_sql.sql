-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    "id" uuid NOT NULL,
    "hostname" text COLLATE "pg_catalog"."default",
    "command" text COLLATE "pg_catalog"."default",
    "description" text COLLATE "pg_catalog"."default",
    "worker_uuid" uuid NOT NULL,
    "timestamp" date
);


-- +goose StatementEnd
