package deliveryhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	app "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain/models"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	useCase app.UseCase
}

type responseEvent struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Duration    int       `json:"duration"`
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
		Duration    int       `json:"duration" validate:"required"`
	}

	type response struct {
		ID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := new(request)

		userID, err := h.extractUserID(r)
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
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := h.extractUserID(r)
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		id, err := h.extractID(r)
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

func (h *Handler) GetList(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := h.extractUserID(r)
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		from, err := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
		if err != nil {
			makeResponseError(w, http.StatusBadRequest, err)
			return
		}

		to, err := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
		if err != nil {
			makeResponseError(w, http.StatusBadRequest, err)
			return
		}

		events, err := h.useCase.GetList(ctx, &app.GetListQuery{
			UserID: *userID,
			From:   from,
			To:     to,
		})
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		response := make([]responseEvent, len(events))

		for k, event := range events {
			response[k] = responseEvent{
				ID:          event.ID.String(),
				Title:       event.Title,
				Description: event.Description,
				Date:        event.Date,
				Duration:    int(event.Duration.Minutes()),
			}
		}

		makeResponse(w, http.StatusOK, response)
	}
}

func (h *Handler) UpdateEvent(ctx context.Context) http.HandlerFunc {
	type request struct {
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description"`
		Date        time.Time `json:"date" validate:"required"`
		Duration    int       `json:"duration" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := new(request)

		userID, err := h.extractUserID(r)
		if err != nil {
			makeResponseError(w, h.resolveErrorCode(err), err)
			return
		}

		id, err := h.extractID(r)
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

func (h *Handler) extractUserID(r *http.Request) (*models.ID, error) {
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

func (h Handler) extractID(r *http.Request) (*models.ID, error) {
	requestID, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, errors.New("event id not specified")
	}

	id, err := uuid.Parse(requestID)
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

	switch {
	case errors.Is(err, errNotFound):
		return http.StatusNotFound
	case errors.Is(err, errUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, errPremissionDenied):
		return http.StatusForbidden
	case errors.Is(err, errDateBusy):
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
