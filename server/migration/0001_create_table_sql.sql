-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    "id" uuid NOT NULL,
    "address" varchar COLLATE "pg_catalog"."default",
    "command" text COLLATE "pg_catalog"."default",
    "hostname" varchar COLLATE "pg_catalog"."default",
    "worker_uuid" uuid NOT NULL,
    "timestamp" timestamp
);
-- +goose StatementEnd
