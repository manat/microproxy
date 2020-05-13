package proxy

// AppConfig provides a way for injecting config from the main package.
var AppConfig *Config

// Config represents structure of MicroProxy config.
type Config struct {
	Route *Route `json:"route"`
}

// NewConfig constructs MicroProxy's Config from files, and environment values.
func NewConfig(r *Route) *Config {
	return &Config{
		Route: r,
	}
}
