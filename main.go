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
	objects = append(objects, geometry.Sphere{Center: geometry.Vector{X: 2, Y: 0, Z: 17}, Radius: 2})
	objects = append(objects, geometry.Sphere{Center: geometry.Vector{X: 4, Y: 2, Z: 14}, Radius: 2})
	objects = append(objects, geometry.Sphere{Center: geometry.Vector{X: -3, Y: 0, Z: 10}, Radius: 2})

	var lights []geometry.Light
	lights = append(lights, geometry.Light{
		Direction: geometry.Vector{X: -100, Y: -100, Z: -100},
		Color:     canvas.Color{R: 150, G: 255, B: 150},
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
