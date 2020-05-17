package proxy

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/manat/microproxy/pkg/config"
)

func RedirectHandler(res http.ResponseWriter, req *http.Request) {
	dest := GetDestination(req)

	serveReverseProxy(dest.URL, res, req)
}

func GetDestination(req *http.Request) *config.Destination {
	var dest *config.Destination

	data, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	sdata := string(data)
	sbody := strings.Split(sdata, "\n")

	for _, r := range config.Instance.Route.Rules {
		if strings.Contains(req.URL.String(), r.Path) {
			route := config.Instance.Route
			dest, err = route.GetDestinationByID(r.DestinationID)
			if err != nil {
				panic(err)
			}

			actualBody := strings.TrimSpace(sbody[len(sbody)-1])
			if actualBody != "" {
				jsonMap := make(map[string]interface{})
				if err = json.Unmarshal([]byte(actualBody), &jsonMap); err != nil {
					panic(err)
				}
				log.Println(jsonMap)
				log.Println(r.Payload)

				if intersect(jsonMap, r.Payload) {
					dest, err = route.GetDestinationByID(r.DestinationID)
					if err != nil {
						panic(err)
					}
					return dest
				}
			}
		}
	}

	return dest
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

// intersect was copied from https://stackoverflow.com/a/31069616/136492
func intersect(as, bs map[string]interface{}) bool {
	for _, a := range as {
		for _, b := range bs {
			if a == b {
				return true
			}
		}
	}

	return false
}
