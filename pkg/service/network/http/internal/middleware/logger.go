package middleware

import "github.com/go-chocolate/chocolate/pkg/service/network/http/internal/logger"

func Logger() Middleware {
	return logger.Logger(logger.WithIgnorePath("/__health__"))
}
