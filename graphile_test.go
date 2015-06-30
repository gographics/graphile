package graphile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathStrip(t *testing.T) {
	name, ext := pathStrip("some/path/cube.obj")
	assert.Equal(t, "cube.obj", name)
	assert.Equal(t, ".obj", ext)
}
