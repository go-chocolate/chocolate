package middleware

import (
	"net/http"
	"time"

	"github.com/go-chocolate/chocolate/pkg/toolkit/netutil"
	"github.com/sirupsen/logrus"
)

type writer struct {
	http.ResponseWriter
	statusCode int
}

func (w *writer) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.statusCode = code
}

func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			newWriter := &writer{ResponseWriter: w, statusCode: http.StatusOK}
			defer logger(newWriter, r, start)
			next.ServeHTTP(newWriter, r)
		})
	}
}

func logger(w *writer, r *http.Request, start time.Time) {
	log := logrus.WithContext(r.Context()).WithFields(map[string]interface{}{
		"method":    r.Method,
		"uri":       r.RequestURI,
		"path":      r.URL.Path,
		"query":     r.URL.RawQuery,
		"status":    w.statusCode,
		"client_ip": netutil.ClientIP(r),
		"time":      time.Since(start),
	})
	if w.statusCode >= 300 {
		log.Error()
	} else {
		log.Info()
	}
}
