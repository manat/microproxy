package main

import (
	"log"
	"net/http"

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

	var c proxy.Config
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
	config := loadConfig("config.json", true)
	log.Println(config)

	// Injecting *config to other packages
	api.AppConfig = config
	proxy.AppConfig = config

	log.Println("Booting server...")
	http.HandleFunc("/config", api.ConfigHandler)
	http.HandleFunc("/", proxy.HandleRequestThenRedirect)

	log.Println("Ready to route HTTP requests")
	if err := http.ListenAndServe(":1338", nil); err != nil {
		panic(err)
	}
}
