package geometry

type Object interface {
	Intersection(ray Ray) (bool, float64)
	SurfaceNormal(point Vector) Vector
}
