package geometry

import (
	"github.com/b-erhart/raytracer/canvas"
)

type Raytracer struct {
	objects    *[]Object
	lights     *[]Light
	background canvas.Color
	ambience   canvas.Color
}

func NewRaytracer(objects *[]Object, lights *[]Light, background, ambience canvas.Color) *Raytracer {
	return &Raytracer{objects, lights, background, ambience}
}

func (r *Raytracer) Render(view *View, canvas *canvas.Canvas) {
	origin := view.Eye()
	lineStart := view.BottomLeft()

	for j := 0; j < canvas.Height(); j++ {
		current := lineStart

		for i := 0; i < canvas.Width(); i++ {
			direction := Sub(current, origin).Normalize()

			ray := Ray{
				Origin:    origin,
				Direction: direction,
				Depth:     0,
			}

			canvas.SetColor(i, j, r.Trace(&ray))

			current = Add(current, view.Du())
		}

		lineStart = Add(lineStart, view.Dv())
	}
}

func (r *Raytracer) Trace(ray *Ray) canvas.Color {
	for _, object := range *r.objects {
		if object.HitBy(ray) {
			return canvas.Color{R: 255, G: 0, B: 0}
		}
	}

	return r.background
}
