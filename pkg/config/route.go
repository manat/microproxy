package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// Route represents info required for forwarding request to multiple destinations,
// based on the specified rules.
type Route struct {
	Destinations []Destination `json:"destinations"`
	Rules        []Rule        `json:"rules"`
}

// Destination represents info of the destination node.
type Destination struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Default bool   `json:"default"`
}

// Rule represents info of how the MicroProxy determines routes based on info such as path,
// and payload.
type Rule struct {
	Path          string                 `json:"path"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	DestinationID string                 `json:"destination_id"`
}

// NewRoute constructs Route based on the given JSON file.
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

func (r *Route) GetDestinationByID(id string) (*Destination, error) {
	for _, d := range r.Destinations {
		if d.ID == id {
			return &d, nil
		}
	}

	msg := fmt.Sprintf("Destination ID: %s does not exist", id)
	return nil, errors.New(msg)
}
