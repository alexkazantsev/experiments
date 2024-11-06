package monitor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Do(t *testing.T) {
	t.Run("should do correct", func(t *testing.T) {
		var a = 0

		assert.NoError(t, Do(func() error {
			a = 1

			return nil
		}))

		assert.Equal(t, a, 1)
	})

	t.Run("should return error", func(t *testing.T) {
		var err = errors.New("dummy")

		assert.ErrorIs(t, Do(func() error {
			return err
		}), err)
	})
}
