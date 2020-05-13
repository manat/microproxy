package proxy

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Route struct {
	Destinations []Destination `json:"destinations"`
	Rules        []Rule        `json:"rules"`
}

type Destination struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Default bool   `json:"default"`
}

type Rule struct {
	Path          string                 `json:"path"`
	Payload       map[string]interface{} `json:"payload"`
	DestinationID string                 `json:"destination_id"`
}

func NewRoute(filePath string) *Route {
	jf, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer jf.Close()

	jv, err := ioutil.ReadAll(jf)
	if err != nil {
		panic(err)
	}

	var route Route
	if err := json.Unmarshal(jv, &route); err != nil {
		panic(err)
	}

	return &route
}
