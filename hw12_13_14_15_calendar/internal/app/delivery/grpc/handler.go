package deliverygrpc

import (
	"context"
	"log"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/delivery/grpc/pb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	useCase app.UseCase
}

func NewHandler(useCase app.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) GetEvent(ctx context.Context, query *pb.GetEventQuery) (*pb.GetEventResponse, error) {
	id, err := uuid.Parse(query.Id)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(query.UserId)
	if err != nil {
		return nil, err
	}

	event, err := h.useCase.GetEvent(ctx, &app.GetEventQuery{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	log.Println(event)

	return &pb.GetEventResponse{
		Event: &pb.EventResponse{
			Id:          event.ID.String(),
			Title:       event.Title,
			Date:        timestamppb.New(event.Date),
			Duration:    durationpb.New(event.Duration),
			Description: event.Description,
		},
	}, nil
}

func (h *Handler) GetList(ctx context.Context, query *pb.GetListQuery) (*pb.GetListResponse, error) {
	return nil, nil
}

func (h *Handler) CreateEvent(ctx context.Context, command *pb.CreateEventCommand) (*pb.CreateEventResponse, error) {
	return nil, nil
}

func (h *Handler) UpdateEvent(ctx context.Context, command *pb.UpdateEventCommand) (*pb.UpdateEventResponse, error) {
	return nil, nil
}
