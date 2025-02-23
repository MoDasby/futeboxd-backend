package pagination

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	t.Run("should return default(50, 1) for invalid values", func(t *testing.T) {
		page, err := WithPage(context.Background(), "invalid", "invalid")

		assert.NoError(t, err)
		assert.Equal(t, 50, page.Size)
		assert.Equal(t, 1, page.Index)

		page, err = WithPage(context.Background(), "-13", "-10")

		assert.NoError(t, err)
		assert.Equal(t, 50, page.Size)
		assert.Equal(t, 1, page.Index)
	})

	t.Run("should return a correct page", func(t *testing.T) {
		page, err := WithPage(context.Background(), "30", "3")

		assert.NoError(t, err)
		assert.Equal(t, 30, page.Size)
		assert.Equal(t, 3, page.Index)
	})

	t.Run("should return error when creating a page bigger than 100", func(t *testing.T) {
		page, err := WithPage(context.Background(), "300", "3")

		assert.Error(t, err)
		assert.Nil(t, page)
	})
}
