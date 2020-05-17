package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/manat/microproxy/pkg/config"
)

// ConfigHandler handle read, and updating config for MicroProxy.
func ConfigHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case "GET":
		cfg, err := json.Marshal(config.Instance)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write(cfg)
	case "PUT":
		r := http.MaxBytesReader(res, req.Body, 200_000)
		b, err := ioutil.ReadAll(r)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer req.Body.Close()

		_, err = config.Instance.Save(b)
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
