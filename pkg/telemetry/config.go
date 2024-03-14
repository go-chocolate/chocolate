package telemetry

type Config struct {
	ServiceName    string
	ServiceVersion string
	Attributes     map[string]string
	Sampler        float64
}
