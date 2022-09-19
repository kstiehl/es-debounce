package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	t.Run("ApplyOptions", func(t *testing.T) {
		// given
		options := Options{}
		assert.Empty(t, options.ListenAddress, "ListenAddress should be empty")

		customOption := func(o *Options) {
			o.ListenAddress = "test"
		}

		// then
		options.ApplyOptions([]Option{customOption})

		// verify
		assert.Equal(t, "test", options.ListenAddress)
	})

	t.Run("DefaultInit", func(t *testing.T) {
		// given
		options := Options{}
		assert.Empty(t, options.ListenAddress, "ListenAddress should be empty")

		// then
		options.InitWithDefaults()

		// verify
		assert.Equal(t, ":8080", options.ListenAddress)
	})
}
