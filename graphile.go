package graphile

import (
	"path/filepath"
)

// Mesh geometry data representation
type Mesh struct {
	Name    string    `json:"name"`
	Vertex  []float32 `json:"vertex"`
	Normal  []float32 `json:"normal"`
	Color   []float32 `json:"color"`
	Texture []float32 `json:"texture"`
	Indices []uint32  `json:"indices"`
}

// GeometryFile store raw geometry before to compile into a mesh structure
type GeometryFile struct {
	name          string
	path          string
	vertex        [][]float32
	vertexTexture [][]float32
	vertexNormal  [][]float32
	triangles     [][]int32

	hasNormals  bool
	hasTextures bool

	vertexSize int
}

// PathStrip strip name and extension from file
func PathStrip(path string) (name, ext string) {
	_, name = filepath.Split(path)
	ext = filepath.Ext(path)
	return name, ext
}

type graphite interface {
	Open() (Mesh, error)
	Save(path string) error
}
