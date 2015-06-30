package graphile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseOBJVertex(t *testing.T) {
	g := GeometryFile{format: "wavefront"}
	g.parseLineOBJ("v 1.0 2.0 3.0")
	assert.Equal(t, []float32{1.0, 2.0, 3.0}, g.vertex[0])
	assert.False(t, g.hasNormals)
	assert.False(t, g.hasTextures)
}

func TestParseOBJVertexNormal(t *testing.T) {
	g := GeometryFile{format: "wavefront"}
	g.parseLineOBJ("vn 1.0 0.0 0.0")
	assert.Equal(t, []float32{1.0, 0.0, 0.0}, g.vertexNormal[0])
	assert.True(t, g.hasNormals)
}

func TestParseOBJVertexTexture(t *testing.T) {
	g := GeometryFile{format: "wavefront"}
	g.parseLineOBJ("vt 1.0 1.0")
	assert.Equal(t, []float32{1.0, 1.0}, g.vertexTexture[0])
	assert.True(t, g.hasTextures)
}

func TestParseOBJFace3(t *testing.T) {
	g := GeometryFile{format: "wavefront"}
	g.parseLineOBJ("f 1 2 3")
	assert.Equal(t, []int32{1, 2, 3}, g.triangles[0])
}

func TestParseOBJFace3Normals(t *testing.T) {
	g := GeometryFile{format: "wavefront", hasNormals: true}
	g.parseLineOBJ("f 0//1 2//3 4//5")
	assert.Equal(t, []int32{0, 1, 2, 3, 4, 5}, g.triangles[0])
}

func TestParseOBJFace3NormalsTexures(t *testing.T) {
	g := GeometryFile{format: "wavefront", hasNormals: true, hasTextures: true}
	g.parseLineOBJ("f 0/1/2 3/4/5 6/7/8")
	assert.Equal(t, []int32{0, 1, 2, 3, 4, 5, 6, 7, 8}, g.triangles[0])
}

func TestParseOBJFace4(t *testing.T) {
	g := GeometryFile{}
	g.parseLineOBJ("f 0 1 2 3")
	assert.Equal(t, []int32{0, 1, 2}, g.triangles[0])
	assert.Equal(t, []int32{2, 3, 0}, g.triangles[1])
}

func TestParseOBJFace4Normals(t *testing.T) {
	g := GeometryFile{format: "wavefront", hasNormals: true}
	g.parseLineOBJ("f 0//1 2//3 4//5 6//7")
	assert.Equal(t, []int32{0, 1, 2, 3, 4, 5}, g.triangles[0])
	assert.Equal(t, []int32{4, 5, 6, 7, 0, 1}, g.triangles[1])
}

func TestParseOBJFace4NormalsTexures(t *testing.T) {
	g := GeometryFile{format: "wavefront", hasNormals: true, hasTextures: true}
	g.parseLineOBJ("f 0/1/2 3/4/5 6/7/8 9/10/11")
	assert.Equal(t, []int32{0, 1, 2, 3, 4, 5, 6, 7, 8}, g.triangles[0])
	assert.Equal(t, []int32{6, 7, 8, 9, 10, 11, 0, 1, 2}, g.triangles[1])
}

func TestCompileOBJ(t *testing.T) {
	g := GeometryFile{name: "triangle", format: "wavefront"}
	g.parseLineOBJ("v 0 0 0")
	g.parseLineOBJ("v 1 0 0")
	g.parseLineOBJ("v 1 1 0")
	g.parseLineOBJ("f 1 2 3")
	mesh, err := g.compile()
	assert.Nil(t, err)
	assert.Equal(t, []float32{0, 0, 0, 1, 0, 0, 1, 1, 0}, mesh.Vertex)
}

func TestCompileOBJNormal(t *testing.T) {
	g := GeometryFile{name: "triangle", format: "wavefront"}
	g.parseLineOBJ("v 0 0 0")
	g.parseLineOBJ("v 1 0 0")
	g.parseLineOBJ("v 1 1 0")
	g.parseLineOBJ("vn 0 0 1")
	g.parseLineOBJ("f 1//1 2//1 3//1")
	mesh, err := g.compile()
	assert.Nil(t, err)
	assert.Equal(t, []float32{0, 0, 0, 1, 0, 0, 1, 1, 0}, mesh.Vertex)
	assert.Equal(t, []float32{0, 0, 1, 0, 0, 1, 0, 0, 1}, mesh.Normal)
}

func TestCompileOBJTexture(t *testing.T) {
	g := GeometryFile{name: "triangle", format: "wavefront"}
	g.parseLineOBJ("v 0 0 0")
	g.parseLineOBJ("v 1 0 0")
	g.parseLineOBJ("v 1 1 0")
	g.parseLineOBJ("vt 0 1")
	g.parseLineOBJ("vt 1 1")
	g.parseLineOBJ("vt 1 0")
	g.parseLineOBJ("f 1/1 2/2 3/3")
	mesh, err := g.compile()
	assert.Nil(t, err)
	assert.Equal(t, []float32{0, 0, 0, 1, 0, 0, 1, 1, 0}, mesh.Vertex)
	assert.Equal(t, []float32{0, 1, 1, 1, 1, 0}, mesh.Texture)
}

func TestCompileOBJNormalTexture(t *testing.T) {
	g := GeometryFile{name: "triangle", format: "wavefront"}
	g.parseLineOBJ("v 0 0 0")
	g.parseLineOBJ("v 1 0 0")
	g.parseLineOBJ("v 1 1 0")
	g.parseLineOBJ("vt 0 0")
	g.parseLineOBJ("vt 1 0")
	g.parseLineOBJ("vt 1 1")
	g.parseLineOBJ("vn 0 0 1")
	g.parseLineOBJ("f 1/1/1 2/2/1 3/3/1")
	mesh, err := g.compile()
	assert.Nil(t, err)
	assert.Equal(t, []float32{0, 0, 0, 1, 0, 0, 1, 1, 0}, mesh.Vertex)
	assert.Equal(t, []float32{0, 0, 1, 0, 0, 1, 0, 0, 1}, mesh.Normal)
	assert.Equal(t, []float32{0, 0, 1, 0, 1, 1}, mesh.Texture)
}
