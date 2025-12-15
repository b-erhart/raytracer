package specification

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/b-erhart/raytracer/internal/canvas"
	"github.com/b-erhart/raytracer/internal/geometry"
	"github.com/b-erhart/raytracer/internal/wavefront"
)

func CreateSceneFromSpecFile(path string) (geometry.Scene, error) {
	spec, err := readSpecFromFile(path)
	if err != nil {
		return geometry.Scene{}, fmt.Errorf("failed to read image specification: %w", err)
	}

	objects, err := createObjects(spec, path)
	if err != nil {
		return geometry.Scene{}, fmt.Errorf("failed to create objects: %w", err)
	}

	canvasWidth := spec.Camera.Resolution.Width
	canvasHeight := spec.Camera.Resolution.Height
	if spec.SSAA {
		canvasWidth *= 2
		canvasHeight *= 2
	}

	canv := canvas.NewCanvas(canvasWidth, canvasHeight)
	view := geometry.NewView(canvasWidth, canvasHeight, spec.Camera.Position, spec.Camera.LookAt, spec.Camera.Up, spec.Camera.Fov)

	return geometry.Scene{
		Canvas:     canv,
		View:       view,
		Objects:    objects,
		Lights:     spec.Lights,
		Background: spec.Background,
		SSAA:       spec.SSAA,
	}, nil
}

func readSpecFromFile(path string) (ImageSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return ImageSpec{}, fmt.Errorf("failed to open specification file: %w", err)
	}

	defer file.Close()

	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		return ImageSpec{}, fmt.Errorf("failed to read specification file: %w", err)
	}

	spec := ImageSpec{}
	err = json.Unmarshal(jsonBytes, &spec)
	if err != nil {
		return ImageSpec{}, fmt.Errorf("failed to parse spec JSON: %w", err)
	}

	err = spec.Validate()
	if err != nil {
		return ImageSpec{}, fmt.Errorf("failed to validate specification: %w", err)
	}

	return spec, nil
}

func createObjects(s ImageSpec, specFilePath string) ([]geometry.Object, error) {
	props, err := createObjectProps(s.SurfaceProps)
	if err != nil {
		return []geometry.Object{}, fmt.Errorf("failed to create object properties: %w", err)
	}

	sphereObjects, err := createSphereObjects(s.Spheres, props)
	if err != nil {
		return []geometry.Object{}, fmt.Errorf("failed to create sphere objects: %w", err)
	}

	triangleObjects, err := createTriangleObjects(s.Triangles, props)
	if err != nil {
		return []geometry.Object{}, fmt.Errorf("failed to create triangle objects: %w", err)
	}

	wavefrontModelObjects, err := createWavefrontModelObjects(s.Models, specFilePath, props)
	if err != nil {
		return []geometry.Object{}, fmt.Errorf("failed to create wavefront model objects: %w", err)
	}

	objs := make([]geometry.Object, 0, len(sphereObjects)+len(triangleObjects)+len(wavefrontModelObjects))
	objs = append(objs, sphereObjects...)
	objs = append(objs, triangleObjects...)
	objs = append(objs, wavefrontModelObjects...)

	return objs, nil
}

func createObjectProps(surfacePropSpecs []SurfacePropSpec) (map[string]geometry.ObjectProps, error) {
	props := make(map[string]geometry.ObjectProps, len(surfacePropSpecs))

	for _, prop := range surfacePropSpecs {
		_, exists := props[prop.Name]
		if exists {
			return nil, fmt.Errorf(
				"multiple surface properties with name %q defined but name must be unique",
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

	return props, nil
}

func createSphereObjects(sphereSpecs []SphereSpec, props map[string]geometry.ObjectProps) ([]geometry.Object, error) {
	sphereObjects := make([]geometry.Object, 0, len(sphereSpecs))

	for _, sphere := range sphereSpecs {
		prop, err := lookupSurfaceProp(sphere.SurfaceProp, props)
		if err != nil {
			return []geometry.Object{}, fmt.Errorf("failed to lookup surface properties for sphere: %w", err)
		}

		sphereObjects = append(sphereObjects, &geometry.Sphere{
			Center:     sphere.Center,
			Radius:     sphere.Radius,
			Properties: prop,
		})
	}

	return sphereObjects, nil
}

func createTriangleObjects(triangleSpecs []TriangleSpec, props map[string]geometry.ObjectProps) ([]geometry.Object, error) {
	triangleObjects := make([]geometry.Object, 0, len(triangleSpecs))

	for _, triangle := range triangleSpecs {
		prop, err := lookupSurfaceProp(triangle.SurfaceProp, props)
		if err != nil {
			return []geometry.Object{}, fmt.Errorf("failed to lookup surface properties for triangle: %w", err)
		}

		triangleObjects = append(triangleObjects, &geometry.Triangle{
			A:          triangle.Corners[0],
			B:          triangle.Corners[1],
			C:          triangle.Corners[2],
			Properties: prop,
		})
	}

	return triangleObjects, nil
}

func createWavefrontModelObjects(modelSpecs []WavefrontModelSpec, specFilePath string, props map[string]geometry.ObjectProps) ([]geometry.Object, error) {
	wavefrontObjects := make([]geometry.Object, 0)

	for _, objModel := range modelSpecs {
		prop, err := lookupSurfaceProp(objModel.SurfaceProp, props)
		if err != nil {
			return []geometry.Object{}, fmt.Errorf("failed to lookup surface properties for wavefront model: %w", err)
		}

		absoluteSpecPath, err := filepath.Abs(specFilePath)
		if err != nil {
			return []geometry.Object{}, fmt.Errorf("failed to get absolute path of specification file: %w", err)
		}

		absolutePath := filepath.Join(filepath.Dir(absoluteSpecPath), objModel.Path)

		wavefrontObjs, err := wavefront.Read(absolutePath, objModel.Center, objModel.Rotation, objModel.Size, prop)

		if err != nil {
			return []geometry.Object{}, fmt.Errorf("failed to read wavefront model: %w", err)
		}

		wavefrontObjects = append(wavefrontObjects, wavefrontObjs...)
	}

	return wavefrontObjects, nil
}

func lookupSurfaceProp(name string, props map[string]geometry.ObjectProps) (geometry.ObjectProps, error) {
	prop, exists := props[name]
	if !exists {
		return geometry.ObjectProps{}, fmt.Errorf(
			"surface properties with name %q do not exist but are assigned to a wavefront model",
			name,
		)
	}

	return prop, nil
}
