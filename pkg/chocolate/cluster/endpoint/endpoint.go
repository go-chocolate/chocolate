package endpoint

const (
	HTTP = "http"
	GRPC = "grpc"
)

type Endpoint struct {
	Protocol string
	Host     string
	Port     uint16
	Metadata map[string]string
	Healthy  HealthyOption
}

type HealthyOption struct {
	Enable     bool
	Protocol   string
	Address    string
	Path       string
	TLS        bool
	HTTPMethod string
}
