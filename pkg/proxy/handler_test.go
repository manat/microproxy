package proxy_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/manat/microproxy/pkg/config"
	"github.com/manat/microproxy/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestGetDestination(t *testing.T) {
	cfg := config.Instance
	cfg.Route = *config.NewRoute("../../test/data/route_1.json")

	t.Run("When path is /expenses, then destination should be mock3", func(t *testing.T) {
		var body []byte
		req, err := http.NewRequest(http.MethodGet, "/expenses", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		expected := cfg.Route.Destinations[2]
		actual := proxy.GetDestination(req)

		assert.Equal(t, expected, *actual)
	})

	t.Run("When path is /orders/escrow_detail, and shopid is 5678, then destination should be mock2",
		func(t *testing.T) {
			payload := map[string]interface{}{"shopid": float64(5678)}
			payloadJson, _ := json.Marshal(payload)
			body := bytes.NewBuffer(payloadJson)

			req, err := http.NewRequest(http.MethodPost, "/orders/escrow_detail", body)
			if err != nil {
				t.Fatal(err)
			}

			expected := cfg.Route.Destinations[1]
			actual := proxy.GetDestination(req)

			assert.Equal(t, expected, *actual)
		})

	t.Run("When path is /orders/escrow_detail, and shopid is 1234, then destination should be mock1",
		func(t *testing.T) {
			payload := map[string]interface{}{"shopid": float64(1234)}
			payloadJson, _ := json.Marshal(payload)
			body := bytes.NewBuffer(payloadJson)

			req, err := http.NewRequest(http.MethodPost, "/orders/escrow_detail", body)
			if err != nil {
				t.Fatal(err)
			}

			expected := cfg.Route.Destinations[0]
			actual := proxy.GetDestination(req)

			assert.Equal(t, expected, *actual)
		})

	t.Run(`When path is /orders/escrow_detail, and shopid is 1234 with lengthy body,
		   then destination should be mock1`,
		func(t *testing.T) {
			payload := map[string]interface{}{"shopid": float64(1234), "partnerid": float64(1232344)}
			payloadJson, _ := json.Marshal(payload)
			body := bytes.NewBuffer(payloadJson)

			req, err := http.NewRequest(http.MethodPost, "/orders/escrow_detail", body)
			if err != nil {
				t.Fatal(err)
			}

			expected := cfg.Route.Destinations[0]
			actual := proxy.GetDestination(req)

			assert.Equal(t, expected, *actual)
		})
}
