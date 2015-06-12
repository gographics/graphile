package graphite

import (
	"path/filepath"
)

type Mesh struct {
	Name    string    `json:"name"`
	Vertex  []float32 `json:"vertex"`
	Normal  []float32 `json:"normal"`
	Color   []float32 `json:"color"`
	Texture []float32 `json:"texture"`
	Indices []uint32  `json:"indices"`
}

type GeometryFile struct {
	name           string
	path           string
	vertex         [][]float32
	vertex_texture [][]float32
	vertex_normal  [][]float32
	triangles      [][]int32

	hasNormals  bool
	hasTextures bool
}

func PathStrip(path string) (name, ext string) {
	_, name = filepath.Split(path)
	ext = filepath.Ext(path)
	return name, ext
}

type graphite interface {
	Open() (Mesh, error)
	Save(path string) error
}
