package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"learn/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type fakeQuerier struct {
	create func(context.Context, db.CreateSubscriptionParams) (db.PracticeSubscription, error)
	get    func(context.Context, int32) (db.PracticeSubscription, error)
	list   func(context.Context, db.ListSubscriptionsParams) ([]db.PracticeSubscription, error)
	update func(context.Context, db.UpdateSubscriptionParams) (db.PracticeSubscription, error)
	delete func(context.Context, int32) error
	total  func(context.Context, db.GetSubscriptionsTotalParams) (int64, error)
}

func (f *fakeQuerier) CreateSubscription(
	ctx context.Context,
	arg db.CreateSubscriptionParams,
) (db.PracticeSubscription, error) {
	return f.create(ctx, arg)
}

func (f *fakeQuerier) GetSubscriptionByID(
	ctx context.Context,
	id int32,
) (db.PracticeSubscription, error) {
	return f.get(ctx, id)
}

func (f *fakeQuerier) ListSubscriptions(
	ctx context.Context,
	arg db.ListSubscriptionsParams,
) ([]db.PracticeSubscription, error) {
	return f.list(ctx, arg)
}

func (f *fakeQuerier) UpdateSubscription(
	ctx context.Context,
	arg db.UpdateSubscriptionParams,
) (db.PracticeSubscription, error) {
	return f.update(ctx, arg)
}

func (f *fakeQuerier) DeleteSubscription(ctx context.Context, id int32) error {
	return f.delete(ctx, id)
}

func (f *fakeQuerier) GetSubscriptionsTotal(
	ctx context.Context,
	arg db.GetSubscriptionsTotalParams,
) (int64, error) {
	return f.total(ctx, arg)
}

func TestHandler_CreateSubscription(t *testing.T) {
	want := subscriptionFixture()
	store := &fakeQuerier{
		create: func(_ context.Context, arg db.CreateSubscriptionParams) (db.PracticeSubscription, error) {
			if arg.ServiceName != want.ServiceName {
				t.Errorf("service name: expected %q, got %q", want.ServiceName, arg.ServiceName)
			}
			if arg.Price != want.Price {
				t.Errorf("price: expected %d, got %d", want.Price, arg.Price)
			}
			if arg.UserID != want.UserID {
				t.Errorf("user ID: expected %s, got %s", want.UserID, arg.UserID)
			}
			if !arg.StartDate.Equal(want.StartDate) {
				t.Errorf("start date: expected %s, got %s", want.StartDate, arg.StartDate)
			}
			return want, nil
		},
	}

	body := `{
		"service_name":"Yandex Plus",
		"price":400,
		"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba",
		"start_date":"07-2025"
	}`
	rec := serve(t, NewHandler(store), http.MethodPost, "/api/v1/subscriptions", body)

	assertStatus(t, rec, http.StatusOK)
	assertSubscriptionID(t, rec, want.ID)
}

func TestHandler_GetSubscriptionByID(t *testing.T) {
	want := subscriptionFixture()
	store := &fakeQuerier{
		get: func(_ context.Context, id int32) (db.PracticeSubscription, error) {
			if id != want.ID {
				t.Errorf("expected ID %d, got %d", want.ID, id)
			}
			return want, nil
		},
	}

	rec := serve(t, NewHandler(store), http.MethodGet, "/api/v1/subscriptions/1", "")

	assertStatus(t, rec, http.StatusOK)
	assertSubscriptionID(t, rec, want.ID)
}

func TestHandler_GetSubscriptionByID_NotFound(t *testing.T) {
	store := &fakeQuerier{
		get: func(context.Context, int32) (db.PracticeSubscription, error) {
			return db.PracticeSubscription{}, pgx.ErrNoRows
		},
	}

	rec := serve(t, NewHandler(store), http.MethodGet, "/api/v1/subscriptions/42", "")

	assertStatus(t, rec, http.StatusNotFound)
	assertError(t, rec, "subscription not found")
}

func TestHandler_GetSubscriptionByID_RejectsInvalidID(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "zero", id: "0"},
		{name: "negative", id: "-1"},
		{name: "not a number", id: "abc"},
		{name: "greater than int32", id: "2147483648"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(
				t,
				NewHandler(nil),
				http.MethodGet,
				"/api/v1/subscriptions/"+tt.id,
				"",
			)

			assertStatus(t, rec, http.StatusBadRequest)
			assertError(t, rec, "id must be a positive integer")
		})
	}
}

func TestHandler_ListSubscriptions(t *testing.T) {
	want := subscriptionFixture()
	store := &fakeQuerier{
		list: func(_ context.Context, arg db.ListSubscriptionsParams) ([]db.PracticeSubscription, error) {
			if arg.Limit != 2 || arg.Offset != 3 {
				t.Errorf("expected limit=2 offset=3, got limit=%d offset=%d", arg.Limit, arg.Offset)
			}
			return []db.PracticeSubscription{want}, nil
		},
	}

	rec := serve(
		t,
		NewHandler(store),
		http.MethodGet,
		"/api/v1/subscriptions?limit=2&offset=3",
		"",
	)

	assertStatus(t, rec, http.StatusOK)
	var response []struct {
		ID int32 `json:"id"`
	}
	decodeResponse(t, rec, &response)
	if len(response) != 1 || response[0].ID != want.ID {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestHandler_UpdateSubscription(t *testing.T) {
	want := subscriptionFixture()
	store := &fakeQuerier{
		get: func(context.Context, int32) (db.PracticeSubscription, error) {
			return want, nil
		},
		update: func(_ context.Context, arg db.UpdateSubscriptionParams) (db.PracticeSubscription, error) {
			if arg.ID != want.ID {
				t.Errorf("expected path ID %d, got %d", want.ID, arg.ID)
			}
			if arg.ServiceName != "Updated service" || arg.Price != 500 {
				t.Errorf("unexpected update params: %+v", arg)
			}
			return want, nil
		},
	}

	body := `{
		"service_name":"Updated service",
		"price":500,
		"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba",
		"start_date":"07-2025"
	}`
	rec := serve(t, NewHandler(store), http.MethodPut, "/api/v1/subscriptions/1", body)

	assertStatus(t, rec, http.StatusOK)
}

func TestHandler_DeleteSubscription(t *testing.T) {
	deleted := false
	store := &fakeQuerier{
		get: func(context.Context, int32) (db.PracticeSubscription, error) {
			return subscriptionFixture(), nil
		},
		delete: func(_ context.Context, id int32) error {
			if id != 1 {
				t.Errorf("expected ID 1, got %d", id)
			}
			deleted = true
			return nil
		},
	}

	rec := serve(t, NewHandler(store), http.MethodDelete, "/api/v1/subscriptions/1", "")

	assertStatus(t, rec, http.StatusOK)
	if !deleted {
		t.Fatal("expected DeleteSubscription to be called")
	}
}

func TestHandler_StorageErrors(t *testing.T) {
	storageErr := errors.New("database unavailable")

	t.Run("create", func(t *testing.T) {
		store := &fakeQuerier{
			create: func(context.Context, db.CreateSubscriptionParams) (db.PracticeSubscription, error) {
				return db.PracticeSubscription{}, storageErr
			},
		}
		body := `{
			"service_name":"Yandex Plus",
			"price":400,
			"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba",
			"start_date":"07-2025"
		}`
		rec := serve(t, NewHandler(store), http.MethodPost, "/api/v1/subscriptions", body)
		assertStatus(t, rec, http.StatusInternalServerError)
	})

	t.Run("get", func(t *testing.T) {
		store := &fakeQuerier{
			get: func(context.Context, int32) (db.PracticeSubscription, error) {
				return db.PracticeSubscription{}, storageErr
			},
		}
		rec := serve(t, NewHandler(store), http.MethodGet, "/api/v1/subscriptions/1", "")
		assertStatus(t, rec, http.StatusInternalServerError)
	})

	t.Run("list", func(t *testing.T) {
		store := &fakeQuerier{
			list: func(context.Context, db.ListSubscriptionsParams) ([]db.PracticeSubscription, error) {
				return nil, storageErr
			},
		}
		rec := serve(t, NewHandler(store), http.MethodGet, "/api/v1/subscriptions?limit=20", "")
		assertStatus(t, rec, http.StatusInternalServerError)
	})
}

func TestHandler_GetSubscriptionsTotal(t *testing.T) {
	userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
	store := &fakeQuerier{
		total: func(_ context.Context, arg db.GetSubscriptionsTotalParams) (int64, error) {
			if arg.ServiceName != "Yandex Plus" {
				t.Errorf("expected service name %q, got %q", "Yandex Plus", arg.ServiceName)
			}
			if arg.UserID != userID {
				t.Errorf("expected user ID %s, got %s", userID, arg.UserID)
			}
			wantStart := time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC)
			if !arg.PeriodStart.Equal(wantStart) {
				t.Errorf("expected start date %s, got %s", wantStart, arg.PeriodStart)
			}
			wantEnd := time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC)
			if !arg.PeriodEnd.Equal(wantEnd) {
				t.Errorf("expected end date %s, got %s", wantEnd, arg.PeriodEnd)
			}
			return int64(1200), nil
		},
	}

	target := "/api/v1/subscriptions/total" +
		"?start_date=07-2025" +
		"&end_date=09-2025" +
		"&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba" +
		"&service_name=Yandex+Plus"
	rec := serve(t, NewHandler(store), http.MethodGet, target, "")

	assertStatus(t, rec, http.StatusOK)
	var response TotalSubscriptionsResponse
	decodeResponse(t, rec, &response)
	if response.TotalPrice != 1200 {
		t.Fatalf("expected total 1200, got %d", response.TotalPrice)
	}
}

func TestHandler_GetSubscriptionsTotal_RejectsInvalidParams(t *testing.T) {
	tests := []struct {
		name   string
		target string
	}{
		{
			name: "missing start date",
			target: "/api/v1/subscriptions/total" +
				"?end_date=09-2025" +
				"&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba" +
				"&service_name=Yandex+Plus",
		},
		{
			name: "invalid end date",
			target: "/api/v1/subscriptions/total" +
				"?start_date=07-2025" +
				"&end_date=wrong" +
				"&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba" +
				"&service_name=Yandex+Plus",
		},
		{
			name: "invalid user ID",
			target: "/api/v1/subscriptions/total" +
				"?start_date=07-2025" +
				"&end_date=09-2025" +
				"&user_id=wrong" +
				"&service_name=Yandex+Plus",
		},
		{
			name: "missing service name",
			target: "/api/v1/subscriptions/total" +
				"?start_date=07-2025" +
				"&end_date=09-2025" +
				"&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(t, NewHandler(nil), http.MethodGet, tt.target, "")
			assertStatus(t, rec, http.StatusBadRequest)
		})
	}
}

func TestHandler_GetSubscriptionsTotal_RejectsReversedPeriod(t *testing.T) {
	store := &fakeQuerier{
		total: func(context.Context, db.GetSubscriptionsTotalParams) (int64, error) {
			return int64(0), nil
		},
	}
	target := "/api/v1/subscriptions/total" +
		"?start_date=09-2025" +
		"&end_date=07-2025" +
		"&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba" +
		"&service_name=Yandex+Plus"

	rec := serve(t, NewHandler(store), http.MethodGet, target, "")

	assertStatus(t, rec, http.StatusBadRequest)
}

func TestHandler_GetSubscriptionsTotal_StorageError(t *testing.T) {
	store := &fakeQuerier{
		total: func(context.Context, db.GetSubscriptionsTotalParams) (int64, error) {
			return 0, errors.New("database unavailable")
		},
	}
	target := "/api/v1/subscriptions/total" +
		"?start_date=07-2025" +
		"&end_date=09-2025" +
		"&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba" +
		"&service_name=Yandex+Plus"

	rec := serve(t, NewHandler(store), http.MethodGet, target, "")

	assertStatus(t, rec, http.StatusInternalServerError)
}

func subscriptionFixture() db.PracticeSubscription {
	return db.PracticeSubscription{
		ID:          1,
		ServiceName: "Yandex Plus",
		Price:       400,
		StartDate:   time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
		UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
	}
}

func serve(
	t *testing.T,
	h *Handler,
	method string,
	target string,
	body string,
) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)
	return rec
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("expected status %d, got %d; body=%s", want, rec.Code, rec.Body.String())
	}
}

func assertSubscriptionID(t *testing.T, rec *httptest.ResponseRecorder, want int32) {
	t.Helper()
	var response struct {
		ID int32 `json:"id"`
	}
	decodeResponse(t, rec, &response)
	if response.ID != want {
		t.Fatalf("expected subscription ID %d, got %d", want, response.ID)
	}
}

func assertError(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	var response ErrorResponse
	decodeResponse(t, rec, &response)
	if response.Error != want {
		t.Fatalf("expected error %q, got %q", want, response.Error)
	}
}

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(dst); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, rec.Body.String())
	}
}
