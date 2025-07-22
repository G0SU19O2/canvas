// Package server contains everything for setting up and running the HTTP server.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Server struct {
	address string
	mux     chi.Router
	server  *http.Server
	log     *zap.Logger
}

type Options struct {
	Host string
	Port int
	Log  *zap.Logger
}

func New(otps Options) *Server {
	if otps.Log == nil {
		otps.Log = zap.NewNop()
	}
	address := net.JoinHostPort(otps.Host, strconv.Itoa(otps.Port))
	mux := chi.NewMux()
	httpServer := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}
	return &Server{address: address, mux: mux, server: httpServer, log: otps.Log}
}

func (s *Server) Start() error {
	s.setupRoutes()
	s.log.Info("Starting", zap.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

func (s *Server) Stop() error {
	s.log.Info("Stopping")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping server: %w", err)
	}
	return nil
}
