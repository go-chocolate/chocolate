package balancer

import "github.com/go-chocolate/chocolate/pkg/service/cluster/endpoint"

type Balancer interface {
	LB(endpoints []*endpoint.Endpoint) (*endpoint.Endpoint, error)
}
