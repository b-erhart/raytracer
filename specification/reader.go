package specification

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/b-erhart/raytracer/canvas"
	"github.com/b-erhart/raytracer/geometry"
	"github.com/b-erhart/raytracer/wavefront"
)

func (s ImageSpec) canvas() canvas.Canvas {
	return *canvas.NewCanvas(s.Camera.Resolution.Width, s.Camera.Resolution.Height)
}

func (s ImageSpec) view() geometry.View {
	return geometry.NewView(
		s.Camera.Resolution.Width,
		s.Camera.Resolution.Height,
		s.Camera.Position,
		s.Camera.LookAt,
		s.Camera.Up,
		s.Camera.Fov,
	)
}

func (s ImageSpec) objects(specFilePath string) ([]geometry.Object, error) {
	props := make(map[string]geometry.ObjectProps, len(s.SurfaceProps))

	for _, prop := range s.SurfaceProps {
		_, exists := props[prop.Name]
		if exists {
			return []geometry.Object{}, fmt.Errorf(
				"multiple surface properties with name \"%s\" defined - name must be unique",
				prop.Name,
			)
		}

		props[prop.Name] = geometry.ObjectProps{
			Color:        prop.Color,
			Reflectivity: prop.Reflectivity,
			Mirror:       prop.Mirror,
			Specular:     prop.Specular,
		}
	}

	objs := make([]geometry.Object, 0, len(s.Spheres)+len(s.Triangles))

	for _, sphere := range s.Spheres {
		prop, exists := props[sphere.SurfaceProp]
		if !exists {
			return []geometry.Object{}, fmt.Errorf(
				"surface properties with name \"%s\" do not exists but are assigned to a sphere",
				sphere.SurfaceProp,
			)
		}

		objs = append(objs, &geometry.Sphere{
			Center:     sphere.Center,
			Radius:     sphere.Radius,
			Properties: prop,
		})
	}

	for _, triangle := range s.Triangles {
		prop, exists := props[triangle.SurfaceProp]
		if !exists {
			return []geometry.Object{}, fmt.Errorf(
				"surface properties with name \"%s\" do not exist but are assigned to a sphere",
				triangle.SurfaceProp,
			)
		}

		objs = append(objs, &geometry.Triangle{
			A:          triangle.Corners[0],
			B:          triangle.Corners[1],
			C:          triangle.Corners[2],
			Properties: prop,
		})
	}

	for _, objModel := range s.Models {
		prop, exists := props[objModel.SurfaceProp]
		if !exists {
			return []geometry.Object{}, fmt.Errorf(
				"surface properties with name \"%s\" do not exist but are assigned to a sphere",
				objModel.SurfaceProp,
			)
		}

		absoluteSpecPath, err := filepath.Abs(specFilePath)

		if err != nil {
			return []geometry.Object{}, err
		}

		absolutePath := filepath.Join(filepath.Dir(absoluteSpecPath), objModel.Path)

		wavefrontObjs, err := wavefront.Read(absolutePath, objModel.Center, objModel.Rotation, objModel.Size, prop)

		if err != nil {
			return []geometry.Object{}, err
		}

		objs = append(objs, wavefrontObjs...)
	}

	return objs, nil
}

func Read(path string) (canvas.Canvas, geometry.View, []geometry.Object, []geometry.Light, canvas.Color, error) {
	file, err := os.Open(path)
	if err != nil {
		return canvas.Canvas{}, geometry.View{}, []geometry.Object{}, []geometry.Light{}, canvas.Color{}, err
	}

	defer file.Close()

	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		return canvas.Canvas{}, geometry.View{}, []geometry.Object{}, []geometry.Light{}, canvas.Color{}, err
	}

	valid := json.Valid(jsonBytes)
	if !valid {
		err = fmt.Errorf("the specification in %s is not valid", path)
		return canvas.Canvas{}, geometry.View{}, []geometry.Object{}, []geometry.Light{}, canvas.Color{}, err
	}

	var spec ImageSpec

	err = json.Unmarshal(jsonBytes, &spec)

	if err != nil {
		return canvas.Canvas{}, geometry.View{}, []geometry.Object{}, []geometry.Light{}, canvas.Color{}, err
	}

	objs, err := spec.objects(path)
	if err != nil {
		return canvas.Canvas{}, geometry.View{}, []geometry.Object{}, []geometry.Light{}, canvas.Color{}, err
	}

	return spec.canvas(), spec.view(), objs, spec.Lights, spec.Background, nil
}
