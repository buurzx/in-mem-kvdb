package compute

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestParse(t *testing.T) {
	l := zap.NewNop()
	c, err := New(l)

	require.NoError(t, err)

	t.Run("empty request", func(t *testing.T) {
		_, err := c.Parse("")
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		assert.Equal(t, errEmptyRequest, err)
	})

	t.Run("invalid command", func(t *testing.T) {
		_, err := c.Parse("invalid")
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		assert.Equal(t, errInvalidCommand, err)
	})

	t.Run("invalid arguments", func(t *testing.T) {
		_, err := c.Parse("GET")
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		assert.Equal(t, errInvalidArguments, err)
	})

	t.Run("valid request", func(t *testing.T) {
		query, err := c.Parse("GET key")
		require.NoError(t, err)

		assert.Equal(t, Query{
			commandID: GetCommandID,
			args:      []string{"key"},
		}, query)
	})
}
