package middleware

import "github.com/go-chocolate/chocolate/pkg/chocolate/chocohttp/internal/middleware/cors"

func CORS(c cors.Config) Middleware {
	return cors.New(c)
}
