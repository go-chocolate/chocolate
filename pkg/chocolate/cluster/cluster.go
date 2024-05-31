package cluster

import (
	"context"
	"github.com/go-chocolate/chocolate/pkg/chocolate/cluster/endpoint"
	"github.com/go-chocolate/chocolate/pkg/chocolate/cluster/registry"
)

type Cluster struct {
	registry registry.Registry
}

func (c *Cluster) WithClusterRegistry(r registry.Registry) {
	c.registry = r
	//c.registry = registry.AutoIPRegistry(r)
}

func (c *Cluster) ClusterRegister(ctx context.Context, name string, endpoint *endpoint.Endpoint) error {
	if c.registry == nil {
		return nil
	}
	return c.registry.Register(ctx, name, endpoint)
}
