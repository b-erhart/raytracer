package geometry

import (
	"github.com/b-erhart/raytracer/canvas"
)

type Light struct {
	Direction Vector
	Color     canvas.Color
}
