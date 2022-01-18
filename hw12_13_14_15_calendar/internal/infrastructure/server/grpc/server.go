package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	appcalendar "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/calendar"
	deliverygrpc "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/calendar/delivery/grpc"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/calendar/delivery/grpc/pb"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"google.golang.org/grpc"
)

type Service struct {
	config  config.GrpcConf
	handler *deliverygrpc.Handler
	logg    app.Logger
	pb.UnimplementedCalendarServer
}

func NewServer(cfg config.GrpcConf, useCase appcalendar.UseCase, logger app.Logger) *Service {
	return &Service{
		config:  cfg,
		handler: deliverygrpc.NewHandler(useCase),
		logg:    logger,
	}
}

func (s *Service) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.config.ListenAddress)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterCalendarServer(server, s)

	s.logg.Info(fmt.Sprintf("grpc server is starting on %s", lsn.Addr()))

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
