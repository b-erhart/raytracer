package canvas

import (
	"fmt"
	"math"
)

// RGB color.
type Color struct {
	R uint8
	G uint8
	B uint8
}

func (a Color) Add(b Color) Color {
	return Color{
		R: addClamped(a.R, b.R),
		G: addClamped(a.G, b.G),
		B: addClamped(a.B, b.B),
	}
}

func (c Color) Mult(f float64) Color {
	if f >= 1 {
		return c
	} else if f <= 0 {
		return Color{}
	}

	return Color{
		R: uint8(float64(c.R) * f),
		G: uint8(float64(c.G) * f),
		B: uint8(float64(c.B) * f),
	}
}

func (c Color) Merge(other Color, factor float64) Color {
	f := math.Min(math.Max(factor, 0), 1)
	t := 1 - f

	return Color{
		R: addClamped(uint8(float64(c.R)*t), uint8(float64(other.R)*f)),
		G: addClamped(uint8(float64(c.G)*t), uint8(float64(other.G)*f)),
		B: addClamped(uint8(float64(c.B)*t), uint8(float64(other.B)*f)),
	}
}

// Get string representation of color.
func (c Color) String() string {
	return fmt.Sprintf("(%3d, %3d, %3d)", c.R, c.G, c.B)
}

func addClamped(a, b uint8) uint8 {
	if math.MaxUint8 - a < b {
		return math.MaxUint8
	}

	return a + b
}
