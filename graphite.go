package graphite

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Data structures
type Mesh struct {
	Name    string    `json:"name"`
	Vertex  []float32 `json:"vertex"`
	Normal  []float32 `json:"normal"`
	Color   []float32 `json:"color"`
	Texture []float32 `json:"texture"`
	Indices []uint32  `json:"indices"`
}

// File types
type PlainTextFile struct {
	path    string
	Name    string
	Ext     string
	Content []byte
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

func StripNameExt(path string) (name, ext string) {
	path_split := strings.Split(path, "/")
	splited := strings.Split(path_split[len(path_split)-1], ".")
	return splited[0], splited[len(splited)-1]
}

func PlainText(path string) (PlainTextFile, error) {
	file := PlainTextFile{path: path}
	file.Name, file.Ext = StripNameExt(file.path)
	content, err := ioutil.ReadFile(file.path)
	if err != nil {
		return PlainTextFile{}, err
	}
	file.Content = content
	return file, nil
}

type loader interface {
	Build() (Mesh, error)   // parse the file and return the vertex array geometry
	compile() (Mesh, error) // compile GeomeryFile structure into vertey array
	parseLine()             // line by line parser
}

func Geometry(path string) (Mesh, error) {
	switch filepath.Ext(path) {
	case ".obj":
		file := Wavefront{path: path}
		return file.Build()
	case ".ply":
		return Mesh{}, errors.New("Pending Implementation: Polygon(.ply)")
	default:
		return Mesh{}, errors.New("Unable to parse file")
	}
}
