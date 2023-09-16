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

	"github.com/b-erhart/raytracer/geometry"
)

func Read(path string, origin, rotation geometry.Vector, scaling float64, props geometry.ObjectProps) ([]geometry.Object, error) {
	logger := log.Default()

	logger.Printf("reading wavefront file \"%s\"\n", path)

	file, err := os.Open(path)
	if err != nil {
		return []geometry.Object{}, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	unsups := make([]string, 0, 5)
	lineNr := 0
	vs := make([]geometry.Vector, 0, 100)
	vns := make([]geometry.Vector, 0, 100)
	fs := make([]geometry.Triangle, 0, 100)
	minV := geometry.Vector{X: math.MaxInt, Y: math.MaxInt, Z: math.MaxInt}
	maxV := geometry.Vector{X: math.MinInt, Y: math.MinInt, Z: math.MinInt}
	for scanner.Scan() {
		lineNr++

		line := scanner.Text()
		words := deleteEmpty(strings.Split(line, " "))

		if len(words) == 0 || words[0][0] == '#' {
			continue
		}

		switch words[0] {
		case "#":
			continue
		case "v":
			v, err := readVector(words)
			if err != nil {
				logger.Printf("line %d: %v\n", lineNr, err)
				return []geometry.Object{}, err
			}

			updateExtremes(&minV, &maxV, v)

			vs = append(vs, v)
		case "vn":
			vn, err := readVector(words)
			if err != nil {
				logger.Printf("line %d: %v\n", lineNr, err)
				return []geometry.Object{}, err
			}

			vns = append(vns, vn)
		case "f":
			f, err := readFace(words, vs, vns)
			if err != nil {
				logger.Printf("line %d: %v\n", lineNr, err)
				return []geometry.Object{}, err
			}

			fs = append(fs, f...)
		default:
			if !slices.Contains(unsups, words[0]) {
				logger.Printf("unsupported directive \"%s\" found - will be ignored\n", words[0])
				unsups = append(unsups, words[0])
			}
		}
	}

	objs := make([]geometry.Object, 0, len(fs))
	centering := geometry.Vector{
		X: -minV.X - (maxV.X-minV.X)/2,
		Y: -minV.Y - (maxV.Y-minV.Y)/2,
		Z: -minV.Z - (maxV.Z-minV.Z)/2,
	}

	trianglesPerCorner := make(map[geometry.Vector][]*geometry.Triangle)
	size := math.Max(maxV.X-minV.X, math.Max(maxV.Y-minV.Y, maxV.Z-minV.Z))
	scalingFactor := scaling / size

	for i := 0; i < len(fs); i++ {
		a := geometry.Vector{
			X: ((fs[i].A.X + centering.X) * scalingFactor),
			Y: ((fs[i].A.Y + centering.Y) * scalingFactor),
			Z: ((fs[i].A.Z + centering.Z) * scalingFactor),
		}
		b := geometry.Vector{
			X: ((fs[i].B.X + centering.X) * scalingFactor),
			Y: ((fs[i].B.Y + centering.Y) * scalingFactor),
			Z: ((fs[i].B.Z + centering.Z) * scalingFactor),
		}
		c := geometry.Vector{
			X: ((fs[i].C.X + centering.X) * scalingFactor),
			Y: ((fs[i].C.Y + centering.Y) * scalingFactor),
			Z: ((fs[i].C.Z + centering.Z) * scalingFactor),
		}

		a = geometry.Add(rotate(a, rotation), origin)
		b = geometry.Add(rotate(b, rotation), origin)
		c = geometry.Add(rotate(c, rotation), origin)

		triangle := &geometry.Triangle{A: a, B: b, C: c, Properties: props}

		objs = append(objs, triangle)
		trianglesPerCorner[a] = append(trianglesPerCorner[a], triangle)
		trianglesPerCorner[b] = append(trianglesPerCorner[b], triangle)
		trianglesPerCorner[c] = append(trianglesPerCorner[c], triangle)
	}

	for corner, triangles := range trianglesPerCorner {
		for i := 0; i < len(triangles); i++ {
			if triangles[i].NormalsSet {
				continue
			}

			normal := &geometry.Vector{}

			tCount := 0

			// only consider triangles where the angle between normals is > 90°
			for j := 0; j < len(triangles); j++ {
				dot := geometry.Dot(triangles[i].TriangleNormal(), triangles[j].TriangleNormal())

				if dot > 0 || i == j {
					tCount += 1
					normal.X += triangles[j].TriangleNormal().X
					normal.Y += triangles[j].TriangleNormal().Y
					normal.Z += triangles[j].TriangleNormal().Z
				}
			}

			normal.X /= float64(tCount)
			normal.Y /= float64(tCount)
			normal.Z /= float64(tCount)

			if triangles[i].A == corner {
				triangles[i].ASurfaceNormal = *normal
				triangles[i].NormalsSet = true
			}
			if triangles[i].B == corner {
				triangles[i].BSurfaceNormal = *normal
				triangles[i].NormalsSet = true
			}
			if triangles[i].C == corner {
				triangles[i].CSurfaceNormal = *normal
				triangles[i].NormalsSet = true
			}
		}
	}

	return objs, nil
}

func readVector(words []string) (geometry.Vector, error) {
	if len(words) < 4 {
		return geometry.Vector{}, fmt.Errorf("invalid vertex definition (less than 3 elements given)")
	}

	x, err := strconv.ParseFloat(words[1], 64)
	if err != nil {
		return geometry.Vector{}, fmt.Errorf("invalid definition (first element is not a valid number)")
	}

	y, err := strconv.ParseFloat(words[2], 64)
	if err != nil {
		return geometry.Vector{}, fmt.Errorf("invalid definition (second element is not a valid number)")
	}

	z, err := strconv.ParseFloat(words[3], 64)
	if err != nil {
		return geometry.Vector{}, fmt.Errorf("invalid definition (third element is not a valid number)")
	}

	return geometry.Vector{X: x, Y: y, Z: z}, nil
}

func readFace(words []string, vs, vns []geometry.Vector) ([]geometry.Triangle, error) {
	if words[0] != "f" {
		return []geometry.Triangle{}, fmt.Errorf("invalid face definition (line does not start with \"v\")")
	} else if len(words) < 4 {
		return []geometry.Triangle{}, fmt.Errorf("invalid vertex definition (faces must have at least 3 corner vetrices)")
	} else if len(words) > 5 {
		// TODO: implement Seidel's algorithm to support any polygon
		return []geometry.Triangle{}, fmt.Errorf("faces with more than four corners are currently not supported")
	}

	corners := make([]geometry.Vector, 0, 3)
	normals := make([]geometry.Vector, 0, 3)

	for i := 1; i < len(words); i++ {
		cornerSpec := strings.Split(words[i], "/")
		vIdxStr := cornerSpec[0]
		vIdx, err := strconv.ParseInt(vIdxStr, 10, 64)
		if err != nil {
			return []geometry.Triangle{}, fmt.Errorf("invalid face definition (element #%d is not a valid number)", i)
		}

		if vIdx > int64(len(vs)) {
			return []geometry.Triangle{}, fmt.Errorf("invalid face definition (vertex #%d is referenced but not defined)", vIdx)
		} else if int(vIdx) == 0 {
			return []geometry.Triangle{}, fmt.Errorf("invalid face definition (vertex number must be greater than 0)", vIdx)
		}

		if vIdx < 0 {
			vIdx = int64(len(vs)) + vIdx + 1
		}

		corners = append(corners, vs[vIdx-1])

		if len(cornerSpec) >= 3 {
			vnIdxStr := cornerSpec[2]
			vnIdx, err := strconv.ParseInt(vnIdxStr, 10, 64)
			if err != nil {
				return []geometry.Triangle{}, fmt.Errorf("invalid face definition (element #%d is not a valid number)", i)
			}

			if vnIdx > int64(len(vns)) {
				return []geometry.Triangle{}, fmt.Errorf("invalid face definition (vertex normal #%d is referenced but not defined)", vIdx)
			} else if int(vnIdx) == 0 {
				return []geometry.Triangle{}, fmt.Errorf("invalid face definition (vertex normal number must be greater than 0)", vIdx)
			}

			if vnIdx < 0 {
				vnIdx = int64(len(vns)) + vnIdx + 1
			}

			normals = append(normals, vns[vnIdx-1])
		}
	}

	triangles := make([]geometry.Triangle, 0, 1)

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
			t2.BSurfaceNormal = normals[1]
			t2.CSurfaceNormal = normals[2]
			t2.NormalsSet = true
		}

		triangles = append(triangles, t2)
	}

	return triangles, nil
}

func rotate(vec, rotation geometry.Vector) geometry.Vector {
	// rotate around x axis
	xRotated := geometry.Vector{
		X: vec.X,
		Y: vec.Y*math.Cos(rotation.X*math.Pi) - vec.Z*math.Sin(rotation.X*math.Pi),
		Z: vec.Y*math.Sin(rotation.X*math.Pi) + vec.Z*math.Cos(rotation.X*math.Pi),
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

func updateExtremes(minV, maxV *geometry.Vector, newV geometry.Vector) {
	minV.X = math.Min(minV.X, newV.X)
	minV.Y = math.Min(minV.Y, newV.Y)
	minV.Z = math.Min(minV.Z, newV.Z)

	maxV.X = math.Max(maxV.X, newV.X)
	maxV.Y = math.Max(maxV.Y, newV.Y)
	maxV.Z = math.Max(maxV.Z, newV.Z)
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
