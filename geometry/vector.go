package geometry

import (
	"fmt"
	"math"
)

// 3-dimensional vector.
type Vector struct {
	X float64
	Y float64
	Z float64
}

// Get the euclidian length of a vector.
func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Get the normalized variant of a vector.
func (v *Vector) Normalize() *Vector {
	normFactor := 1 / v.Length()
	return Sprod(v, normFactor)
}

// Add two vectors.
func Add(a, b *Vector) *Vector {
	return &Vector{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Subtract the second vector from the first vector.
func Sub(a, b *Vector) *Vector {
	return &Vector{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Get the cross product of two vectors.
func Cross(a, b *Vector) *Vector {
	return &Vector{
		a.Y*b.Z - b.Y*a.Z,
		b.X*a.Z - a.X*b.Z,
		a.X*b.Y - b.X*a.Y,
	}
}

// Get the dot product of two vectors.
func Dot(a, b *Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Get the product of a vector and a scalar.
func Sprod(v *Vector, s float64) *Vector {
	return &Vector{s * v.X, s * v.Y, s * v.Z}
}

// Get the string representation of a vector.
func (v Vector) String() string {
	return fmt.Sprintf("(%f, %f, %f)", v.X, v.Y, v.Z)
}
