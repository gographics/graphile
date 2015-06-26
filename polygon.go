package graphile

import (
	"strconv"
	"strings"
)

type polygon GeometryFile

func (p *polygon) parseLine(line string) {
	if skipLine(line) {
		return
	}
	larray := strings.Split(strings.TrimSpace(line), " ")
	if larray[0] == "element" {
		if larray[1] == "vertex" {
			size, _ := strconv.Atoi(larray[2])
			p.vertexSize = size
		}
		return
	}
	if len(p.vertex) < p.vertexSize {
		x, _ := strconv.ParseFloat(larray[0], 32)
		y, _ := strconv.ParseFloat(larray[1], 32)
		z, _ := strconv.ParseFloat(larray[2], 32)
		p.vertex = append(p.vertex, []float32{float32(x), float32(y), float32(z)})
	} else {
		faces, _ := strconv.Atoi(larray[0])
		a, _ := strconv.Atoi(larray[1])
		b, _ := strconv.Atoi(larray[2])
		c, _ := strconv.Atoi(larray[3])
		p.triangles = append(p.triangles, []int32{int32(a), int32(b), int32(c)})
		if faces == 4 {
			d, _ := strconv.Atoi(larray[4])
			p.triangles = append(p.triangles, []int32{int32(c), int32(d), int32(a)})
		}
	}
}

func skipLine(s string) bool {
	return strings.Contains(s, "ply") || strings.Contains(s, "comment") || strings.Contains(s, "format") || strings.Contains(s, "property") || strings.Contains(s, "end_header")
}
