package proxy

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/manat/microproxy/pkg/config"
)

func RedirectHandler(res http.ResponseWriter, req *http.Request) {
	dest := GetDestination(req)
	log.Println(dest)

	if dest == nil {
		serveReverseProxy(req.URL.String(), res, req)
		return
	}

	if dest.ID == "" {
		serveReverseProxy(req.URL.String(), res, req)
		return
	}

	if dest.Host == "" {
		serveReverseProxy(dest.URL, res, req)
		return
	}

	re := regexp.MustCompile(`(https?://.*):(\d*)`)
	newURL := re.ReplaceAllString(req.URL.String(), dest.Host)
	serveReverseProxy(newURL, res, req)
	return
}

func GetDestination(req *http.Request) *config.Destination {
	var dest *config.Destination

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	for _, r := range config.Instance.Route.Rules {
		if strings.Contains(req.URL.String(), r.Path) {
			route := config.Instance.Route
			dest, err = route.GetDestinationByID(r.DestinationID)
			if err != nil {
				panic(err)
			}

			if len(body) > 0 {
				jsonMap := make(map[string]interface{})
				if err = json.Unmarshal([]byte(body), &jsonMap); err != nil {
					panic(err)
				}
				log.Println(jsonMap)
				log.Println(r.Payload)

				if intersect(jsonMap, r.Payload) {
					dest, err = route.GetDestinationByID(r.DestinationID)
					if err != nil {
						panic(err)
					}
					log.Println(dest)
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
	ps := httputil.NewSingleHostReverseProxy(url)

	// Updates the headers to allow for SSL redirection
	req.URL.Host = url.Host
	// req.URL.Path = ""
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	ps.ServeHTTP(res, req)
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
