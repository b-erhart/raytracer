package canvas

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

// Canvas to draw RGB values to.
type Canvas struct {
	width  int
	height int
	// Red pixel values
	R [][]uint8
	// Green pixel values
	G [][]uint8
	// Blue pixel values
	B [][]uint8
}

// Create a new canvas with specified width and height. Initialize R, G and B
// slices accordingly.
func NewCanvas(width, height int) *Canvas {
	canvas := Canvas{height: height, width: width}

	canvas.R = make([][]uint8, width)
	canvas.G = make([][]uint8, width)
	canvas.B = make([][]uint8, width)

	for i := 0; i < width; i++ {
		canvas.R[i] = make([]uint8, height)
		canvas.G[i] = make([]uint8, height)
		canvas.B[i] = make([]uint8, height)
	}

	return &canvas
}

// Get canvas width in pixels.
func (canvas *Canvas) Width() int {
	return canvas.width
}

// Get canvas height in pixels.
func (canvas *Canvas) Height() int {
	return canvas.height
}

// Set r, g and b values of the pixel at coordinates (x, y). Panics if the
// pixels given are out of bounds.
func (canvas *Canvas) SetPixel(x, y int, r, g, b uint8) error {
	if x < 0 || y < 0 || x >= canvas.width || y >= canvas.height {
		panic(fmt.Sprintf("pixel coordinates out of bounds - tried to access pixel (%d, %d) in a %dx%d canvas", x, y, canvas.width, canvas.height))
	}

	canvas.R[x][y] = r
	canvas.G[x][y] = g
	canvas.B[x][y] = b

	return nil
}

func (canvas *Canvas) SetPixel(x, y int, color Color) error {
	return canvas.SetPixel(x, y, color.R, color.G, color.B)
}

// Write the canvas to a PPM (P6) file. If a file exists at the given path, it
// is moved to "<path>.bak". Return an error if writing the file fails.
func (canvas *Canvas) WriteToPpm(path string) error {
	err := os.Rename(path, path+".bak")

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(fmt.Sprintf("P6\n%d %d\n%d\n", canvas.width, canvas.height, math.MaxUint8))

	if err != nil {
		return err
	}

	for j := 0; j < canvas.height; j++ {
		for i := 0; i < canvas.width; i++ {
			for _, byte := range []byte{canvas.R[i][j], canvas.G[i][j], canvas.B[i][j]} {
				err = writer.WriteByte(byte)

				if err != nil {
					return err
				}
			}
		}

		if err != nil {
			return err
		}
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	return nil
}

// Create string representation of the canvas.
func (canvas Canvas) String() string {
	var strBuilder strings.Builder

	strBuilder.WriteString(fmt.Sprintf("canvas (%dx%d) = [\n", canvas.width, canvas.height))
	for j := 0; j < canvas.height; j++ {
		strBuilder.WriteString("\t[")
		for i := 0; i < canvas.width; i++ {
			strBuilder.WriteString(fmt.Sprintf("(%d, %d, %d)", canvas.R[i][j], canvas.G[i][j], canvas.B[i][j]))

			if i < canvas.width-1 {
				strBuilder.WriteString(", ")
			}
		}
		strBuilder.WriteString("]\n")
	}
	strBuilder.WriteString("]")

	return strBuilder.String()
}
