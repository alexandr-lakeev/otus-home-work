package internalhttp

import (
	"context"
)

type Server struct { // TODO
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Panic(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
