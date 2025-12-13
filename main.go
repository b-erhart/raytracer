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
	canv, view, objs, lights, background, ssaa, err := specification.Read("SPEC/image.json")

	if err != nil {
		log.Fatalf("failed to read image specification: %v", err)
	}

	log.Println("Image spec read sucessfully!")

	if ssaa {
		log.Println("SSAA enabled - rendering at doubled resolution...")
	}

	raytracer := geometry.NewRaytracer(objs, lights, background)

	log.Println("Rendering image...")
	start := time.Now()
	raytracer.Render(view, &canv)
	elapsed := time.Since(start)
	log.Printf("Rendering done! (took %s)\n", elapsed)
	if ssaa {
		canv = *canv.ApplySSAA()
	}

	log.Println("Writing PPM file...")
	err = canv.WriteToPpm("./output.ppm")

	if err != nil {
		log.Fatalf("failed to write PPM file: %v", err)
	}
	log.Println("Done!")
}
