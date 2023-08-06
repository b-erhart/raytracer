package geometry

import (
	"math"

	"github.com/b-erhart/raytracer/canvas"
)

type View struct {
	eye    Vector
	lookAt Vector
	up     Vector
	fov    float64

	u          Vector
	v          Vector
	du         Vector
	dv         Vector
	bottomLeft Vector
}

func NewView(canv *canvas.Canvas, eye, lookAt, up Vector, fov float64) View {
	var view View

	view.eye = eye
	view.lookAt = lookAt
	view.up = up
	view.fov = fov

	lxup := Cross(lookAt, up)
	view.u = Sprod(lxup, -1/lxup.Length()).Normalize()

	lxu := Cross(lookAt, view.u)
	view.v = Sprod(lxu, -1/lxu.Length()).Normalize()

	aspectRatio := float64(canv.Height()) / float64(canv.Width())

	uLen := math.Tan(fov * (math.Pi / 180))
	vLen := uLen * aspectRatio

	view.du = Sprod(view.u, uLen/float64(canv.Width()-1))
	view.dv = Sprod(view.v, vLen/float64(canv.Height()-1))

	centerToLeft := Sprod(view.du, float64(-1*(canv.Width()/2)))
	centerToBottom := Sprod(view.dv, float64(-1*(canv.Height()/2)))
	view.bottomLeft = Add(Add(eye, lookAt), Add(centerToLeft, centerToBottom))

	return view
}

func (v View) Eye() Vector {
	return v.eye
}

func (v View) LookAt() Vector {
	return v.lookAt
}

func (v View) Up() Vector {
	return v.up
}

func (v View) Fov() float64 {
	return v.fov
}

func (v View) U() Vector {
	return v.u
}

func (v View) V() Vector {
	return v.v
}

func (v View) Du() Vector {
	return v.du
}

func (v View) Dv() Vector {
	return v.dv
}

func (v View) BottomLeft() Vector {
	return v.bottomLeft
}
