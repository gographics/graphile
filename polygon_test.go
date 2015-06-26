package graphile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseElementVertex(t *testing.T) {
	p := polygon{}
	p.parseLine("element vertex 6906")
	assert.Equal(t, 6906, p.vertexSize)
}

func TestParseVector(t *testing.T) {
	p := polygon{}
	p.parseLine("element vertex 1")
	p.parseLine("1.56444 -0.204025 0.346805")
	assert.Equal(t, 1, p.vertexSize)
	assert.Equal(t, []float32{1.56444, -0.204025, 0.346805}, p.vertex[len(p.vertex)-1])
}

func TestParseVectorScientificPrecition(t *testing.T) {
	p := polygon{}
	p.parseLine("element vertex 1")
	p.parseLine("0.167500E-02 0.505000E-02 0.000000E+00")
	assert.Equal(t, 1, p.vertexSize)
	assert.Equal(t, []float32{0.001675, 0.00505, 0.0}, p.vertex[len(p.vertex)-1])
}

func TestParseFaceTriangles(t *testing.T) {
	p := polygon{}
	p.parseLine("element vertex 1")
	p.parseLine("1.0 0.0 0.0")
	p.parseLine("3 1 2 3")
	assert.Equal(t, []int32{1, 2, 3}, p.triangles[0])
}

func TestParseFaceCuads(t *testing.T) {
	p := polygon{}
	p.parseLine("element vertex 1")
	p.parseLine("1.0 0.0 0.0")
	p.parseLine("4 1 2 3 4")
	assert.Equal(t, []int32{1, 2, 3}, p.triangles[0])
	assert.Equal(t, []int32{3, 4, 1}, p.triangles[1])
}
