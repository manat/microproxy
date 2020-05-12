package main

import (
	"net/http"

	"github.com/manat/microproxy/pkg/proxy"
)

func main() {
	// start server
	http.HandleFunc("/", proxy.HandleRequestThenRedirect)
	if err := http.ListenAndServe(":1338", nil); err != nil {
		panic(err)
	}
}
