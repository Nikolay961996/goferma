package router

import (
	"context"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/services"
	"github.com/Nikolay961996/goferma/internal/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type (
	responseDate struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		data *responseDate
	}
)

func (w *loggingResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.data.status = status
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.data.size += size
	return size, err
}

func WithLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		rd := &responseDate{
			status: 200,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			data:           rd,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		utils.Log.WithFields(logrus.Fields{
			"method":   method,
			"uri":      uri,
			"duration": duration,
			"status":   lw.data.status,
			"size":     lw.data.size,
		}).Info("request log")
	})
}

func WithAuth(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := services.GetUserID(r.Header.Get("Authorization"), secretKey)
			if err != nil {
				utils.Log.Error("error login/password pair:", err.Error())
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), models.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
