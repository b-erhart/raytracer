package geometry

import "github.com/b-erhart/raytracer/internal/canvas"

type Scene struct {
	Canvas     canvas.Canvas
	View       View
	Objects    []Object
	Lights     []Light
	Background canvas.Color
	SSAA       bool
}
