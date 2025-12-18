package wavefront

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/b-erhart/raytracer/internal/geometry"
)

type fileContent struct {
	vertices      []geometry.Vector
	vertexNormals []geometry.Vector
	faces         []geometry.Triangle
	maxVertex     geometry.Vector
	minVertex     geometry.Vector
}

func Read(path string, origin, rotation geometry.Vector, scaling float64, props geometry.ObjectProps) ([]geometry.Object, error) {
	logger := log.Default()
	logger.Printf("reading wavefront file \"%s\"\n", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := parseFile(file, logger)
	if err != nil {
		return nil, err
	}

	return turnContentToObjects(content, scaling, rotation, origin, props)
}

func parseFile(file *os.File, logger *log.Logger) (fileContent, error) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	unsupportedDirectives := make([]string, 0)
	lineNr := 0

	content := fileContent{
		vertices:      make([]geometry.Vector, 0),
		vertexNormals: make([]geometry.Vector, 0),
		faces:         make([]geometry.Triangle, 0),
		maxVertex:     geometry.Vector{X: math.Inf(-1), Y: math.Inf(-1), Z: math.Inf(-1)},
		minVertex:     geometry.Vector{X: math.Inf(1), Y: math.Inf(1), Z: math.Inf(1)},
	}

	for scanner.Scan() {
		lineNr++

		line := scanner.Text()
		words := strings.Fields(line)

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "#":
			continue
		case "v":
			newVertex, err := readVector(words)
			if err != nil {
				return fileContent{}, fmt.Errorf("unable to parse vertex on line %d: %v", lineNr, err)
			}

			updateExtremes(&content.minVertex, &content.maxVertex, newVertex)

			content.vertices = append(content.vertices, newVertex)
		case "vn":
			newVertexNormal, err := readVector(words)
			if err != nil {
				return fileContent{}, fmt.Errorf("unable to parse vertex normal on line %d: %v", lineNr, err)
			}

			content.vertexNormals = append(content.vertexNormals, newVertexNormal)
		case "f":
			newFace, err := readFace(words, content.vertices, content.vertexNormals)
			if err != nil {
				return fileContent{}, err
			}

			content.faces = append(content.faces, newFace...)
		default:
			if !slices.Contains(unsupportedDirectives, words[0]) {
				logger.Printf("unsupported directive \"%s\" found - will be ignored", words[0])
				unsupportedDirectives = append(unsupportedDirectives, words[0])
			}
		}
	}

	return content, nil
}

func turnContentToObjects(content fileContent, scaling float64, rotation geometry.Vector, origin geometry.Vector, props geometry.ObjectProps) ([]geometry.Object, error) {
	objs := make([]geometry.Object, 0, len(content.faces))
	centering := geometry.Vector{
		X: -content.minVertex.X - (content.maxVertex.X-content.minVertex.X)/2,
		Y: -content.minVertex.Y - (content.maxVertex.Y-content.minVertex.Y)/2,
		Z: -content.minVertex.Z - (content.maxVertex.Z-content.minVertex.Z)/2,
	}

	trianglesPerCorner := make(map[geometry.Vector][]*geometry.Triangle)
	size := math.Max(content.maxVertex.X-content.minVertex.X, math.Max(content.maxVertex.Y-content.minVertex.Y, content.maxVertex.Z-content.minVertex.Z))
	scalingFactor := scaling / size

	for _, f := range content.faces {
		a := geometry.Vector{
			X: ((f.A.X + centering.X) * scalingFactor),
			Y: ((f.A.Y + centering.Y) * scalingFactor),
			Z: ((f.A.Z + centering.Z) * scalingFactor),
		}
		b := geometry.Vector{
			X: ((f.B.X + centering.X) * scalingFactor),
			Y: ((f.B.Y + centering.Y) * scalingFactor),
			Z: ((f.B.Z + centering.Z) * scalingFactor),
		}
		c := geometry.Vector{
			X: ((f.C.X + centering.X) * scalingFactor),
			Y: ((f.C.Y + centering.Y) * scalingFactor),
			Z: ((f.C.Z + centering.Z) * scalingFactor),
		}

		a = geometry.Add(rotate(a, rotation), origin)
		b = geometry.Add(rotate(b, rotation), origin)
		c = geometry.Add(rotate(c, rotation), origin)

		triangle := &geometry.Triangle{A: a, B: b, C: c, Properties: props}

		if f.NormalsSet {
			triangle.ASurfaceNormal = rotate(f.ASurfaceNormal, rotation).Normalize()
			triangle.BSurfaceNormal = rotate(f.BSurfaceNormal, rotation).Normalize()
			triangle.CSurfaceNormal = rotate(f.CSurfaceNormal, rotation).Normalize()
			triangle.NormalsSet = true
		}

		objs = append(objs, triangle)
		trianglesPerCorner[a] = append(trianglesPerCorner[a], triangle)
		trianglesPerCorner[b] = append(trianglesPerCorner[b], triangle)
		trianglesPerCorner[c] = append(trianglesPerCorner[c], triangle)
	}

	for _, obj := range objs {
		triangle, isTriangle := obj.(*geometry.Triangle)
		if !isTriangle {
			continue
		}

		if !triangle.NormalsSet {
			triangle.ASurfaceNormal = calculateCornerNormal(triangle.A, triangle, trianglesPerCorner)
			triangle.BSurfaceNormal = calculateCornerNormal(triangle.B, triangle, trianglesPerCorner)
			triangle.CSurfaceNormal = calculateCornerNormal(triangle.C, triangle, trianglesPerCorner)
			triangle.NormalsSet = true
		}
	}

	return objs, nil
}

func calculateCornerNormal(corner geometry.Vector, triangle *geometry.Triangle, trianglesPerCorner map[geometry.Vector][]*geometry.Triangle) geometry.Vector {
	normal := triangle.TriangleNormal()

	for _, otherTriangle := range trianglesPerCorner[corner] {
		if otherTriangle == triangle {
			continue
		}

		dot := geometry.Dot(triangle.TriangleNormal(), otherTriangle.TriangleNormal())

		if dot > 0+geometry.Epsilon {
			normal = geometry.Add(normal, otherTriangle.TriangleNormal())
		}
	}
	normal = normal.Normalize()

	return normal
}

func readVector(words []string) (geometry.Vector, error) {
	if len(words) < 4 {
		return geometry.Vector{}, fmt.Errorf("invalid vertex definition: expected 3 elements but got %d", len(words)-1)
	}

	elements := make([]float64, 3)
	for i := range elements {
		var err error
		elements[i], err = strconv.ParseFloat(words[i+1], 64)
		if err != nil {
			return geometry.Vector{}, fmt.Errorf("invalid vertex definition: element #%d is not a valid number", i+1)
		}
	}

	return geometry.Vector{X: elements[0], Y: elements[1], Z: elements[2]}, nil
}

func readFace(words []string, vertices, vertexNormals []geometry.Vector) ([]geometry.Triangle, error) {
	if words[0] != "f" {
		panic("got a face definition line that does not start with 'f'")
	} else if len(words) < 4 {
		return nil, fmt.Errorf("invalid face definition: faces must have at least 3 corner vertices")
	} else if len(words) > 5 {
		// TODO: implement Seidel's algorithm to support any polygon
		return nil, fmt.Errorf("invalid face definition: faces with more than four corners are currently not supported")
	}

	corners := make([]geometry.Vector, 0, 3)
	normals := make([]geometry.Vector, 0, 3)

	for i := 1; i < len(words); i++ {
		cornerSpec := strings.Split(words[i], "/")
		vIndexStr := cornerSpec[0]
		vIndex, err := strconv.ParseInt(vIndexStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid face definition: element #%d is not a valid number", i)
		}

		if vIndex > int64(len(vertices)) {
			return nil, fmt.Errorf("invalid face definition: vertex #%d is referenced but not defined", vIndex)
		} else if int(vIndex) == 0 {
			return nil, fmt.Errorf("invalid face definition: vertex number must be greater than 0 but is %d", vIndex)
		}

		if vIndex < 0 {
			vIndex = int64(len(vertices)) + vIndex + 1
		}

		corners = append(corners, vertices[vIndex-1])

		if len(cornerSpec) >= 3 {
			vnIndexStr := cornerSpec[2]
			vnIndex, err := strconv.ParseInt(vnIndexStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid face definition: element #%d is not a valid number", i)
			}

			if vnIndex > int64(len(vertexNormals)) {
				return nil, fmt.Errorf("invalid face definition: vertex normal #%d is referenced but not defined", vIndex)
			} else if int(vnIndex) == 0 {
				return nil, fmt.Errorf("invalid face definition: vertex normal number must be greater than 0 but is %d", vIndex)
			}

			if vnIndex < 0 {
				vnIndex = int64(len(vertexNormals)) + vnIndex + 1
			}

			normals = append(normals, vertexNormals[vnIndex-1])
		}
	}

	triangles := make([]geometry.Triangle, 0)

	t := geometry.Triangle{A: corners[0], B: corners[1], C: corners[2]}

	if len(corners) == len(normals) {
		t.ASurfaceNormal = normals[0]
		t.BSurfaceNormal = normals[1]
		t.CSurfaceNormal = normals[2]
		t.NormalsSet = true
	}

	triangles = append(triangles, t)

	if len(corners) == 4 {
		t2 := geometry.Triangle{A: corners[0], B: corners[2], C: corners[3]}

		if len(corners) == len(normals) {
			t2.ASurfaceNormal = normals[0]
			t2.BSurfaceNormal = normals[2]
			t2.CSurfaceNormal = normals[3]
			t2.NormalsSet = true
		}

		triangles = append(triangles, t2)
	}

	return triangles, nil
}

func rotate(vertex, rotation geometry.Vector) geometry.Vector {
	// rotate around x axis
	xRotated := geometry.Vector{
		X: vertex.X,
		Y: vertex.Y*math.Cos(rotation.X*math.Pi) - vertex.Z*math.Sin(rotation.X*math.Pi),
		Z: vertex.Y*math.Sin(rotation.X*math.Pi) + vertex.Z*math.Cos(rotation.X*math.Pi),
	}

	// rotate around y axis
	xyRotated := geometry.Vector{
		X: xRotated.X*math.Cos(rotation.Y*math.Pi) + xRotated.Z*math.Sin(rotation.Y*math.Pi),
		Y: xRotated.Y,
		Z: -xRotated.X*math.Sin(rotation.Y*math.Pi) + xRotated.Z*math.Cos(rotation.Y*math.Pi),
	}

	// rotate around z axis
	xyzRotated := geometry.Vector{
		X: xyRotated.X*math.Cos(rotation.Z*math.Pi) - xyRotated.Y*math.Sin(rotation.Z*math.Pi),
		Y: xyRotated.X*math.Sin(rotation.Z*math.Pi) + xyRotated.Y*math.Cos(rotation.Z*math.Pi),
		Z: xyRotated.Z,
	}

	return xyzRotated
}

func updateExtremes(minVertex, maxVertex *geometry.Vector, newVertex geometry.Vector) {
	minVertex.X = math.Min(minVertex.X, newVertex.X)
	minVertex.Y = math.Min(minVertex.Y, newVertex.Y)
	minVertex.Z = math.Min(minVertex.Z, newVertex.Z)

	maxVertex.X = math.Max(maxVertex.X, newVertex.X)
	maxVertex.Y = math.Max(maxVertex.Y, newVertex.Y)
	maxVertex.Z = math.Max(maxVertex.Z, newVertex.Z)
}
