package chocorpc

import (
	"fmt"
	"time"

	"github.com/go-chocolate/chocolate/pkg/chocolate/chocohttp/internal/middleware/cors"
)

type Config struct {
	Name    string
	Addr    string
	TLS     *TLSConfig
	Options OptionsConfig
}

type TLSConfig struct {
	KeyFile  string
	CertFile string
}

type OptionsConfig struct {
	CORS      CORSConfig
	Logger    LoggerConfig
	RateLimit RateLimitConfig
}

type CORSConfig struct {
	Enable          bool
	AbortOnError    bool
	AllowAllOrigins bool

	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// Default value is ["*"]
	AllowedOrigins []string

	//// AllowOriginFunc is a custom function to validate the origin. It take the origin
	//// as argument and returns true if allowed or false otherwise. If this option is
	//// set, the content of AllowedOrigins is ignored.
	//AllowOriginFunc func(origin string) bool

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET and POST)
	AllowedMethods []string

	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge string
}

func (c CORSConfig) build() cors.Config {
	cfg := cors.DefaultConfig()
	cfg.AbortOnError = c.AbortOnError
	cfg.AllowAllOrigins = c.AllowAllOrigins
	cfg.AllowCredentials = c.AllowCredentials
	if c.AllowedOrigins != nil {
		cfg.AllowedOrigins = c.AllowedOrigins
	}
	if c.AllowedMethods != nil {
		cfg.AllowedMethods = c.AllowedMethods
	}
	if c.AllowedHeaders != nil {
		cfg.AllowedHeaders = c.AllowedHeaders
	}
	if c.ExposedHeaders != nil {
		cfg.ExposedHeaders = c.ExposedHeaders
	}
	if c.MaxAge != "" {
		var err error
		cfg.MaxAge, err = time.ParseDuration(c.MaxAge)
		if err != nil {
			panic(fmt.Errorf("chocohttp.Config.Middleware.CORS.MaxAge:cannot decode string value '%s' to duration, %v", c.MaxAge, err))
		}
	}
	return cfg
}

type LoggerConfig struct {
	Enable             bool
	RecordHeader       []string
	RecordRequestBody  bool
	RecordResponseBody bool
}

type RateLimitConfig struct {
	Enable bool
	Limit  int
}
