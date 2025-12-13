package geometry

import "fmt"

type Ray struct {
	Origin    Vector
	Direction Vector
	Depth     int
}

func (r Ray) At(distance float64) Vector {
	return Add(r.Origin, Sprod(r.Direction.Normalize(), distance))
}

func (r Ray) String() string {
	return fmt.Sprintf("%v + %d*%v", r.Origin, r.Depth, r.Direction)
}
