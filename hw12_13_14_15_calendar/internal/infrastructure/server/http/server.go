package internalhttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app"
	deliveryhttp "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/delivery/http"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/gorilla/mux"
)

type Server struct {
	server  *http.Server
	usecase app.UseCase
	logg    app.Logger
}

func NewServer(cfg config.ServerConf, usecase app.UseCase, logger app.Logger) *Server {
	router := mux.NewRouter()
	handler := deliveryhttp.NewHandler(usecase)

	router.Use(newLoggingMiddleware(logger))

	ctx := context.Background()
	router.HandleFunc("/api/calendar/v1/events", handler.CreateEvent(ctx)).Methods("POST")
	router.HandleFunc("/api/calendar/v1/events/{id:[0-9\\-a-f]+}", handler.GetEvent(ctx)).Methods("GET")
	router.HandleFunc("/api/calendar/v1/events/{id:[0-9\\-a-f]+}", handler.UpdateEvent(ctx)).Methods("POST")

	return &Server{
		server: &http.Server{
			Handler:      router,
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
