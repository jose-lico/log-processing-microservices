package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jose-lico/log-processing-microservices/common/config"
	"github.com/jose-lico/log-processing-microservices/common/logging"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

var Validate = validator.New()

type RESTServer struct {
	Router *chi.Mux
}

func NewRESTServer() *RESTServer {
	return &RESTServer{Router: chi.NewRouter()}
}

func (s *RESTServer) UseDefaultMiddleware(cfg *config.RESTConfig) {
	cors := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		ExposedHeaders:   cfg.AllowedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           cfg.MaxAge,
	})

	s.Router.Use(cors.Handler)

	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
}

func (s *RESTServer) Run(ctx context.Context, cfg *config.RESTConfig) error {
	if cfg.Port == "" {
		return errors.New("no port provided")
	}

	addr := ":" + cfg.Port

	srv := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}

	errChan := make(chan error, 1)

	go func() {
		logging.Logger.Info("Starting REST server", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		logging.Logger.Info("Shutting down REST server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logging.Logger.Error("REST server shutdown failed", zap.Error(err))
			return err
		}
		logging.Logger.Info("REST server shut down")
		return nil
	case err := <-errChan:
		if err != nil {
			logging.Logger.Error("REST server encountered an error", zap.Error(err))
			return err
		}
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
