package http

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
