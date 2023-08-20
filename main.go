package main

import (
	"fmt"
	"os"

	"github.com/b-erhart/raytracer/canvas"
	"github.com/b-erhart/raytracer/geometry"
)

func main() {
	canv := canvas.NewCanvas(1920, 1080)

	eye := geometry.Vector{X: 0, Y: 0, Z: 0}
	lookAt := geometry.Vector{X: 0, Y: 0, Z: 1.2}
	up := geometry.Vector{X: 0, Y: 1, Z: 0}

	view := geometry.NewView(canv, eye, lookAt, up, 55)

	var objects []geometry.Object
	objects = append(objects, geometry.Sphere{Center: geometry.Vector{X: 2, Y: 0, Z: 17}, Radius: 2, Col: canvas.Color{R: 42, G: 106, B: 245}, Refl: 0.75, Mirr: 0.5, Spec: 0.5})
	objects = append(objects, geometry.Sphere{Center: geometry.Vector{X: 4, Y: 2, Z: 14}, Radius: 2, Col: canvas.Color{R: 230, G: 32, B: 183}, Refl: 0.66, Mirr: 0.2, Spec: 0.2})
	objects = append(objects, geometry.Sphere{Center: geometry.Vector{X: -3, Y: 0, Z: 10}, Radius: 2, Col: canvas.Color{R: 224, G: 38, B: 9}, Refl: 0.45, Mirr: 0.05, Spec: 0.05})
	objects = append(objects, geometry.Sphere{Center: geometry.Vector{X: 0, Y: -300, Z: 80}, Radius: 300, Col: canvas.Color{R: 6, G: 117, B: 13}, Refl: 0.25, Mirr: 0.0025, Spec: 0.0025})

	var lights []geometry.Light
	lights = append(lights, geometry.Light{
		Direction: geometry.Vector{X: 0.4, Y: -0.6, Z: 0.75},
		Color:     canvas.Color{R: 255, G: 200, B: 210},
	})

	background := canvas.Color{R: 21, G: 21, B: 21}
	ambience := canvas.Color{R: 10, G: 10, B: 30}

	raytracer := geometry.NewRaytracer(objects, &lights, background, ambience)

	raytracer.Render(view, canv)

	err := canv.WriteToPpm("./output.ppm")

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
