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

func (c Color) Merge(other Color, factor float64) Color {
	f := math.Min(math.Max(factor, 0), 1)
	t := 1 - f

	return Color{
		R: uint8(float64(c.R)*t) + uint8(float64(other.R)*f),
		G: uint8(float64(c.G)*t) + uint8(float64(other.G)*f),
		B: uint8(float64(c.B)*t) + uint8(float64(other.B)*f),
	}
}

// Get string representation of color.
func (c Color) String() string {
	return fmt.Sprintf("(%3d, %3d, %3d)", c.R, c.G, c.B)
}
