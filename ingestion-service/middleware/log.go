package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()
}

func LoggingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			Logger.Info("HTTP request received",
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
