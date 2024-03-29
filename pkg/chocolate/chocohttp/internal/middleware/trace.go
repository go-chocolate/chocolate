package middleware

import (
	"net/http"

	"github.com/go-chocolate/chocolate/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	headerTraceId = "X-Trace-Id"
)

func Trace(name string, options ...otelhttp.Option) Middleware {
	return otelhttp.NewMiddleware(name, options...)
}

func TraceId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set(headerTraceId, telemetry.TraceIDFromContext(request.Context()))
			next.ServeHTTP(writer, request)
		})
	}
}
