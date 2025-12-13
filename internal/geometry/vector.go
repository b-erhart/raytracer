package geometry

import (
	"fmt"
	"math"

	"github.com/b-erhart/raytracer/internal/canvas"
)

// 3-dimensional vector.
type Vector struct {
	X float64
	Y float64
	Z float64
}

// Get the euclidian length of a vector.
func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Get the normalized variant of a vector.
func (v Vector) Normalize() Vector {
	normFactor := 1 / v.Length()
	return Sprod(v, normFactor)
}

// Add two vectors.
func Add(a, b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Subtract the second vector from the first vector.
func Sub(a, b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Get the cross product of two vectors.
func Cross(a, b Vector) Vector {
	return Vector{
		a.Y*b.Z - b.Y*a.Z,
		b.X*a.Z - a.X*b.Z,
		a.X*b.Y - b.X*a.Y,
	}
}

// Get the dot product of two vectors.
func Dot(a, b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Get the product of a vector and a scalar.
func Sprod(v Vector, s float64) Vector {
	return Vector{s * v.X, s * v.Y, s * v.Z}
}

func (v Vector) ToColor() canvas.Color {
	nv := v.Normalize()
	var r, g, b uint8

	if nv.X > 1 {
		r = 255
	} else {
		r = uint8((0.5 * (nv.X + 1)) * 255)
	}

	if nv.Y > 1 {
		g = 255
	} else {
		g = uint8((0.5 * (nv.Y + 1)) * 255)
	}

	if nv.Z > 1 {
		b = 255
	} else {
		b = uint8((0.5 * (nv.Z + 1)) * 255)
	}

	return canvas.Color{R: r, G: g, B: b}
}

// Get the string representation of a vector.
func (v Vector) String() string {
	return fmt.Sprintf("(%f, %f, %f)", v.X, v.Y, v.Z)
}
