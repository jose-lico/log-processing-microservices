package middleware

import (
	"net/http"
	"time"

	"github.com/jose-lico/log-processing-microservices/common/logging"

	"go.uber.org/zap"
)

func LoggingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			logging.Logger.Info("HTTP request received",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", ww.statusCode),
				zap.Duration("duration", duration),
				zap.String("client_ip", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
