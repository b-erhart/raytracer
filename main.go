package main

import (
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/b-erhart/raytracer/internal/geometry"
	"github.com/b-erhart/raytracer/internal/specification"
)

func main() {
	f, err := os.Create("raytracer.prof")
	if err != nil {
		log.Fatalf("failed to create profiling file: %v", err)
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	scene, err := specification.CreateSceneFromSpecFile("SPEC/image.json")

	if err != nil {
		log.Fatalf("failed to read image specification: %v", err)
	}

	log.Println("Image spec read sucessfully!")

	if scene.SSAA {
		log.Println("SSAA enabled - rendering at doubled resolution...")
	}

	raytracer := geometry.NewRaytracer(scene.Objects, scene.Lights, scene.Background)

	log.Println("Rendering image...")
	start := time.Now()
	raytracer.Render(scene.View, &scene.Canvas)
	elapsed := time.Since(start)
	log.Printf("Rendering done! (took %s)\n", elapsed)
	if scene.SSAA {
		scene.Canvas = *scene.Canvas.ApplySSAA()
	}

	log.Println("Writing PPM file...")
	err = scene.Canvas.WriteToPpm("./output.ppm")

	if err != nil {
		log.Fatalf("failed to write PPM file: %v", err)
	}
	log.Println("Done!")
}
