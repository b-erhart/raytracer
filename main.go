package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/b-erhart/raytracer/geometry"
	"github.com/b-erhart/raytracer/specification"
)

func main() {
	f, err := os.Create("raytracer.prof")
	if err != nil {
		fmt.Println(err)
		return
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	canv, view, objs, lights, background, err := specification.Read("SPEC/image.json")

	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(2)
	}

	fmt.Println("Image spec read sucessfully!")

	raytracer := geometry.NewRaytracer(objs, lights, background)

	fmt.Println("Rendering image...")
	start := time.Now()
	raytracer.Render(view, &canv)
	elapsed := time.Since(start)
	fmt.Printf("Rendering done! (took %s)\n", elapsed)

	fmt.Println("Writing PPM file...")
	err = canv.WriteToPpm("./output.ppm")

	if err != nil {
		fmt.Println("[ERROR]", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}
