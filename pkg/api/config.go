package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/manat/microproxy/pkg/proxy"
)

// ProxyConfig provides a way for injecting config from the main package.
var ProxyConfig *proxy.Config

// ConfigHandler handle read, and updating config for MicroProxy.
func ConfigHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case "GET":
		cfg, err := json.Marshal(ProxyConfig)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(cfg)
	case "PUT":
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer req.Body.Close()

		_, err = ProxyConfig.Save(b)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(`{"message": "Config was updated successfully"}`))
	default:
		res.WriteHeader(http.StatusNotImplemented)
		res.Write([]byte(`{"message": "Only GET and PUT is supported"}`))
	}
}
