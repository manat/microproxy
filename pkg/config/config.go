package config

import (
	"log"
	"os"
	"sync"
)

// Instance holds reference of singleton Config
var Instance *Config

// once indicates wheter or not Instance has been created
var once sync.Once

// Config represents structure of MicroProxy config.
type Config struct {
	FilePath string `json:"file_path"`
	Server   Server `json:"server"`
	Route    Route  `json:"route"`
}

// Make sure Instance is ready when config package is referred.
func init() {
	newConfig()
}

// newConfig makes sure config.Instance can be created only once.
func newConfig() *Config {
	once.Do(func() {
		Instance = new(Config)
	})

	return Instance
}

// Save persists the specified value into the file, as specified on Config.FilePath.
func (c *Config) Save(b []byte) (int, error) {
	if c.FilePath == "" {
		panic("*Config must have FilePath value specified.")
	}

	f, err := os.OpenFile(c.FilePath, os.O_WRONLY, os.ModeExclusive)
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
