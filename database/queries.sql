-- name: CreateSubscription :one
INSERT INTO practice.subscription (service_name,
                                   price,
                                   start_date,
                                   end_date,
                                   user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetSubscriptionByID :one
SELECT *
FROM practice.subscription
WHERE id = $1;

-- name: ListSubscriptions :many
SELECT *
FROM practice.subscription
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: UpdateSubscription :one
UPDATE practice.subscription
SET service_name = $2,
    price        = $3,
    start_date   = $4,
    end_date     = $5,
    user_id      = $6,
    updated_at   = now()
WHERE id = $1
RETURNING
    *;

-- name: DeleteSubscription :exec
DELETE
FROM practice.subscription
WHERE id = $1;

-- name: GetSubscriptionsTotal :one
SELECT COALESCE(SUM(s.price), 0)::bigint AS total_price
FROM practice.subscription AS s
         JOIN generate_series(
        @period_start::date,
        @period_end::date,
        interval '1 month'
              ) AS months(month_start)
              ON s.start_date <= months.month_start
                  AND (
                     s.end_date IS NULL
                         OR s.end_date >= months.month_start
                     )
WHERE s.service_name = $1
  AND s.user_id = $2;