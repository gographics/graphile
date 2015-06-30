package graphile

import (
	"strconv"
	"strings"
)

func (g *GeometryFile) parseLineOBJ(line string) {
	newLine := strings.Replace(strings.TrimSpace(line), "/", " ", -1)
	larray := strings.Split(newLine, " ")

	switch larray[0] {
	case "v":
		x, _ := strconv.ParseFloat(larray[1], 32)
		y, _ := strconv.ParseFloat(larray[2], 32)
		z, _ := strconv.ParseFloat(larray[3], 32)
		g.vertex = append(g.vertex, []float32{float32(x), float32(y), float32(z)})
	case "vt":
		x, _ := strconv.ParseFloat(larray[1], 32)
		y, _ := strconv.ParseFloat(larray[2], 32)
		g.vertexTexture = append(g.vertexTexture, []float32{float32(x), float32(y)})
		g.hasTextures = g.hasTextures || true
	case "vn":
		x, _ := strconv.ParseFloat(larray[1], 32)
		y, _ := strconv.ParseFloat(larray[2], 32)
		z, _ := strconv.ParseFloat(larray[3], 32)
		g.vertexNormal = append(g.vertexNormal, []float32{float32(x), float32(y), float32(z)})
		g.hasNormals = g.hasNormals || true
	case "f":
		buffer := []int32{}
		for _, faceIdx := range larray[1:] {
			v, err := strconv.Atoi(faceIdx)
			if err == nil {
				buffer = append(buffer, int32(v))
			}
		}

		offset := g.offset()
		pivot := offset * 3

		switch len(buffer) / offset {
		case 3:
			g.triangles = append(g.triangles, buffer)
		case 4:
			buffer = append(buffer, buffer[0:offset]...)
			g.triangles = append(g.triangles, buffer[:pivot])        // T1
			g.triangles = append(g.triangles, buffer[pivot-offset:]) // T2
		default:
			return
		}
	default:
		return
	}
}
