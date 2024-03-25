package dependency

import (
	"context"

	"github.com/go-chocolate/chocolate/example/internal/app/config"
)

type dependencyContextKey struct{}

var _dependencyContextKey = &dependencyContextKey{}

func (dep *Dependency) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, _dependencyContextKey, dep)
}

func WithContext(ctx context.Context) context.Context {
	return instance.WithContext(ctx)
}

func FromContext(ctx context.Context) *Dependency {
	d, _ := ctx.Value(_dependencyContextKey).(*Dependency)
	return d
}

func Get() *Dependency {
	return instance
}

func Context() context.Context {
	return instance.WithContext(context.Background())
}

func GetGlobalConfig() config.Config {
	return instance.Config
}
