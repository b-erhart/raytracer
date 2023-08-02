package geometry

import (
	"sync"

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

func (r *Raytracer) Render(view *View, canv *canvas.Canvas) {
	origin := view.Eye()
	lineStart := view.BottomLeft()

	var wg sync.WaitGroup

	for j := 0; j < canv.Height(); j++ {
		current := lineStart

		for i := 0; i < canv.Width(); i++ {
			wg.Add(1)

			go func(i, j int, current *Vector) {
				defer wg.Done()

				direction := Sub(current, origin).Normalize()

				ray := Ray{
					Origin:    origin,
					Direction: direction,
					Depth:     0,
				}

				canv.SetColor(i, j, r.Trace(&ray))
			}(i, j, current)
			
			current = Add(current, view.Du())
		}

		lineStart = Add(lineStart, view.Dv())
	}

	wg.Wait()
}

func (r *Raytracer) Trace(ray *Ray) canvas.Color {
	for _, object := range *r.objects {
		if object.HitBy(ray) {
			return canvas.Color{R: 255, G: 0, B: 0}
		}
	}

	return r.background
}
