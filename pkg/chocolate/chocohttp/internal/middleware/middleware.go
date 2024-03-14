package middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

var nopMiddleware = Middleware(func(next http.Handler) http.Handler {
	return next
})

func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares); i > 0; i-- {
			next = middlewares[i-1](next)
		}
		return next
	}
}
