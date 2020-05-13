package proxy_test

import (
	"testing"

	"github.com/manat/microproxy/pkg/proxy"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("Config should hold values required for MicroProxy app", func(t *testing.T) {
		route := proxy.NewRoute("../../test/data/route_1.json")
		cfg := proxy.NewConfig(route)

		assert.Equal(t, cfg.Route, route)
		assert.Equal(t, 2, len(cfg.Route.Destinations))
		assert.Equal(t, 2, len(cfg.Route.Rules))
	})
}
