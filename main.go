package main

import (
	"fmt"
	"os"
	"time"

	"github.com/b-erhart/raytracer/geometry"
	"github.com/b-erhart/raytracer/specification"
)

func main() {
	canv, view, objects, lights, background, err := specification.Read("image.json")

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(2)
	}

	fmt.Println("Image spec read sucessfully!")

	raytracer := geometry.NewRaytracer(objects, lights, background)

	fmt.Println("Rendering image...")
	start := time.Now()
	raytracer.Render(view, &canv)
	elapsed := time.Since(start)
	fmt.Printf("Rendering done! (took %s)\n", elapsed)

	fmt.Println("Writing PPM file...")
	err = canv.WriteToPpm("./output.ppm")

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}
