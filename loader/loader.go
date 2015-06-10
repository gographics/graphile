package loader

import (
	"errors"
	"path/filepath"
)

type GeometryFile struct {
	name           string
	path           string
	vertex         [][]float32
	vertex_texture [][]float32
	vertex_normal  [][]float32
	triangles      [][]int32
}

type Wavefront GeometryFile
type Polygon GeometryFile

type loader interface {
	BuildMesh() (geometry.Mesh, error) // parse the file and return the vertex array geometry
	Compile() (geometry.Mesh, error)   // compile GeomeryFile structure into vertey array
	lineParse()                        // line by line parser
	hasTextures() bool
	hasNormals() bool
}

func Geometry(path string) (geometry.Mesh, error) {
	switch filepath.Ext(path) {
	case ".obj":
		file := Wavefront{path: path}
		return file.BuildMesh()
	case ".ply":
		return geometry.Mesh{}, errors.New("Pending Implementation: Polygon(.ply)")
	default:
		return geometry.Mesh{}, errors.New("Unable to parse file")
	}
}
