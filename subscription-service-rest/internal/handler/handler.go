package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"learn/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	queries db.Querier
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandler(queries db.Querier) *Handler {
	return &Handler{queries: queries}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(
		"GET /swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)
	mux.HandleFunc("GET /api/v1/subscriptions", h.listSubscriptions)
	mux.HandleFunc("POST /api/v1/subscriptions", h.createSubscription)
	mux.HandleFunc("GET /api/v1/subscriptions/{id}", h.getSubscriptionByID)
	mux.HandleFunc("PUT /api/v1/subscriptions/{id}", h.updateSubscriptionByID)
	mux.HandleFunc("DELETE /api/v1/subscriptions/{id}", h.deleteSubscriptionByID)
	mux.HandleFunc("GET /api/v1/subscriptions/total", h.getSubscriptionsTotal)
	return mux
}

// createSubscription creates a new subscription.
//
//	@Summary Create a new subscription
//	@Description Create a new subscription.
//	@Accept json
//	@Produce json
//	@Param subscription body CreateSubscriptionParams true "Subscription to create"
//	@Success 200 {object} SubscriptionResponse "Created subscription"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 500 {object} ErrorResponse "Internal server error"
//	@Tags subscriptions
//	@Router /subscriptions [post]
func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	var req CreateSubscriptionParams
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	req.ServiceName = strings.TrimSpace(req.ServiceName)
	if req.ServiceName == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("service_name is required"))
		return
	}
	if req.Price <= 0 {
		writeError(w, http.StatusBadRequest, fmt.Errorf("price must be greater than 0"))
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid user_id: %w", err))
		return
	}
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid start_date: %w", err))
		return
	}
	var endDate *time.Time

	if req.EndDate != nil {
		parsedEndDate, err := parseDate(*req.EndDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid end_date: %w", err))
			return
		}
		if parsedEndDate.Before(startDate) {
			writeError(w, http.StatusBadRequest, fmt.Errorf("end_date must be after start_date"))
			return
		}
		endDate = &parsedEndDate
	}

	sub, err := h.queries.CreateSubscription(
		r.Context(),
		db.CreateSubscriptionParams{
			ServiceName: req.ServiceName,
			Price:       req.Price,
			StartDate:   startDate,
			EndDate:     endDate,
			UserID:      userID,
		})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, sub)
}

// getSubscriptionByID returns a subscription by its ID.
//
//	@Summary Get a subscription by ID
//	@Description Get a subscription by its ID.
//	@Accept json
//	@Produce json
//	@Param id path int true "Subscription ID"
//	@Success 200 {object} SubscriptionResponse "Subscription"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 404 {object} ErrorResponse "Subscription not found"
//	@Tags subscriptions
//	@Router /subscriptions/{id} [get]
func (h *Handler) getSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("id must be a positive integer"))
		return
	}
	sub, err := h.queries.GetSubscriptionByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("subscription not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	writeJSON(w, http.StatusOK, sub)
}

// updateSubscriptionByID updates a subscription by its ID.
//
//	@Summary Update a subscription by ID
//	@Description Update a subscription by its ID.
//	@Accept json
//	@Produce json
//	@Param id path int true "Subscription ID"
//	@Param subscription body UpdateSubscriptionParams true "Subscription to update"
//	@Success 200 {object} SubscriptionResponse "Updated subscription"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 404 {object} ErrorResponse "Subscription not found"
//	@Tags subscriptions
//	@Router /subscriptions/{id} [put]
func (h *Handler) updateSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("id must be a positive integer"))
		return
	}
	var req UpdateSubscriptionParams
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	req.ServiceName = strings.TrimSpace(req.ServiceName)
	if req.ServiceName == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("service_name is required"))
		return
	}
	if req.Price <= 0 {
		writeError(w, http.StatusBadRequest, fmt.Errorf("price must be greater than 0"))
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid user_id: %w", err))
		return
	}
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid start_date: %w", err))
		return
	}
	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err := parseDate(*req.EndDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid end_date: %w", err))
			return
		}
		if parsedEndDate.Before(startDate) {
			writeError(w, http.StatusBadRequest, fmt.Errorf("end_date must be after start_date"))
			return
		}
		endDate = &parsedEndDate
	}
	sub, err := h.queries.UpdateSubscription(
		r.Context(),
		db.UpdateSubscriptionParams{
			ID:          int32(id),
			ServiceName: req.ServiceName,
			Price:       req.Price,
			StartDate:   startDate,
			EndDate:     endDate,
			UserID:      userID,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("subscription not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, sub)
}

// deleteSubscriptionByID deletes a subscription by its ID.
//
//	@Summary Delete a subscription by ID
//	@Description Delete a subscription by its ID.
//	@Accept json
//	@Produce json
//	@Param id path int true "Subscription ID"
//	@Success 200 {object} nil "Subscription deleted"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 404 {object} ErrorResponse "Subscription not found"
//	@Tags subscriptions
//	@Router /subscriptions/{id} [delete]
func (h *Handler) deleteSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("id must be a positive integer"))
		return
	}
	_, err = h.queries.GetSubscriptionByID(r.Context(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("subscription not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
		return
	}
	err = h.queries.DeleteSubscription(r.Context(), int32(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, nil)
}

// listSubscriptions returns a list of subscriptions.
//
//	@Summary List subscriptions
//	@Description List subscriptions.
//	@Accept json
//	@Produce json
//	@Param limit query int false "Limit the number of subscriptions returned"
//	@Param offset query int false "Offset the subscriptions returned"
//	@Success 200 {array} SubscriptionResponse "List of subscriptions"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 500 {object} ErrorResponse "Internal server error"
//	@Tags subscriptions
//	@Router /subscriptions [get]
func (h *Handler) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}
	offset := r.URL.Query().Get("offset")
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}
	subs, err := h.queries.ListSubscriptions(
		r.Context(),
		db.ListSubscriptionsParams{
			Limit:  limitInt,
			Offset: offsetInt,
		})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, subs)
}

// getSubscriptionsTotal returns the total price of subscriptions within a specified period.
//
//	@Summary Get total price of subscriptions
//	@Description Get total price of subscriptions within a specified period.
//	@Accept json
//	@Produce json
//	@Param start_date query string true "Start date of the period (format: 01-2006)"
//	@Param end_date query string true "End date of the period (format: 01-2006)"
//	@Param user_id query string true "User ID"
//	@Param service_name query string true "Service name"
//	@Success 200 {object} TotalSubscriptionsResponse "Total price of subscriptions"
//	@Failure 400 {object} ErrorResponse "Invalid request"
//	@Failure 500 {object} ErrorResponse "Internal server error"
//	@Tags subscriptions
//	@Router /subscriptions/total [get]
func (h *Handler) getSubscriptionsTotal(w http.ResponseWriter, r *http.Request) {
	startDate, err := parseDate(r.URL.Query().Get("start_date"))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid start_date: %w", err))
		return
	}
	endDate, err := parseDate(r.URL.Query().Get("end_date"))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid end_date: %w", err))
		return
	}
	if endDate.Before(startDate) {
		writeError(
			w,
			http.StatusBadRequest,
			errors.New("end_date must not be before start_date"),
		)
		return
	}
	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid user_id: %w", err))
		return
	}
	serviceName := r.URL.Query().Get("service_name")
	serviceName = strings.TrimSpace(serviceName)
	if serviceName == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("service_name is required"))
		return
	}
	total, err := h.queries.GetSubscriptionsTotal(
		r.Context(),
		db.GetSubscriptionsTotalParams{
			ServiceName: serviceName,
			UserID:      userID,
			PeriodStart: startDate,
			PeriodEnd:   endDate,
		},
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, TotalSubscriptionsResponse{TotalPrice: total})
}
