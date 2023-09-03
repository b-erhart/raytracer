package geometry

import "math"

type extremes struct {
	minX float64
	minY float64
	minZ float64
	maxX float64
	maxY float64
	maxZ float64
}

func merge(a, b extremes) extremes {
	return extremes{
		minX: math.Min(a.minX, b.minX),
		minY: math.Min(a.minY, b.minY),
		minZ: math.Min(a.minZ, b.minZ),
		maxX: math.Max(a.maxX, b.maxX),
		maxY: math.Max(a.maxY, b.maxY),
		maxZ: math.Max(a.maxZ, b.maxZ),
	}
}
