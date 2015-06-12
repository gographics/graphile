package graphite

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Wavefront GeometryFile

func (w Wavefront) Open() (Mesh, error) {
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

func (w *Wavefront) parseLine(line string) {
	new_line := strings.Replace(strings.TrimSpace(line), "/", " ", -1)
	larray := strings.Split(new_line, " ")

	switch larray[0] {
	case "v":
		x, _ := strconv.ParseFloat(larray[1], 32)
		y, _ := strconv.ParseFloat(larray[2], 32)
		z, _ := strconv.ParseFloat(larray[3], 32)
		w.vertex = append(w.vertex, []float32{float32(x), float32(y), float32(z)})
	case "vt":
		x, _ := strconv.ParseFloat(larray[1], 32)
		y, _ := strconv.ParseFloat(larray[2], 32)
		w.vertex_texture = append(w.vertex_texture, []float32{float32(x), float32(y)})
		w.hasTextures = w.hasTextures || true
	case "vn":
		x, _ := strconv.ParseFloat(larray[1], 32)
		y, _ := strconv.ParseFloat(larray[2], 32)
		z, _ := strconv.ParseFloat(larray[3], 32)
		w.vertex_normal = append(w.vertex_normal, []float32{float32(x), float32(y), float32(z)})
		w.hasNormals = w.hasNormals || true
	case "f":
		buffer := []int32{}
		for _, face_index := range larray[1:] {
			v, err := strconv.Atoi(face_index)
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

func (w *Wavefront) compile() (Mesh, error) {
	var (
		out_vertex_list  []float32
		out_texture_list []float32
		out_normal_list  []float32
	)

	offset := w.offset()

	// sorting from faces indexes
	switch offset {
	case 1:
		for _, triangle := range w.triangles {
			for idx := 0; idx < len(triangle); idx++ {
				out_vertex_list = append(out_vertex_list, w.vertex[triangle[idx]-1]...)
			}
		}
		if len(out_vertex_list) != 9*len(w.triangles) {
			return Mesh{}, errors.New("Compilation Error: mismatch length on vertex array")
		}
	case 2:
		if w.hasTextures {
			for _, triangle := range w.triangles {
				for idx := 0; idx < len(triangle); idx += offset {
					out_vertex_list = append(out_vertex_list, w.vertex[triangle[idx]-1]...)
					out_texture_list = append(out_texture_list, w.vertex_texture[triangle[idx+1]-1]...)
				}
			}
			if len(out_vertex_list)/3 != len(out_texture_list)/2 {
				return Mesh{}, errors.New("Compilation Error: mismatch length between vertex and texture arrays")
			}
		}
		if w.hasNormals {
			for _, triangle := range w.triangles {
				for idx := 0; idx < len(triangle); idx += offset {
					out_vertex_list = append(out_vertex_list, w.vertex[triangle[idx]-1]...)
					out_normal_list = append(out_normal_list, w.vertex_normal[triangle[idx+1]-1]...)
				}
			}
			if len(out_vertex_list) != len(out_normal_list) {
				return Mesh{}, errors.New("Compilation Error: mismatch length between vertex and normal arrays")
			}
		}
	case 3:
		for _, triangle := range w.triangles {
			for idx := 0; idx < len(triangle); idx += offset {
				out_vertex_list = append(out_vertex_list, w.vertex[triangle[idx]-1]...)
				out_texture_list = append(out_texture_list, w.vertex_texture[triangle[idx+1]-1]...)
				out_normal_list = append(out_normal_list, w.vertex_normal[triangle[idx+2]-1]...)
			}
		}
		if len(out_vertex_list) != len(out_normal_list) || len(out_vertex_list)/3 != len(out_texture_list)/2 {
			return Mesh{}, errors.New("Compilation Error: mismatch length vertex arrays")
		}
	default:
		return Mesh{}, errors.New("Offset invalid")
	}

	return Mesh{Name: w.name, Vertex: out_vertex_list, Normal: out_normal_list, Texture: out_texture_list}, nil
}

func (w *Wavefront) offset() int {
	var offset int = 1
	if w.hasNormals {
		offset++
	}
	if w.hasTextures {
		offset++
	}
	return offset
}
