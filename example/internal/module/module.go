package module

import (
	"github.com/go-chocolate/chocolate/example/internal/app/dependency"
	"github.com/go-chocolate/chocolate/example/internal/module/auth"
)

func Setup() error {
	services.Auth = auth.NewService(dependency.Get().KVStorage)

	return nil
}

var services = &struct {
	Auth auth.Service
}{}

func AuthService() auth.Service {
	return services.Auth
}
