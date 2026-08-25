-- +goose Up
CREATE SCHEMA IF NOT EXISTS practice;

CREATE TABLE IF NOT EXISTS practice.subscription
(
    id           SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price        INTEGER      NOT NULL,
    start_date   DATE         NOT NULL,
    end_date     DATE,
    user_id      UUID         NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS practice.subscription;

DROP SCHEMA IF EXISTS practice;
