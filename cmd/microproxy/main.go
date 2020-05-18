package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"

	"github.com/manat/microproxy/pkg/api"
	"github.com/manat/microproxy/pkg/config"
	"github.com/manat/microproxy/pkg/proxy"
)

func loadConfig(reload bool) {
	k := koanf.New(".")

	// Load default values
	k.Load(confmap.Provider(map[string]interface{}{
		"filepath":    "config.json",
		"server.port": "1338",
	}, "."), nil)

	// Load from Env
	k.Load(env.Provider("MICROPROXY_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "MICROPROXY_")), "_", ".", -1)
	}), nil)

	// Load from filepath (.json)
	filePath := k.String("filepath")
	f := file.Provider(filePath)
	if err := k.Load(f, json.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
		panic(err)
	}

	log.Println(k.All())
	c := config.Instance
	c.FilePath = filePath
	k.UnmarshalWithConf("", &c, koanf.UnmarshalConf{Tag: "json"})

	if reload {
		f.Watch(func(event interface{}, err error) {
			if err != nil {
				log.Printf("watch error: %v", err)
				return
			}

			log.Println("Config changed. Reloading...")
			k.Load(f, json.Parser())
			k.UnmarshalWithConf("", &c, koanf.UnmarshalConf{Tag: "json"})
		})
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	loadConfig(true)

	http.HandleFunc("/config", api.ConfigHandler)
	http.HandleFunc("/", proxy.RedirectHandler)

	log.Println("Booting server...")
	server := &http.Server{
		Addr:         ":" + config.Instance.Server.Port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	log.Println("Ready to route HTTP request")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
