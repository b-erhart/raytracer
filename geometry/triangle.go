package geometry

type Triangle struct {
	A          Vector
	B          Vector
	C          Vector
	Properties ObjectProps
	edges bool
	edge1      Vector
	edge2      Vector
}

// Calculate intersection between triangle and ray.
// Source: https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
func (triangle *Triangle) Intersection(ray Ray) (bool, float64) {
	if !triangle.edges {
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
	t.edges = true
}

func (t *Triangle) SurfaceNormal(point Vector) Vector {
	edge1 := Sub(t.B, t.A)
	edge2 := Sub(t.C, t.A)

	return Cross(edge1, edge2)
}

func (t *Triangle) Props() ObjectProps {
	return t.Properties
}
