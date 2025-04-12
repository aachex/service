package logging

import (
	"log/slog"
	"net/http"
)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func Middleware(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("incoming request",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("userAgent", r.UserAgent()))

		lrw := &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next(lrw, r)

		logger.Debug("request handled",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.Int("statusCode", lrw.statusCode),
			slog.String("responseBody", string(lrw.body)))
	}
}
