package http_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/manat/microproxy/pkg/http"

	"github.com/stretchr/testify/assert"
)

func TestReadRoute(t *testing.T) {
	jf, err := os.Open("../../test/data/route_1.json")
	if err != nil {
		panic(err)
	}
	defer jf.Close()

	jv, err := ioutil.ReadAll(jf)
	if err != nil {
		panic(err)
	}

	var route http.Route
	json.Unmarshal(jv, &route)
	log.Println(route.Rules[0].Payload["shopid"])

	t.Run("route_1.json should be marshall into Route object properly", func(t *testing.T) {
		expected := &http.Route{
			Destinations: []http.Destination{
				{ID: "mock1", URL: "http://example1.com", Default: true},
				{ID: "mock2", URL: "http://example2.com"},
			},
			Rules: []http.Rule{
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
		actual := &route

		assert.Equal(t, expected, actual)
	})
}
