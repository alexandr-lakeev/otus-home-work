package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	useCase app.UseCase
}

func NewHandler(useCase app.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) CreateEvent(ctx context.Context) http.HandlerFunc {
	type request struct {
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description"`
		Date        time.Time `json:"date" validate:"required"`
		Duration    int       `json:"duration"`
	}

	type response struct {
		ID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := new(request)

		userID, err := h.extractUserId(r)
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			makeResponseError(w, http.StatusInternalServerError, err)
			return
		}

		if err := validator.New().Struct(request); err != nil {
			makeResponseError(w, http.StatusBadRequest, err)
			return
		}

		id := uuid.New()

		err = h.useCase.CreateEvent(ctx, &app.CreateEventCommand{
			ID:          id,
			UserID:      *userID,
			Title:       request.Title,
			Description: request.Description,
			Date:        request.Date,
			Duration:    time.Minute * time.Duration(request.Duration),
		})
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		makeResponse(w, http.StatusOK, response{
			ID: id.String(),
		})
	}
}

func (h *Handler) GetEvent(ctx context.Context) http.HandlerFunc {
	type responseEvent struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Date        time.Time `json:"date"`
		Duration    int       `json:"duration"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := h.extractUserId(r)
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		id, err := h.extractId(r)
		if err != nil {
			makeResponseError(w, http.StatusBadRequest, err)
			return
		}

		event, err := h.useCase.GetEvent(ctx, &app.GetEventQuery{
			ID:     *id,
			UserID: *userID,
		})
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		makeResponse(w, http.StatusOK, &responseEvent{
			ID:          event.ID.String(),
			Title:       event.Title,
			Description: event.Description,
			Date:        event.Date,
			Duration:    int(event.Duration.Minutes()),
		})
	}
}

func (h *Handler) UpdateEvent(ctx context.Context) http.HandlerFunc {
	type request struct {
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description"`
		Date        time.Time `json:"date" validate:"required"`
		Duration    int       `json:"duration"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := new(request)

		userID, err := h.extractUserId(r)
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		id, err := h.extractId(r)
		if err != nil {
			makeResponseError(w, http.StatusBadRequest, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			makeResponseError(w, http.StatusInternalServerError, err)
			return
		}

		if err := validator.New().Struct(request); err != nil {
			makeResponseError(w, http.StatusBadRequest, err)
			return
		}

		err = h.useCase.UpdateEvent(ctx, &app.UpdateEventCommand{
			ID:          *id,
			UserID:      *userID,
			Title:       request.Title,
			Description: request.Description,
			Date:        request.Date,
			Duration:    time.Minute * time.Duration(request.Duration),
		})
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		makeResponse(w, http.StatusNoContent, nil)
	}
}

func (h *Handler) extractUserId(r *http.Request) (*models.ID, error) {
	headerUserID := r.Header.Get("x-user-id")
	if headerUserID == "" {
		return nil, domain.ErrUnauthorized
	}

	userID, err := uuid.Parse(headerUserID)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (h Handler) extractId(r *http.Request) (*models.ID, error) {
	requestId, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, errors.New("event id not specified")
	}

	id, err := uuid.Parse(requestId)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (h *Handler) resolveErrorCode(err error) int {
	errNotFound := domain.ErrEventNotFound
	errUnauthorized := domain.ErrUnauthorized
	errPremissionDenied := domain.ErrPremissionDenied
	errDateBusy := domain.ErrDateBusy

	if errors.Is(err, errNotFound) {
		return http.StatusNotFound
	} else if errors.Is(err, errUnauthorized) {
		return http.StatusUnauthorized
	} else if errors.Is(err, errPremissionDenied) {
		return http.StatusForbidden
	} else if errors.Is(err, errDateBusy) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
