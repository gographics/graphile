package graphite

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseLineVertex(t *testing.T) {
	w := Wavefront{}
	w.parseLine("v 1.0 2.0 3.0")
	assert.Equal(t, []float32{1.0, 2.0, 3.0}, w.vertex[0])
}

func TestParseVertexNormalLine(t *testing.T) {
	w := Wavefront{}
	w.parseLine("vn 1.0 2.0 3.0")
	assert.Equal(t, []float32{1.0, 2.0, 3.0}, w.vertex_normal[0])
}

func TestParseVertexTextureLine(t *testing.T) {
	w := Wavefront{}
	w.parseLine("vt 1.0 1.0")
	assert.Equal(t, []float32{1.0, 1.0}, w.vertex_texture[0])
}
