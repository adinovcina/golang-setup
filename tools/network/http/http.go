package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/adinovcina/golang-setup/config"
	"github.com/adinovcina/golang-setup/tools/logger"
	"github.com/adinovcina/golang-setup/tools/network"
)

// shutdownTimeout is the time given for outstanding requests to finish before shutdown.
const shutdownTimeout = 1 * time.Second

// Instance of the service.
type Server struct {
	router        *chi.Mux
	server        *http.Server
	Configuration network.Config
}

// Setup network.
func NewServer(conf *config.Config) *Server {
	httpConfig := network.Config{
		HTTPPort:         ":" + conf.Service.Port,
		HTTPReadTimeout:  conf.Timeouts.ReadDuration,
		HTTPWriteTimeout: conf.Timeouts.WriteDuration,
	}

	router := chi.NewRouter()

	server := &http.Server{
		Addr:         httpConfig.HTTPPort,
		ReadTimeout:  httpConfig.HTTPReadTimeout,
		WriteTimeout: httpConfig.HTTPWriteTimeout,
		Handler:      router,
	}

	return &Server{
		Configuration: httpConfig,
		router:        router,
		server:        server,
	}
}

// Serve will start serving and listening http reequests.
func (s *Server) Serve() error {
	logger.Info().Msgf("Serving on port: %v", s.Configuration.HTTPPort)

	return s.server.ListenAndServe()
}

func (s *Server) Get() *chi.Mux {
	return s.router
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	defer cancel()

	return s.server.Shutdown(ctx)
}
