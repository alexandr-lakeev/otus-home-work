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

type errorResponse struct {
	errors map[string]string
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

		headerUserId := r.Header.Get("x-user-id")
		if headerUserId == "" {
			makeResponse(w, r, http.StatusUnauthorized, nil)
			return
		}

		userId, err := uuid.Parse(headerUserId)
		if err != nil {
			makeResponseError(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			makeResponseError(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := validator.New().Struct(request); err != nil {
			makeResponseError(w, r, http.StatusBadRequest, err.(validator.ValidationErrors))
			return
		}

		id := uuid.New()

		err = h.useCase.CreateEvent(ctx, &app.CreateEventCommand{
			ID:          id,
			UserID:      userId,
			Title:       request.Title,
			Description: request.Description,
			Date:        request.Date,
			Duration:    time.Minute * time.Duration(request.Duration),
		})
		if err != nil {
			errDateBusy := domain.ErrDateBusy

			if errors.Is(err, errDateBusy) {
				makeResponseError(w, r, http.StatusBadRequest, err)
			} else {
				makeResponseError(w, r, http.StatusInternalServerError, err)
			}

			return
		}

		makeResponse(w, r, http.StatusOK, response{
			ID: id.String(),
		})
	}
}
