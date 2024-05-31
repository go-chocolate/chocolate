package registry

import (
	"context"
	"github.com/go-chocolate/chocolate/pkg/chocolate/cluster/endpoint"
	"github.com/go-chocolate/chocolate/pkg/toolkit/netutil"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Registry interface {
	Register(ctx context.Context, name string, endpoint *endpoint.Endpoint) error
}

type RegistryFunc func(ctx context.Context, name string, endpoint *endpoint.Endpoint) error

func (f RegistryFunc) Register(ctx context.Context, name string, endpoint *endpoint.Endpoint) error {
	return f(ctx, name, endpoint)
}

var (
	// NoneRegistry do nothing
	NoneRegistry = RegistryFunc(func(ctx context.Context, name string, endpoint *endpoint.Endpoint) error {
		return nil
	})
)

// Registries register service to multi data center
type Registries []Registry

func (r Registries) Register(ctx context.Context, name string, end *endpoint.Endpoint) error {

	// 这里是为了将闭包函数移到循环体外
	// 循环体内写闭包函数在某些情况下可能导致预料之外的异常执行结果（go1.23+已优化）
	var f = func(r Registry, c context.Context, name string, end *endpoint.Endpoint) func() error {
		return func() error {
			return r.Register(c, name, end)
		}
	}

	group, groupCtx := errgroup.WithContext(ctx)
	for _, v := range r {
		group.Go(f(v, groupCtx, name, end))
	}
	return group.Wait()

}

type autoIPRegistry struct {
	parent Registry
}

func (r *autoIPRegistry) Register(ctx context.Context, name string, end *endpoint.Endpoint) error {
	cloned := new(endpoint.Endpoint)
	*cloned = *end
	cloned.Host = netutil.FigureOutIP(end.Host)
	if cloned.Host != end.Host {
		logrus.Infof("registerd address repleced: %s => %s", end.Host, cloned.Host)
	}
	return r.parent.Register(ctx, name, end)
}

// AutoIPRegistry
// 当未设置ip地址时尝试自动读取本机网卡ip
func AutoIPRegistry(registry Registry) Registry {
	return &autoIPRegistry{parent: registry}
}
