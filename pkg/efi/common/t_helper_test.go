package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	i := New[int]()
	assert.Equal(t, i, 0)

	pi := New[*int]()
	require.NotNil(t, pi)
	assert.Equal(t, *pi, 0)
}

func TestNil(t *testing.T) {
	i := Nil[int]()
	assert.Equal(t, i, 0)

	pi := Nil[*int]()
	assert.Nil(t, pi)
}
