package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath2URI(t *testing.T) {
	uri1 := Path2URI("/dir", "subdir/filename with spaces")
	assert.Equal(t, uri1, "file:///dir/subdir/filename%20with%20spaces")
}
