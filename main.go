package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/b-erhart/raytracer/geometry"
	"github.com/b-erhart/raytracer/specification"
	"github.com/b-erhart/raytracer/wavefront"
)

func main() {
	f, err := os.Create("raytracer.prof")
	if err != nil {
		fmt.Println(err)
		return
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	canv, view, _, lights, background, err := specification.Read("image.json")

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(2)
	}

	fmt.Println("Image spec read sucessfully!")

	wavefrontFile := "teapot.obj"
	wavefrontObjects, err := wavefront.Read(wavefrontFile, geometry.Vector{X: 0, Y: 0, Z: 10}, geometry.Vector{X: 0, Y: -0.3, Z: 0}, 4)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(5)
	}

	fmt.Printf("%s read successfully\n", wavefrontFile)

	raytracer := geometry.NewRaytracer(wavefrontObjects, lights, background)

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
