-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS slot (
    slot_id SERIAL PRIMARY KEY,
    slot_description text NOT NULL,
    total_display integer DEFAULT 1
);

CREATE TABLE IF NOT EXISTS banner (
    banner_id SERIAL PRIMARY KEY,
    banner_description text NOT NULL,
    total_display integer DEFAULT  1
);

CREATE TABLE IF NOT EXISTS banner_to_slot (
    banner_to_slot_id SERIAL PRIMARY KEY,
    banner_id integer REFERENCES banner (banner_id),
    slot_id integer REFERENCES slot (slot_id),
    UNIQUE (slot_id, banner_id)
);

CREATE TABLE IF NOT EXISTS social_group (
    social_group_id SERIAL PRIMARY KEY,
    social_description text NOT NULL
);

CREATE TABLE IF NOT EXISTS statistics (
    statistics_id SERIAL PRIMARY KEY,
    banner_id integer REFERENCES banner (banner_id) NOT NULL,
    slot_id integer REFERENCES slot (slot_id) NOT NULL,
    social_group_id integer REFERENCES social_group (social_group_id) NOT NULL,
    display integer DEFAULT 1,
    click integer DEFAULT 0,
    UNIQUE (slot_id, banner_id, social_group_id)
);

-- +goose StatementEnd
