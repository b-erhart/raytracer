package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/b-erhart/raytracer/internal/geometry"
	"github.com/b-erhart/raytracer/internal/specification"
)

func main() {
	f, err := os.Create("raytracer.prof")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create profiling file: %v", err)
		os.Exit(1)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	scene, err := specification.CreateSceneFromSpecFile("SPEC/image.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read image specification: %v", err)
		os.Exit(1)
	}

	fmt.Println("Image spec read successfully!")

	if scene.SSAA {
		fmt.Println("SSAA enabled - rendering at doubled resolution...")
	}

	raytracer := geometry.NewRaytracer(scene.Objects, scene.Lights, scene.Background)

	fmt.Println("Rendering image...")
	start := time.Now()
	raytracer.Render(scene.View, scene.Canvas)
	elapsed := time.Since(start)
	fmt.Printf("Rendering done! (took %s)\n", elapsed)
	if scene.SSAA {
		scene.Canvas = scene.Canvas.CreateSSAACanvas()
	}

	fmt.Println("Writing PPM file...")
	err = scene.Canvas.WriteToPpm("./output.ppm")

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write PPM file: %v", err)
		os.Exit(1)
	}
	fmt.Println("Done!")
}
