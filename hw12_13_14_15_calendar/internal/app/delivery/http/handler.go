package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
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

// TODO replace log with logger
// TODO return bad request errors
func (h *Handler) CreateEvent(ctx context.Context) http.HandlerFunc {
	type createEventRequest struct {
		ID          string    `json:"id" validate:"required"`
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description"`
		Date        time.Time `json:"date" validate:"required"`
		Duration    int       `json:"duration"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := new(createEventRequest)

		headerUserId := r.Header.Get("x-user-id")

		if headerUserId == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := validator.New().Struct(request); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(request.ID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userId, err := uuid.Parse(headerUserId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = h.useCase.CreateEvent(ctx, &app.CreateEventCommand{
			ID:          id,
			UserID:      userId,
			Title:       request.Title,
			Description: request.Description,
			Date:        request.Date,
			Duration:    time.Minute * time.Duration(request.Duration),
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
