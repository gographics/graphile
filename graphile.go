package graphile

import (
	"bufio"
	"errors"
	"os"
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
	format        string
	vertex        [][]float32
	vertexTexture [][]float32
	vertexNormal  [][]float32
	triangles     [][]int32

	hasNormals  bool
	hasTextures bool

	vertexSize int
}

func Open(filepath string) (Mesh, error) {
	name, ext := pathStrip(filepath)
	g := GeometryFile{name: name}
	switch ext {
	case ".obj":
		g.format = "wavefront"
	case ".ply":
		g.format = "polygon"
	default:
		return Mesh{}, errors.New("File format is not implemented")
	}
	file, err := os.Open(g.path)
	if err != nil {
		return Mesh{}, err
	}

	scanner := bufio.NewScanner(file)
	switch g.format {
	case "wavefront":
		for scanner.Scan() {
			g.parseLineOBJ(scanner.Text())
		}
	case "polygon":
		for scanner.Scan() {
			g.parseLinePLY(scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return Mesh{}, err
	}

	return g.compile()
}

func (g *GeometryFile) compile() (Mesh, error) {
	var (
		outVertexList  []float32
		outTextureList []float32
		outNormalList  []float32
	)

	offset := g.offset()

	// sorting from faces indexes
	switch offset {
	case 1:
		for _, triangle := range g.triangles {
			for idx := 0; idx < len(triangle); idx++ {
				outVertexList = append(outVertexList, g.vertex[triangle[idx]-1]...)
			}
		}
		if len(outVertexList) != 9*len(g.triangles) {
			return Mesh{}, errors.New("Compilation Error: mismatch length on vertex array")
		}
	case 2:
		if g.hasTextures {
			for _, triangle := range g.triangles {
				for idx := 0; idx < len(triangle); idx += offset {
					outVertexList = append(outVertexList, g.vertex[triangle[idx]-1]...)
					outTextureList = append(outTextureList, g.vertexTexture[triangle[idx+1]-1]...)
				}
			}
			if len(outVertexList)/3 != len(outTextureList)/2 {
				return Mesh{}, errors.New("Compilation Error: mismatch length between vertex and texture arrays")
			}
		}
		if g.hasNormals {
			for _, triangle := range g.triangles {
				for idx := 0; idx < len(triangle); idx += offset {
					outVertexList = append(outVertexList, g.vertex[triangle[idx]-1]...)
					outNormalList = append(outNormalList, g.vertexNormal[triangle[idx+1]-1]...)
				}
			}
			if len(outVertexList) != len(outNormalList) {
				return Mesh{}, errors.New("Compilation Error: mismatch length between vertex and normal arrays")
			}
		}
	case 3:
		for _, triangle := range g.triangles {
			for idx := 0; idx < len(triangle); idx += offset {
				outVertexList = append(outVertexList, g.vertex[triangle[idx]-1]...)
				outTextureList = append(outTextureList, g.vertexTexture[triangle[idx+1]-1]...)
				outNormalList = append(outNormalList, g.vertexNormal[triangle[idx+2]-1]...)
			}
		}
		if len(outVertexList) != len(outNormalList) || len(outVertexList)/3 != len(outTextureList)/2 {
			return Mesh{}, errors.New("Compilation Error: mismatch length vertex arrays")
		}
	default:
		return Mesh{}, errors.New("Offset invalid")
	}

	return Mesh{Name: g.name, Vertex: outVertexList, Normal: outNormalList, Texture: outTextureList}, nil
}

func (g *GeometryFile) offset() int {
	offset := 1
	if g.hasNormals {
		offset++
	}
	if g.hasTextures {
		offset++
	}
	return offset
}

func pathStrip(path string) (name, ext string) {
	_, name = filepath.Split(path)
	ext = filepath.Ext(path)
	return name, ext
}
