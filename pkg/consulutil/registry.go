package consulutil

import (
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
)

var (
	POD_ID = uuid.NewString()
)

type consulRegistry struct {
	client *capi.Client
}

//func (re *consulRegistry) Register(ctx context.Context, name string, end *endpoint.Endpoint) error {
//	id := POD_ID
//	return re.client.Agent().ServiceRegister(&capi.AgentServiceRegistration{
//		Kind:    capi.ServiceKindTypical,
//		ID:      id,
//		Name:    name,
//		Address: end.Host,
//		Port:    int(end.Port),
//		Meta:    end.Metadata,
//		Check:   newConsulHealthChecker(id, end),
//		//Tags:              nil,
//		//SocketPath:        "",
//		//TaggedAddresses:   nil,
//		//EnableTagOverride: false,
//		//Weights:           nil,
//		//Checks:            nil,
//		//Proxy:             nil,
//		//Connect:           nil,
//		//Namespace:         "",
//		//Partition:         "",
//		//Locality:          nil,
//	})
//}
//
//func newConsulHealthChecker(id string, end *endpoint.Endpoint) *capi.AgentServiceCheck {
//	if !end.Healthy.Enable {
//		return nil
//	}
//	check := &capi.AgentServiceCheck{
//		Interval:                       "10s",
//		Timeout:                        "1s",
//		DeregisterCriticalServiceAfter: "5s",
//	}
//
//	switch end.Healthy.Protocol {
//	case endpoint.HTTP:
//		address := end.Healthy.Address
//		if address == "" {
//			if end.Healthy.TLS {
//				address = fmt.Sprintf("https://%s:%d/_healthy_", end.Host, end.Port)
//			} else {
//				address = fmt.Sprintf("http://%s:%d/_healthy_", end.Host, end.Port)
//			}
//		}
//		method := http.MethodGet
//		if end.Healthy.HTTPMethod != "" {
//			method = end.Healthy.HTTPMethod
//		}
//
//		check.CheckID = id + "_http"
//		check.Name = "health_check_http"
//		check.HTTP = address
//		check.Method = method
//		check.TLSSkipVerify = !end.Healthy.TLS
//	case endpoint.GRPC:
//		check.CheckID = id + "_grpc"
//		check.Name = "health_check_grpc"
//		check.GRPC = fmt.Sprintf("%s:%d", end.Host, end.Port)
//		check.GRPCUseTLS = end.Healthy.TLS
//	default:
//		return nil
//	}
//	return check
//}
