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
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT subscription_service_name_not_empty CHECK (BTRIM(service_name) <> ''),
    CONSTRAINT subscription_price_positive CHECK (price > 0),
    CONSTRAINT subscription_valid_period CHECK (
        end_date IS NULL
            OR end_date >= start_date
        )
);