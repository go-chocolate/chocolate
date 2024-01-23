package registry

import (
	"context"

	"github.com/go-chocolate/chocolate/pkg/service/cluster/endpoint"
)

type Registry interface {
	Register(ctx context.Context, end *endpoint.Endpoint) error
}
