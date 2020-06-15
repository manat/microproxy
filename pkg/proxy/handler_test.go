package proxy_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manat/microproxy/pkg/config"
	"github.com/manat/microproxy/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestGetRedirectHandler(t *testing.T) {
	cfg := config.Instance
	cfg.Route = *config.NewRoute("../../test/data/route_2.json")

	t.Run("When path is /orders/basics, and shopid is 250725, then it should redirect using mock3", func(t *testing.T) {
		url := "http://localhost:12011/api/v1/orders/basics"
		payload := map[string]interface{}{"shopid": float64(250725)}
		payloadJson, _ := json.Marshal(payload)
		body := bytes.NewBuffer(payloadJson)

		req, err := http.NewRequest(http.MethodPost, url, body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(proxy.RedirectHandler)
		handler.ServeHTTP(rr, req)

		expected := "http://dev-channel-bridge-1a-0.acommercedev.platform:12021/api/v1/orders/basics"
		actual := req.URL.String()

		assert.Equal(t, expected, actual)
	})

}

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
	t.Run(`When path is /expenses/12345, and shopid is 5678 with lengthy body,
		then destination should be mock4's host with the remaining path`,
		func(t *testing.T) {
			payload := map[string]interface{}{"shopid": float64(5678), "partnerid": float64(1232344)}
			payloadJson, _ := json.Marshal(payload)
			body := bytes.NewBuffer(payloadJson)

			req, err := http.NewRequest(http.MethodPost, "/expenses/12345", body)
			if err != nil {
				t.Fatal(err)
			}

			expected := cfg.Route.Destinations[3]
			actual := proxy.GetDestination(req)

			assert.Equal(t, expected, *actual)
		})
}
