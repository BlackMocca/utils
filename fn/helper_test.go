package fn_test

import (
	"os"
	"testing"

	"github.com/BlackMocca/utils/fn"
	"github.com/stretchr/testify/assert"
)

const (
	envKey       = "EXAMPLE_TEST"
	envValue     = "abcd1234"
	defaultValue = "defaultvalue1234"
)

func TestLookupEnv(t *testing.T) {
	t.Run("lookup not found in env", func(t *testing.T) {
		val := fn.LookupEnv(envKey, defaultValue)
		assert.Equal(t, val, defaultValue)
	})
	t.Run("lookup found in env", func(t *testing.T) {
		os.Setenv(envKey, envValue)

		val := fn.LookupEnv(envKey, defaultValue)
		assert.Equal(t, val, envValue)
	})
}
