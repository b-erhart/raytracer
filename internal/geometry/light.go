package geometry

import (
	"github.com/b-erhart/raytracer/internal/canvas"
)

type Light struct {
	Direction Vector
	Color     canvas.Color
}
