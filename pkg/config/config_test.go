package config_test

import (
	"testing"

	"github.com/manat/microproxy/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCreatingConfig(t *testing.T) {
	t.Run("Config can only be created once", func(t *testing.T) {
		c1 := config.Instance
		c1.FilePath = "33434343"
		c2 := config.Instance
		c3 := config.Instance

		assert.Equal(t, c1, c2)
		assert.Equal(t, c2, c3)
		assert.Equal(t, c1, c3)
	})
}
