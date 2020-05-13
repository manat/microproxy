package proxy

// Config is a centralized object that holds every configuration value.
type Config struct {
	Route *Route
}

// NewConfig constructs MicroProxy's Config from files, and environment values.
func NewConfig(r *Route) *Config {
	return &Config{
		Route: r,
	}
}
