package handler

type SubscriptionResponse struct {
	ID          int32   `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int32   `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	UserID      string  `json:"user_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateSubscriptionParams struct {
	ServiceName string  `json:"service_name"`
	Price       int32   `json:"price"`
	StartDate   string  `json:"start_date"`
	UserID      string  `json:"user_id"`
	EndDate     *string `json:"end_date,omitempty"`
}

type UpdateSubscriptionParams struct {
	ServiceName string  `json:"service_name"`
	Price       int32   `json:"price"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	UserID      string  `json:"user_id"`
}

type TotalSubscriptionsResponse struct {
	TotalPrice int64 `json:"total_price"`
}
