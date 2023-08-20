package geometry

import (
	"math"
	"sync"

	"github.com/b-erhart/raytracer/canvas"
)

type Raytracer struct {
	objects    []Object
	lights     *[]Light
	background canvas.Color
	ambience   canvas.Color
}

func NewRaytracer(objects []Object, lights *[]Light, background, ambience canvas.Color) *Raytracer {
	return &Raytracer{objects, lights, background, ambience}
}

func (r *Raytracer) Render(view View, canv *canvas.Canvas) {
	origin := view.Eye()
	lineStart := view.BottomLeft()

	var wg sync.WaitGroup

	for j := 0; j < canv.Height(); j++ {
		current := lineStart

		for i := 0; i < canv.Width(); i++ {
			wg.Add(1)

			go func(i, j int, current Vector) {
				defer wg.Done()

				direction := Sub(current, origin).Normalize()

				ray := Ray{
					Origin:    origin,
					Direction: direction,
					Depth:     0,
				}

				canv.SetColor(i, j, r.Trace(ray))
			}(i, j, current)

			current = Add(current, view.Du())
		}

		lineStart = Add(lineStart, view.Dv())
	}

	wg.Wait()
}

func (r *Raytracer) Trace(ray Ray) canvas.Color {
	if ray.Depth >= 10 {
		return canvas.Color{}
	}

	var closestObj Object
	var tMin float64

	for _, object := range r.objects {
		intersects, t := object.Intersection(ray)

		if intersects && t >= 0.00001 && (closestObj == nil || t < tMin) {
			closestObj = object
			tMin = t
		}
	}

	if closestObj == nil && ray.Depth == 0 {
		return r.background
	} else if closestObj == nil {
		return canvas.Color{}
	} else if closestObj.Reflectivity() <= 0 {
		return closestObj.Color()
	}

	color := closestObj.Color()
	point := ray.At(tMin)
	normal := closestObj.SurfaceNormal(point)
	reflect := Sub(ray.Direction, Sprod(Sprod(normal, Dot(normal, ray.Direction)), 2))
	reflectedRay := Ray{
		Origin:    point,
		Direction: reflect,
		Depth:     ray.Depth + 1,
	}

	Lights:
	for _, light := range *r.lights {
		towardsLight := Sprod(light.Direction, -1).Normalize()
		rayToLight := Ray{
			Origin:    point,
			Direction: towardsLight,
			Depth:     0,
		}

		for _, object := range r.objects {
			intersects, t := object.Intersection(rayToLight)

			if intersects && t >= 0.00001 {
				continue Lights
			}
		}

		ld := Dot(towardsLight, normal.Normalize())

		if ld > 0 {
			color = color.Merge(light.Color, ld*closestObj.Reflectivity())
		}

		spec := Dot(reflectedRay.Direction.Normalize(), towardsLight.Normalize())

		if spec > 0 {
			spec = math.Pow(math.Pow(math.Pow(spec, 2), 2), 2)
			spec *= closestObj.Specular()
			specColor := light.Color.Mult(spec)
			color = color.Add(specColor)
		}
	}

	reflection := r.Trace(reflectedRay)

	return color.Merge(reflection, closestObj.Mirror())
}
