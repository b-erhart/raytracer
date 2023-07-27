package geometry

import (
	"github.com/b-erhart/raytracer/canvas"
	"math"
)

type View struct {
	Eye    Vector
	LookAt Vector
	Up     Vector
	Fov    float64

	u          Vector
	v          Vector
	du         Vector
	dv         Vector
	bottomLeft Vector
}

func NewView(canvas *canvas.Canvas, eye, lookAt, up *Vector, fov float64) *View {
	var view *View

	view.Eye = *eye
	view.LookAt = *lookAt
	view.Up = *up
	view.Fov = fov

	lxup := Cross(lookAt, up)
	view.u = *Sprod(lxup, 1/lxup.Length()).Normalize()

	lxu := Cross(lookAt, &view.u)
	view.v = *Sprod(lxu, 1/lxu.Length()).Normalize()

	aspectRatio := float64(canvas.Width()) / float64(canvas.Height())

	uLen := math.Tan(fov * (math.Pi / 180))
	vLen := uLen * aspectRatio

	view.du = *Sprod(&view.u, uLen/float64(canvas.Width()-1))
	view.dv = *Sprod(&view.v, vLen/float64(canvas.Height()-1))

	return view
}
