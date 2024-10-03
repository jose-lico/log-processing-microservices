package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jose-lico/log-processing-microservices/common/config"
	"github.com/jose-lico/log-processing-microservices/common/logging"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

var Validate = validator.New()

type RESTServer struct {
	Router *chi.Mux
	cfg    *config.RESTConfig
}

func NewRESTServer(cfg *config.RESTConfig) *RESTServer {
	return &RESTServer{Router: chi.NewRouter(), cfg: cfg}
}

func (s *RESTServer) UseDefaultMiddleware() {
	cors := cors.New(cors.Options{
		AllowedOrigins:   s.cfg.AllowedOrigins,
		AllowedMethods:   s.cfg.AllowedMethods,
		AllowedHeaders:   s.cfg.AllowedHeaders,
		ExposedHeaders:   s.cfg.AllowedHeaders,
		AllowCredentials: s.cfg.AllowCredentials,
		MaxAge:           s.cfg.MaxAge,
	})

	s.Router.Use(cors.Handler)

	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
}

func (s *RESTServer) Run() error {
	if s.cfg.Port == "" {
		return errors.New("no port provided")
	}

	addr := ":" + s.cfg.Port
	logging.Logger.Info("Starting API Server", zap.String("port", s.cfg.Port))
	return http.ListenAndServe(addr, s.Router)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
