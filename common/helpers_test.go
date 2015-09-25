package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileSizeFromPath(t *testing.T) {

	tests := []struct {
		path   string
		width  int
		height int
	}{
		{"media/catalog/product/detail/cotton-11-14315-001-h.jpg", 0, 0},
		{"catalog/product/cache/2/small_image/218x258/9df78eab33525d08d6e5fb8d27136e95/0/3/03793_224_v.jpg", 218, 258},
		{"small_image/138x165/9df78eab33525d08d6e5fb8d27136e95/detail/favourites-29598-096-v.jpg", 138, 165},
		{"small_image/135x/9df78eab33525d08d6e5fb8d27136e95/detail/favourites-29598-096-v.jpg", 135, 135},
	}

	for i, test := range tests {
		w, h := FileSizeFromPath(test.path)
		assert.Equal(t, test.width, w, "Index %d", i)
		assert.Equal(t, test.height, h, "Index %d", i)
	}
}
