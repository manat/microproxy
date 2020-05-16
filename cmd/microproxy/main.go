package main

import (
	"log"
	"net/http"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"

	"github.com/manat/microproxy/pkg/api"
	"github.com/manat/microproxy/pkg/proxy"
)

func loadConfig(filePath string, reload bool) *proxy.Config {
	k := koanf.New(".")
	f := file.Provider(filePath)
	if err := k.Load(f, json.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Println("Route = ", k.String("route"))

	c := proxy.Config{FilePath: &filePath}
	k.Unmarshal("", &c)

	if reload {
		f.Watch(func(event interface{}, err error) {
			if err != nil {
				log.Printf("watch error: %v", err)
				return
			}

			log.Println("Config changed. Reloading...")
			k.Load(f, json.Parser())
			k.Unmarshal("", &c)
		})
	}

	return &c
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config := loadConfig("config.json", true)
	log.Println(config)

	// Injecting *config to other packages
	api.ProxyConfig = config

	http.HandleFunc("/config", api.ConfigHandler)
	http.HandleFunc("/", proxy.HandleRequestThenRedirect)

	log.Println("Booting server...")
	server := &http.Server{
		Addr:         ":1338",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	log.Println("Ready to route HTTP request")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
