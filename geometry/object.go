package geometry

import "github.com/b-erhart/raytracer/canvas"

type Object interface {
	Intersection(ray Ray) (bool, float64)
	SurfaceNormal(point Vector) Vector
	Color() canvas.Color
	Reflectivity() float64
	Mirror() float64
}
