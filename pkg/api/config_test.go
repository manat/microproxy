package api_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/manat/microproxy/pkg/api"
	"github.com/manat/microproxy/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestConfigHandler(t *testing.T) {
	const apiConfigPath = "/config"

	configFilePath := "../../test/data/route_1.json"
	api.ProxyConfig = &proxy.Config{
		FilePath: &configFilePath,
		Route:    proxy.NewRoute(configFilePath),
	}

	t.Run("GET config should return 200 OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(api.ConfigHandler)
		req, err := http.NewRequest(http.MethodGet, apiConfigPath, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rec, req)

		expectedCode := http.StatusOK
		actualCode := rec.Code
		assert.Equal(t, expectedCode, actualCode)

		expectedBody, err := json.Marshal(&api.ProxyConfig)
		if err != nil {
			panic(err)
		}

		actualBody := rec.Body.String()
		assert.Equal(t, string(expectedBody), actualBody)
	})

	t.Run("PUT config should return 201 CREATED", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(api.ConfigHandler)
		// So that we don't have to be concerned of messing the original data file.
		tmpFile, err := ioutil.TempFile("", "config")
		if err != nil {
			panic(err)
		}
		defer os.Remove(tmpFile.Name())
		tmpFilePath := tmpFile.Name()
		api.ProxyConfig.FilePath = &tmpFilePath

		cfg, err := json.Marshal(api.ProxyConfig)
		if err != nil {
			panic(err)
		}
		body := bytes.NewBuffer(cfg)

		req, err := http.NewRequest(http.MethodPut, apiConfigPath, body)
		if err != nil {
			t.Fatal(err)
		}
		handler.ServeHTTP(rec, req)

		expectedCode := http.StatusCreated
		actualCode := rec.Code
		assert.Equal(t, expectedCode, actualCode)

		expectedBody := `{"message": "Config was updated successfully"}`
		actualBody := rec.Body.String()
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("PATCH config should return 501 NOT IMPLEMENTED", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(api.ConfigHandler)
		req, err := http.NewRequest(http.MethodPatch, apiConfigPath, nil)
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
