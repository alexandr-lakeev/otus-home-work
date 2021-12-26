package deliverygrpc

import (
	"context"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/delivery/grpc/pb"
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
	return nil, nil
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
