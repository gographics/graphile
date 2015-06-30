package graphile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsePLYElementVertex(t *testing.T) {
	g := GeometryFile{format: "polygon"}
	g.parseLinePLY("element vertex 6906")
	assert.Equal(t, 6906, g.vertexSize)
}

func TestParsePLYVector(t *testing.T) {
	g := GeometryFile{format: "polygon"}
	g.parseLinePLY("element vertex 1")
	g.parseLinePLY("1.56444 -0.204025 0.346805")
	assert.Equal(t, 1, g.vertexSize)
	assert.Equal(t, []float32{1.56444, -0.204025, 0.346805}, g.vertex[len(g.vertex)-1])
}

func TestParsePLYVectorScientificPrecition(t *testing.T) {
	g := GeometryFile{format: "polygon"}
	g.parseLinePLY("element vertex 1")
	g.parseLinePLY("0.167500E-02 0.505000E-02 0.000000E+00")
	assert.Equal(t, 1, g.vertexSize)
	assert.Equal(t, []float32{0.001675, 0.00505, 0.0}, g.vertex[len(g.vertex)-1])
}

func TestParsePLYFaceTriangles(t *testing.T) {
	g := GeometryFile{format: "polygon"}
	g.parseLinePLY("element vertex 1")
	g.parseLinePLY("1.0 0.0 0.0")
	g.parseLinePLY("3 1 2 3")
	assert.Equal(t, []int32{1, 2, 3}, g.triangles[0])
}

func TestParsePLYFaceCuads(t *testing.T) {
	g := GeometryFile{format: "polygon"}
	g.parseLinePLY("element vertex 1")
	g.parseLinePLY("1.0 0.0 0.0")
	g.parseLinePLY("4 1 2 3 4")
	assert.Equal(t, []int32{1, 2, 3}, g.triangles[0])
	assert.Equal(t, []int32{3, 4, 1}, g.triangles[1])
}

func TestPolygonCompile(t *testing.T) {
	g := GeometryFile{format: "polygon"}
	g.parseLinePLY("element vertex 3")
	g.parseLinePLY("0 0 0")
	g.parseLinePLY("1 0 0")
	g.parseLinePLY("1 1 0")
	g.parseLinePLY("3 1 2 3")
	mesh, err := g.compile()
	assert.Nil(t, err)
	assert.Equal(t, []float32{0, 0, 0, 1, 0, 0, 1, 1, 0}, mesh.Vertex)
}
