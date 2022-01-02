package deliverygrpc

import (
	"context"

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
	userID, err := uuid.Parse(query.UserId)
	if err != nil {
		return nil, err
	}

	events, err := h.useCase.GetList(ctx, &app.GetListQuery{
		UserID: userID,
		From:   query.From.AsTime(),
		To:     query.To.AsTime(),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*pb.EventResponse, len(events))
	for key, event := range events {
		result[key] = &pb.EventResponse{
			Id:          event.ID.String(),
			Title:       event.Title,
			Date:        timestamppb.New(event.Date),
			Duration:    durationpb.New(event.Duration),
			Description: event.Description,
		}
	}

	return &pb.GetListResponse{
		Events: result,
	}, nil
}

func (h *Handler) CreateEvent(ctx context.Context, command *pb.CreateEventCommand) (*pb.CreateEventResponse, error) {
	id := uuid.New()

	userID, err := uuid.Parse(command.UserId)
	if err != nil {
		return nil, err
	}

	event := command.Event

	err = h.useCase.CreateEvent(ctx, &app.CreateEventCommand{
		ID:          id,
		UserID:      userID,
		Title:       event.Title,
		Date:        event.Date.AsTime(),
		Duration:    event.Duration.AsDuration(),
		Description: event.Description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateEventResponse{
		Id: id.String(),
	}, nil
}

func (h *Handler) UpdateEvent(ctx context.Context, command *pb.UpdateEventCommand) (*pb.UpdateEventResponse, error) {
	id, err := uuid.Parse(command.Id)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(command.UserId)
	if err != nil {
		return nil, err
	}

	event := command.Event

	err = h.useCase.UpdateEvent(ctx, &app.UpdateEventCommand{
		ID:          id,
		UserID:      userID,
		Title:       event.Title,
		Date:        event.Date.AsTime(),
		Duration:    event.Duration.AsDuration(),
		Description: event.Description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateEventResponse{}, nil
}
