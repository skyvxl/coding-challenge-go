-- +goose Up
ALTER TABLE practice.subscription
    ADD CONSTRAINT subscription_service_name_not_empty
        CHECK (BTRIM(service_name) <> '');

ALTER TABLE practice.subscription
    ADD CONSTRAINT subscription_price_positive
        CHECK (price > 0);

ALTER TABLE practice.subscription
    ADD CONSTRAINT subscription_valid_period
        CHECK (
            end_date IS NULL
                OR end_date >= start_date
            );

-- +goose Down
ALTER TABLE practice.subscription
    DROP CONSTRAINT subscription_service_name_not_empty;

ALTER TABLE practice.subscription
    DROP CONSTRAINT subscription_price_positive;

ALTER TABLE practice.subscription
    DROP CONSTRAINT subscription_valid_period;
