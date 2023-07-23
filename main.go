package main

import (
	"fmt"
	"os"

	"github.com/b-erhart/raytracer/canvas"
)

func main() {
	canvas := canvas.NewCanvas(120, 80)

	for j := 0; j < 10; j++ {
		for i := 0; i < canvas.Width(); i++ {
			if i < canvas.Width()/2 {
				canvas.SetPixel(i, j, 0, 0, 0)
			} else {
				canvas.SetPixel(i, j, 255, 255, 255)
			}
		}
	}

	for j := 10; j < canvas.Height(); j++ {
		for i := 0; i < canvas.Width(); i++ {
			switch {
			case i < 1*canvas.Width()/6:
				canvas.SetPixel(i, j, 115, 11, 219)
			case i < 2*canvas.Width()/6:
				canvas.SetPixel(i, j, 11, 105, 219)
			case i < 3*canvas.Width()/6:
				canvas.SetPixel(i, j, 54, 138, 12)
			case i < 4*canvas.Width()/6:
				canvas.SetPixel(i, j, 222, 212, 16)
			case i < 5*canvas.Width()/6:
				canvas.SetPixel(i, j, 222, 129, 16)
			default:
				canvas.SetPixel(i, j, 222, 16, 16)
			}
		}
	}

	err := canvas.WriteToPpm("./output.ppm")

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
