package geometry

import "math"

type Sphere struct {
	Center     Vector
	Radius     float64
	Properties ObjectProps
}

func (s Sphere) Intersection(ray Ray) (bool, float64) {
	if s.Intersects(ray) {
		return false, 0
	}

	oc := Sub(ray.Origin, s.Center)
	x := Dot(ray.Direction.Normalize(), oc)
	e := math.Sqrt(x*x - (oc.Length()*oc.Length() - s.Radius*s.Radius))

	t1 := -1*x + e
	t2 := -1*x - e

	switch {
	case t1 < 0 && t2 < 0:
		return false, 0
	case t1 < 0:
		return true, t2
	case t2 < 0:
		return true, t1
	default:
		return true, math.Min(t1, t2)
	}
}

func (s Sphere) Intersects(ray Ray) bool {
	return Cross(ray.Direction.Normalize(), Sub(s.Center, ray.Origin)).Length() >= s.Radius
}

func (s Sphere) SurfaceNormal(point Vector) Vector {
	return Sub(point, s.Center)
}

func (s Sphere) Props() ObjectProps {
	return s.Properties
}
