package api

import (
	"net/http"
)

func ConfigHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case "GET":
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"message": "here is the current config"}`))
	case "POST", "PUT":
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(`{"message": "created successfully"}`))
	default:
		res.WriteHeader(http.StatusNotImplemented)
		res.Write([]byte(`{"message": "Only GET and POST is supported"}`))
	}
}
