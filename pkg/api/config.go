package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/manat/microproxy/pkg/proxy"
)

// AppConfig provides a way for injecting config from the main package.
var AppConfig *proxy.Config

// ConfigHandler handle read, and updating config for MicroProxy.
func ConfigHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case "GET":
		cfg, err := json.Marshal(AppConfig)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(cfg)
	case "PUT":
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		f, err := os.OpenFile("config.json", os.O_WRONLY, os.ModeExclusive)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer f.Close()

		n, err := f.Write(b)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		log.Printf("Successfully saved config.json. Bytes written: %d", n)

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(`{"message": "updated successfully"}`))
	default:
		res.WriteHeader(http.StatusNotImplemented)
		res.Write([]byte(`{"message": "Only GET and PUT is supported"}`))
	}
}
