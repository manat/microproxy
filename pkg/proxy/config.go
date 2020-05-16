package proxy

import (
	"log"
	"os"
)

// AppConfig provides a way for injecting config from the main package.
var AppConfig *Config

// Config represents structure of MicroProxy config.
type Config struct {
	FilePath *string
	Route    *Route `json:"route"`
}

// Save persists the specified value into file, as specified on Config.FilePath.
func (c *Config) Save(b []byte) (int, error) {
	if c.FilePath == nil {
		panic("*Config must have FilePath value specified.")
	}

	f, err := os.OpenFile(*c.FilePath, os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	n, err := f.Write(b)
	if err != nil {
		return -1, err
	}

	log.Printf("Successfully saved config to %s. Bytes written: %d", f.Name(), n)
	return n, nil
}
