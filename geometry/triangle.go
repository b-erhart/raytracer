package geometry

import "math"

type Triangle struct {
	A                Vector
	B                Vector
	C                Vector
	Properties       ObjectProps
	ASurfaceNormal   Vector
	BSurfaceNormal   Vector
	CSurfaceNormal   Vector
	NormalsSet       bool
	edgesCalculated  bool
	edge1            Vector
	edge2            Vector
	extrmsCalculated bool
	extrms           extremes
}

// Calculate intersection between triangle and ray.
// Source: https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
func (triangle *Triangle) Intersection(ray Ray) (bool, float64) {
	if !triangle.edgesCalculated {
		triangle.calculateEdges()
	}

	h := Cross(ray.Direction, triangle.edge2)
	a := Dot(triangle.edge1, h)

	if a > -Epsilon && a < Epsilon {
		return false, 0
	}

	f := 1 / a
	s := Sub(ray.Origin, triangle.A)
	u := f * Dot(s, h)

	if u < 0 || u > 1 {
		return false, 0
	}

	q := Cross(s, triangle.edge1)
	v := f * Dot(ray.Direction, q)

	if v < 0 || u+v > 1 {
		return false, 0
	}

	t := f * Dot(triangle.edge2, q)

	if t > Epsilon {
		return true, t
	}

	return false, 0
}

func (t *Triangle) calculateEdges() {
	t.edge1 = Sub(t.B, t.A)
	t.edge2 = Sub(t.C, t.A)
	t.edgesCalculated = true
}

func (t *Triangle) SurfaceNormal(point Vector) Vector {
	bary := t.bary(point)

	interpolated := Add(Add(Sprod(t.ASurfaceNormal, bary.X), Sprod(t.BSurfaceNormal, bary.Y)), Sprod(t.CSurfaceNormal, bary.Z))
	return interpolated.Normalize()
}

func (t *Triangle) TriangleNormal() Vector {
	if !t.edgesCalculated {
		t.calculateEdges()
	}

	return Cross(t.edge1, t.edge2)
}

func (t *Triangle) Props() ObjectProps {
	return t.Properties
}

func (t *Triangle) extremes() extremes {
	if !t.extrmsCalculated {
		t.calculateExtremes()
	}

	return t.extrms
}

func (t *Triangle) bary(p Vector) Vector {
	var bary Vector

	normal := t.TriangleNormal()

	areaABC := Dot(normal, Cross(Sub(t.B, t.A), Sub(t.C, t.A)))
	areaPBC := Dot(normal, Cross(Sub(t.B, p), Sub(t.C, p)))
	areaPCA := Dot(normal, Cross(Sub(t.C, p), Sub(t.A, p)))

	bary.X = areaPBC / areaABC     // alpha
	bary.Y = areaPCA / areaABC     // beta
	bary.Z = 1.0 - bary.X - bary.Y // gamma

	return bary
}

func (t *Triangle) calculateExtremes() {
	t.extrms = extremes{
		minX: math.Min(t.A.X, math.Min(t.B.X, t.C.X)),
		minY: math.Min(t.A.Y, math.Min(t.B.Y, t.C.Y)),
		minZ: math.Min(t.A.Z, math.Min(t.B.Z, t.C.Z)),
		maxX: math.Max(t.A.X, math.Max(t.B.X, t.C.X)),
		maxY: math.Max(t.A.Y, math.Max(t.B.Y, t.C.Y)),
		maxZ: math.Max(t.A.Z, math.Max(t.B.Z, t.C.Z)),
	}
	t.extrmsCalculated = true
}
