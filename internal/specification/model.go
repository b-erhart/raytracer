package specification

import (
	"fmt"

	"github.com/b-erhart/raytracer/internal/canvas"
	"github.com/b-erhart/raytracer/internal/geometry"
)

type ImageSpec struct {
	Camera       Camera
	Background   canvas.Color
	Lights       []geometry.Light
	SurfaceProps []SurfacePropSpec
	Spheres      []SphereSpec
	Triangles    []TriangleSpec
	Models       []WavefrontModelSpec
	SSAA         bool
}

func (i ImageSpec) Validate() error {
	err := validateMany(
		i.Camera.Validate(),
		validate(len(i.Lights) > 0, "at least one light source must be defined"),
	)
	if err != nil {
		return err
	}

	for _, prop := range i.SurfaceProps {
		if err = prop.Validate(); err != nil {
			return err
		}
	}

	for _, sphere := range i.Spheres {
		if err = sphere.Validate(); err != nil {
			return err
		}
	}

	for _, triangle := range i.Triangles {
		if err = triangle.Validate(); err != nil {
			return err
		}
	}

	for _, model := range i.Models {
		if err = model.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type Camera struct {
	Resolution struct {
		Width  int
		Height int
	}
	Position geometry.Vector
	LookAt   geometry.Vector
	Up       geometry.Vector
	Fov      float64
}

func (c Camera) Validate() error {
	return validateMany(
		validate(c.Resolution.Width > 0, "camera resolution width must be greater than 0"),
		validate(c.Resolution.Height > 0, "camera resolution height must be greater than 0"),
		validate(c.Fov > 0 && c.Fov < 180, "camera FOV must be between 0 and 180 degrees"),
		validate(c.Up != geometry.Vector{}, "camera up vector must not be zero vector"),
		validate(c.LookAt != geometry.Vector{}, "camera lookAt vector must not be zero vector"),
		validate(c.Position != c.LookAt, "camera position and lookAt vector must different"),
	)
}

type SurfacePropSpec struct {
	Name         string
	Color        canvas.Color
	Reflectivity float64
	Mirror       float64
	Specular     float64
}

func (p SurfacePropSpec) Validate() error {
	return validateMany(
		validate(p.Name != "", "surface property name must not be empty"),
		validate(p.Reflectivity >= 0 && p.Reflectivity <= 1, "surface property reflectivity must be between 0 and 1"),
		validate(p.Mirror >= 0 && p.Mirror <= 1, "surface property mirror must be between 0 and 1"),
		validate(p.Specular >= 0 && p.Specular <= 1, "surface property specular must be between 0 and 1"),
	)
}

type SphereSpec struct {
	Center      geometry.Vector
	Radius      float64
	SurfaceProp string
}

func (s SphereSpec) Validate() error {
	return validateMany(
		validate(s.Radius > 0, "sphere radius must be greater than 0"),
		validate(s.SurfaceProp != "", "sphere must have a surface property assigned"),
	)
}

type TriangleSpec struct {
	Corners     [3]geometry.Vector
	SurfaceProp string
}

func (t TriangleSpec) Validate() error {
	cornersDifferent := t.Corners[0] != t.Corners[1] && t.Corners[1] != t.Corners[2] && t.Corners[2] != t.Corners[0]

	return validateMany(
		validate(cornersDifferent, "triangle corners have different coordinates"),
		validate(t.SurfaceProp != "", "triangle must have a surface property assigned"),
	)
}

type WavefrontModelSpec struct {
	Path        string
	Size        float64
	Center      geometry.Vector
	Rotation    geometry.Vector
	SurfaceProp string
}

func (o WavefrontModelSpec) Validate() error {
	return validateMany(
		validate(o.Path != "", "model path must not be empty"),
		validate(o.Size > 0, "model size must be greater than 0"),
		validate(o.SurfaceProp != "", "model must have a surface property assigned"),
	)
}

func validateMany(assertions ...error) error {
	for _, err := range assertions {
		if err != nil {
			return err
		}
	}

	return nil
}

func validate(assertion bool, format string, args ...interface{}) error {
	if !assertion {
		return fmt.Errorf(format, args...)
	}

	return nil
}
