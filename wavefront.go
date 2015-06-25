package graphile

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type wavefront GeometryFile

func (w wavefront) Open() (Mesh, error) {
	file, err := os.Open(w.path)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w.parseLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return w.compile()
}

func (w *wavefront) parseLine(line string) {
	newLine := strings.Replace(strings.TrimSpace(line), "/", " ", -1)
	larray := strings.Split(newLine, " ")

	switch larray[0] {
	case "v":
		x, _ := strconv.ParseFloat(larray[1], 32)
		y, _ := strconv.ParseFloat(larray[2], 32)
		z, _ := strconv.ParseFloat(larray[3], 32)
		w.vertex = append(w.vertex, []float32{float32(x), float32(y), float32(z)})
	case "vt":
		x, _ := strconv.ParseFloat(larray[1], 32)
		y, _ := strconv.ParseFloat(larray[2], 32)
		w.vertexTexture = append(w.vertexTexture, []float32{float32(x), float32(y)})
		w.hasTextures = w.hasTextures || true
	case "vn":
		x, _ := strconv.ParseFloat(larray[1], 32)
		y, _ := strconv.ParseFloat(larray[2], 32)
		z, _ := strconv.ParseFloat(larray[3], 32)
		w.vertexNormal = append(w.vertexNormal, []float32{float32(x), float32(y), float32(z)})
		w.hasNormals = w.hasNormals || true
	case "f":
		buffer := []int32{}
		for _, faceIdx := range larray[1:] {
			v, err := strconv.Atoi(faceIdx)
			if err == nil {
				buffer = append(buffer, int32(v))
			}
		}

		offset := w.offset()
		pivot := offset * 3

		switch len(buffer) / offset {
		case 3:
			w.triangles = append(w.triangles, buffer)
		case 4:
			buffer = append(buffer, buffer[0:offset]...)
			w.triangles = append(w.triangles, buffer[:pivot])        // T1
			w.triangles = append(w.triangles, buffer[pivot-offset:]) // T2
		default:
			return
		}
	default:
		return
	}
}

func (w *wavefront) compile() (Mesh, error) {
	var (
		outVertexList  []float32
		outTextureList []float32
		outNormalList  []float32
	)

	offset := w.offset()

	// sorting from faces indexes
	switch offset {
	case 1:
		for _, triangle := range w.triangles {
			for idx := 0; idx < len(triangle); idx++ {
				outVertexList = append(outVertexList, w.vertex[triangle[idx]-1]...)
			}
		}
		if len(outVertexList) != 9*len(w.triangles) {
			return Mesh{}, errors.New("Compilation Error: mismatch length on vertex array")
		}
	case 2:
		if w.hasTextures {
			for _, triangle := range w.triangles {
				for idx := 0; idx < len(triangle); idx += offset {
					outVertexList = append(outVertexList, w.vertex[triangle[idx]-1]...)
					outTextureList = append(outTextureList, w.vertexTexture[triangle[idx+1]-1]...)
				}
			}
			if len(outVertexList)/3 != len(outTextureList)/2 {
				return Mesh{}, errors.New("Compilation Error: mismatch length between vertex and texture arrays")
			}
		}
		if w.hasNormals {
			for _, triangle := range w.triangles {
				for idx := 0; idx < len(triangle); idx += offset {
					outVertexList = append(outVertexList, w.vertex[triangle[idx]-1]...)
					outNormalList = append(outNormalList, w.vertexNormal[triangle[idx+1]-1]...)
				}
			}
			if len(outVertexList) != len(outNormalList) {
				return Mesh{}, errors.New("Compilation Error: mismatch length between vertex and normal arrays")
			}
		}
	case 3:
		for _, triangle := range w.triangles {
			for idx := 0; idx < len(triangle); idx += offset {
				outVertexList = append(outVertexList, w.vertex[triangle[idx]-1]...)
				outTextureList = append(outTextureList, w.vertexTexture[triangle[idx+1]-1]...)
				outNormalList = append(outNormalList, w.vertexNormal[triangle[idx+2]-1]...)
			}
		}
		if len(outVertexList) != len(outNormalList) || len(outVertexList)/3 != len(outTextureList)/2 {
			return Mesh{}, errors.New("Compilation Error: mismatch length vertex arrays")
		}
	default:
		return Mesh{}, errors.New("Offset invalid")
	}

	return Mesh{Name: w.name, Vertex: outVertexList, Normal: outNormalList, Texture: outTextureList}, nil
}

func (w *wavefront) offset() int {
	offset := 1
	if w.hasNormals {
		offset++
	}
	if w.hasTextures {
		offset++
	}
	return offset
}
