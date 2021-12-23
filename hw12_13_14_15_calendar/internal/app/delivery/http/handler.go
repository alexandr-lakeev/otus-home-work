package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/domain"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

		headerUserID := r.Header.Get("x-user-id")
		if headerUserID == "" {
			makeResponse(w, http.StatusUnauthorized, nil)
			return
		}

		userID, err := uuid.Parse(headerUserID)
		if err != nil {
			makeResponseError(w, http.StatusInternalServerError, err)
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
			UserID:      userID,
			Title:       request.Title,
			Description: request.Description,
			Date:        request.Date,
			Duration:    time.Minute * time.Duration(request.Duration),
		})
		if err != nil {
			errDateBusy := domain.ErrDateBusy

			if errors.Is(err, errDateBusy) {
				makeResponseError(w, http.StatusBadRequest, err)
			} else {
				makeResponseError(w, http.StatusInternalServerError, err)
			}

			return
		}

		makeResponse(w, http.StatusOK, response{
			ID: id.String(),
		})
	}
}
