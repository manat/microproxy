package main

import (
	"log"
	"net/http"

	"github.com/manat/microproxy/pkg/api"

	"github.com/manat/microproxy/pkg/proxy"
)

func main() {
	log.Println("Booting server...")
	http.HandleFunc("/config", api.ConfigHandler)
	http.HandleFunc("/", proxy.HandleRequestThenRedirect)

	log.Println("Ready to route HTTP requests")
	if err := http.ListenAndServe(":1338", nil); err != nil {
		panic(err)
	}
}
