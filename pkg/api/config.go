package api

import (
	"net/http"
)

// ConfigHandler handle read, and updating config for MicroProxy.
func ConfigHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case "GET":
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"message": "here is the current config"}`))
	case "PUT":
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(`{"message": "updated successfully"}`))
	default:
		res.WriteHeader(http.StatusNotImplemented)
		res.Write([]byte(`{"message": "Only GET and PUT is supported"}`))
	}
}
