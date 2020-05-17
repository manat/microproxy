package config_test

import (
	"log"
	"testing"

	"github.com/manat/microproxy/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewRoute(t *testing.T) {
	t.Run("Route can be constructed from a JSON file", func(t *testing.T) {
		r := config.NewRoute("../../test/data/route_1.json")

		assert.IsType(t, &config.Route{}, r)
	})

	t.Run("An error is raised when the specified file is not JSON", func(t *testing.T) {
		assert.Panics(t, func() { config.NewRoute("test.csv") })
	})

	t.Run("An error is raised when the specified file does not exist", func(t *testing.T) {
		assert.Panics(t, func() { config.NewRoute("does_not_exist.json") })
	})
}

func TestReadRoute(t *testing.T) {
	route := config.NewRoute("../../test/data/route_1.json")
	log.Println(route.Rules[0].Payload["shopid"])

	t.Run("route_1.json should be marshall into Route object properly", func(t *testing.T) {
		expected := &config.Route{
			Destinations: []config.Destination{
				{ID: "mock1", URL: "http://example1.com", Default: true},
				{ID: "mock2", URL: "http://example2.com"},
			},
			Rules: []config.Rule{
				{
					Path:          "/orders/get_escrow_detail",
					Payload:       map[string]interface{}{"shopid": float64(1234)},
					DestinationID: "mock1",
				},
				{
					Path:          "/orders/get_escrow_detail",
					Payload:       map[string]interface{}{"shopid": float64(5678)},
					DestinationID: "mock2",
				},
			},
		}
		actual := route

		assert.Equal(t, expected, actual)
	})
}
