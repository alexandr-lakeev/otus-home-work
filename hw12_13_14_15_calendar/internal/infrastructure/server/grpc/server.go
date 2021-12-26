package internalgrpc

import (
	"context"
	"log"
	"net"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	deliverygrpc "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/delivery/grpc"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/delivery/grpc/pb"
	"google.golang.org/grpc"
)

type Service struct {
	handler *deliverygrpc.Handler
	logg    app.Logger
	pb.UnimplementedCalendarServer
}

func NewServer(useCase app.UseCase, logger app.Logger) *Service {
	handler := deliverygrpc.NewHandler(useCase)

	return &Service{
		handler: handler,
		logg:    logger,
	}
}

func (s *Service) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", ":50051") // nolint
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterCalendarServer(server, new(Service))

	log.Printf("grpc server is starting on %s", lsn.Addr().String())

	return server.Serve(lsn)
}

func (s *Service) GetEvent(ctx context.Context, query *pb.GetEventQuery) (*pb.GetEventResponse, error) {
	return s.handler.GetEvent(ctx, query)
}

func (s *Service) GetList(ctx context.Context, query *pb.GetListQuery) (*pb.GetListResponse, error) {
	return s.handler.GetList(ctx, query)
}

func (s *Service) CreateEvent(ctx context.Context, command *pb.CreateEventCommand) (*pb.CreateEventResponse, error) {
	return s.handler.CreateEvent(ctx, command)
}

func (s *Service) UpdateEvent(ctx context.Context, command *pb.UpdateEventCommand) (*pb.UpdateEventResponse, error) {
	return s.handler.UpdateEvent(ctx, command)
}
