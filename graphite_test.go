package graphite

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStripNameExt(t *testing.T) {
	name, ext := StripNameExt("some/path/cube.obj")
	assert.Equal(t, "cube", name)
	assert.Equal(t, "obj", ext)
}
