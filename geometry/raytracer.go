package geometry

import (
	"math"
	"sync"

	"github.com/b-erhart/raytracer/canvas"
)

const Epsilon = 0.0000001

type Raytracer struct {
	objects    []Object
	lights     []Light
	background canvas.Color
	bvhTree    BvhTree
	// ambience   canvas.Color
}

func NewRaytracer(objects []Object, lights []Light, background canvas.Color) *Raytracer {
	return &Raytracer{objects, lights, background, ConstructBvhTree(objects)}
}

func (r *Raytracer) Render(view View, canv *canvas.Canvas) {
	// fmt.Printf("%v\n", r.bvhTree)
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

	relevantObjs := r.bvhTree.GetRelevantObjects(ray)

	for i := 0; i < len(relevantObjs); i++ {
		intersects, t := relevantObjs[i].Intersection(ray)

		if intersects && t >= Epsilon && (closestObj == nil || t < tMin) {
			closestObj = relevantObjs[i]
			tMin = t
		}
	}

	if closestObj == nil && ray.Depth == 0 {
		return r.background
	} else if closestObj == nil {
		return canvas.Color{}
	} else if closestObj.Props().Reflectivity <= 0 {
		return closestObj.Props().Color
	}

	surface := closestObj.Props()
	color := surface.Color
	point := ray.At(tMin)
	normal := closestObj.SurfaceNormal(point)

	reflect := Sub(ray.Direction, Sprod(Sprod(normal, Dot(normal, ray.Direction)), 2))
	reflectedRay := Ray{
		Origin:    point,
		Direction: reflect,
		Depth:     ray.Depth + 1,
	}

Lights:
	for i := 0; i < len(r.lights); i++ {
		towardsLight := Sprod(r.lights[i].Direction, -1).Normalize()
		rayToLight := Ray{
			Origin:    point,
			Direction: towardsLight,
			Depth:     0,
		}

		lightRelevantObjs := r.bvhTree.GetRelevantObjects(rayToLight)

		for j := 0; j < len(lightRelevantObjs); j++ {
			intersects, t := lightRelevantObjs[j].Intersection(rayToLight)

			if intersects && t >= Epsilon {
				continue Lights
			}
		}

		ld := Dot(towardsLight, normal.Normalize())

		if ld > 0 {
			color = color.Merge(r.lights[i].Color, ld*surface.Reflectivity)
		}

		spec := Dot(reflectedRay.Direction.Normalize(), towardsLight.Normalize())

		if spec > 0 {
			spec = math.Pow(math.Pow(math.Pow(spec, 2), 2), 2)
			spec *= surface.Specular
			specColor := r.lights[i].Color.Mult(spec)
			color = color.Add(specColor)
		}
	}

	reflection := r.Trace(reflectedRay)

	return color.Merge(reflection, surface.Mirror)
}
