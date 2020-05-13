package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manat/microproxy/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestConfigHandler(t *testing.T) {
	const ConfigPath = "/config"

	t.Run("GET config should return 200 OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(api.ConfigHandler)
		req, err := http.NewRequest(http.MethodGet, ConfigPath, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rec, req)

		expectedCode := http.StatusOK
		actualCode := rec.Code
		assert.Equal(t, expectedCode, actualCode)

		expectedBody := `{"message": "here is the current config"}`
		actualBody := rec.Body.String()
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("PUT config should return 201 CREATED", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(api.ConfigHandler)
		req, err := http.NewRequest(http.MethodPut, ConfigPath, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rec, req)

		expectedCode := http.StatusCreated
		actualCode := rec.Code
		assert.Equal(t, expectedCode, actualCode)

		expectedBody := `{"message": "updated successfully"}`
		actualBody := rec.Body.String()
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("PATCH config should return 501 NOT IMPLEMENTED", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(api.ConfigHandler)
		req, err := http.NewRequest(http.MethodPatch, ConfigPath, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rec, req)

		expectedCode := http.StatusNotImplemented
		actualCode := rec.Code
		assert.Equal(t, expectedCode, actualCode)

		expectedBody := `{"message": "Only GET and PUT is supported"}`
		actualBody := rec.Body.String()
		assert.Equal(t, expectedBody, actualBody)
	})
}
