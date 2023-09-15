package specification

import (
	"github.com/b-erhart/raytracer/canvas"
	"github.com/b-erhart/raytracer/geometry"
)

type ImageSpec struct {
	Camera       Camera       `validate:"required"`
	Background   canvas.Color `validate:"required"`
	Lights       []geometry.Light
	SurfaceProps []SurfaceProp
	Spheres      []Sphere
	Triangles    []Triangle
	Models       []ObjModel
}

type Camera struct {
	Resolution struct {
		Width  int `validate:"required"`
		Height int `validate:"required"`
	} `validate:"required"`
	Position geometry.Vector `validate:"required"`
	LookAt   geometry.Vector `validate:"required"`
	Up       geometry.Vector `validate:"required"`
	Fov      float64         `validate:"required"`
}

type SurfaceProp struct {
	Name         string       `validate:"required"`
	Color        canvas.Color `validate:"required"`
	Reflectivity float64      `validate:"required"`
	Mirror       float64      `validate:"required"`
	Specular     float64      `validate:"required"`
}

type Sphere struct {
	Center      geometry.Vector `validate:"required"`
	Radius      float64         `validate:"required"`
	SurfaceProp string          `validate:"required"`
}

type Triangle struct {
	Corners     [3]geometry.Vector `validate:"required"`
	SurfaceProp string             `validate:"required"`
}

type ObjModel struct {
	Path        string          `validate:"required"`
	Size        float64         `validate:"required"`
	Center      geometry.Vector `validate:"required"`
	Rotation    geometry.Vector
	SurfaceProp string          `validate:"required"`
}
