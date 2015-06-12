package graphite

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseVertex(t *testing.T) {
	w := Wavefront{}
	w.parseLine("v 1.0 2.0 3.0")
	assert.Equal(t, []float32{1.0, 2.0, 3.0}, w.vertex[0])
	assert.False(t, w.hasNormals)
	assert.False(t, w.hasTextures)
}

func TestParseVertexNormal(t *testing.T) {
	w := Wavefront{}
	w.parseLine("vn 1.0 0.0 0.0")
	assert.Equal(t, []float32{1.0, 0.0, 0.0}, w.vertex_normal[0])
	assert.True(t, w.hasNormals)
}

func TestParseVertexTexture(t *testing.T) {
	w := Wavefront{}
	w.parseLine("vt 1.0 1.0")
	assert.Equal(t, []float32{1.0, 1.0}, w.vertex_texture[0])
	assert.True(t, w.hasTextures)
}

func TestParseFace3(t *testing.T) {
	w := Wavefront{}
	w.parseLine("f 1 2 3")
	assert.Equal(t, []int32{1, 2, 3}, w.triangles[0])
}

func TestParseFace3Normals(t *testing.T) {
	w := Wavefront{hasNormals: true}
	w.parseLine("f 0//1 2//3 4//5")
	assert.Equal(t, []int32{0, 1, 2, 3, 4, 5}, w.triangles[0])
}

func TestParseFace3NormalsTexures(t *testing.T) {
	w := Wavefront{hasNormals: true, hasTextures: true}
	w.parseLine("f 0/1/2 3/4/5 6/7/8")
	assert.Equal(t, []int32{0, 1, 2, 3, 4, 5, 6, 7, 8}, w.triangles[0])
}

func TestParseFace4(t *testing.T) {
	w := Wavefront{}
	w.parseLine("f 0 1 2 3")
	assert.Equal(t, []int32{0, 1, 2}, w.triangles[0])
	assert.Equal(t, []int32{2, 3, 0}, w.triangles[1])
}

func TestParseFace4Normals(t *testing.T) {
	w := Wavefront{hasNormals: true}
	w.parseLine("f 0//1 2//3 4//5 6//7")
	assert.Equal(t, []int32{0, 1, 2, 3, 4, 5}, w.triangles[0])
	assert.Equal(t, []int32{4, 5, 6, 7, 0, 1}, w.triangles[1])
}

func TestParseFace4NormalsTexures(t *testing.T) {
	w := Wavefront{hasNormals: true, hasTextures: true}
	w.parseLine("f 0/1/2 3/4/5 6/7/8 9/10/11")
	assert.Equal(t, []int32{0, 1, 2, 3, 4, 5, 6, 7, 8}, w.triangles[0])
	assert.Equal(t, []int32{6, 7, 8, 9, 10, 11, 0, 1, 2}, w.triangles[1])
}

func TestCompile(t *testing.T) {
	w := Wavefront{name: "triangle"}
	w.parseLine("v 0 0 0")
	w.parseLine("v 1 0 0")
	w.parseLine("v 1 1 0")
	w.parseLine("f 1 2 3")
	mesh, err := w.compile()
	assert.Nil(t, err)
	assert.Equal(t, []float32{0, 0, 0, 1, 0, 0, 1, 1, 0}, mesh.Vertex)
}

func TestCompileNormal(t *testing.T) {
	w := Wavefront{name: "triangle"}
	w.parseLine("v 0 0 0")
	w.parseLine("v 1 0 0")
	w.parseLine("v 1 1 0")
	w.parseLine("vn 0 0 1")
	w.parseLine("f 1//1 2//1 3//1")
	mesh, err := w.compile()
	assert.Nil(t, err)
	assert.Equal(t, []float32{0, 0, 0, 1, 0, 0, 1, 1, 0}, mesh.Vertex)
	assert.Equal(t, []float32{0, 0, 1, 0, 0, 1, 0, 0, 1}, mesh.Normal)
}

func TestCompileTexture(t *testing.T) {
	w := Wavefront{name: "triangle"}
	w.parseLine("v 0 0 0")
	w.parseLine("v 1 0 0")
	w.parseLine("v 1 1 0")
	w.parseLine("vt 0 1")
	w.parseLine("vt 1 1")
	w.parseLine("vt 1 0")
	w.parseLine("f 1/1 2/2 3/3")
	mesh, err := w.compile()
	assert.Nil(t, err)
	assert.Equal(t, []float32{0, 0, 0, 1, 0, 0, 1, 1, 0}, mesh.Vertex)
	assert.Equal(t, []float32{0, 1, 1, 1, 1, 0}, mesh.Texture)
}

func TestCompileNormalTexture(t *testing.T) {
	w := Wavefront{name: "triangle"}
	w.parseLine("v 0 0 0")
	w.parseLine("v 1 0 0")
	w.parseLine("v 1 1 0")
	w.parseLine("vt 0 0")
	w.parseLine("vt 1 0")
	w.parseLine("vt 1 1")
	w.parseLine("vn 0 0 1")
	w.parseLine("f 1/1/1 2/2/1 3/3/1")
	mesh, err := w.compile()
	assert.Nil(t, err)
	assert.Equal(t, []float32{0, 0, 0, 1, 0, 0, 1, 1, 0}, mesh.Vertex)
	assert.Equal(t, []float32{0, 0, 1, 0, 0, 1, 0, 0, 1}, mesh.Normal)
	assert.Equal(t, []float32{0, 0, 1, 0, 1, 1}, mesh.Texture)
}
