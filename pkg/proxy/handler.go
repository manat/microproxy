package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func RedirectHandler(res http.ResponseWriter, req *http.Request) {
	url := getDestination(req)

	serveReverseProxy(url, res, req)

}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Updates the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func getDestination(req *http.Request) string {
	data, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(string(data), "\"partnerid\": 12345") {
		url := "http://127.0.0.1:8002/test"
		log.Println(url)

		return url
	}

	return "http://127.0.0.1:8000/test"
}
