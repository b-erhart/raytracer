package canvas

import "fmt"

// RGB color.
type Color struct {
	R uint8
	G uint8
	B uint8
}

// Get string representation of color.
func (c Color) String() string {
	return fmt.Sprintf("(%3d, %3d, %3d)", c.R, c.G, c.B)
}
