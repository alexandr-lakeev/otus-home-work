package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
)

type Server struct {
	server  *http.Server
	usecase app.UseCase
	logg    Logger
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
	Panic(msg string)
}

func NewServer(cfg config.ServerConf, usecase app.UseCase, logger Logger) *Server {
	return &Server{
		server: &http.Server{
			Handler:      &dummyHandler{},
			Addr:         cfg.BindAddress,
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		logg:    logger,
		usecase: usecase,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logg.Info(fmt.Sprintf("server is starting on %s", s.server.Addr))

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

type dummyHandler struct{}

func (h *dummyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}
