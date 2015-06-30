package graphile

import (
	"strconv"
	"strings"
)

func (g *GeometryFile) parseLinePLY(line string) {
	if skipLine(line) {
		return
	}
	larray := strings.Split(strings.TrimSpace(line), " ")
	if larray[0] == "element" {
		if larray[1] == "vertex" {
			size, _ := strconv.Atoi(larray[2])
			g.vertexSize = size
		}
		return
	}
	if len(g.vertex) < g.vertexSize {
		x, _ := strconv.ParseFloat(larray[0], 32)
		y, _ := strconv.ParseFloat(larray[1], 32)
		z, _ := strconv.ParseFloat(larray[2], 32)
		g.vertex = append(g.vertex, []float32{float32(x), float32(y), float32(z)})
	} else {
		faces, _ := strconv.Atoi(larray[0])
		a, _ := strconv.Atoi(larray[1])
		b, _ := strconv.Atoi(larray[2])
		c, _ := strconv.Atoi(larray[3])
		g.triangles = append(g.triangles, []int32{int32(a), int32(b), int32(c)})
		if faces == 4 {
			d, _ := strconv.Atoi(larray[4])
			g.triangles = append(g.triangles, []int32{int32(c), int32(d), int32(a)})
		}
	}
}

func skipLine(s string) bool {
	return strings.Contains(s, "ply") || strings.Contains(s, "comment") || strings.Contains(s, "format") || strings.Contains(s, "property") || strings.Contains(s, "end_header")
}
