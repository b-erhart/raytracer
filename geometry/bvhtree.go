package geometry

import (
	"fmt"
	"math"
	"sort"
)

// Bounding Volume Tree
type BvhTree struct {
	root bvhTreeElement
}

func (t BvhTree) String() string {
	return "{root:" + fmt.Sprintf("%v", t.root) + "}"
}

type bvhTreeElement interface {
	getRelevantObjects(ray Ray) []Object
}

type bvhTreeNode struct {
	box   bvhBoundingBox
	left  bvhTreeElement
	right bvhTreeElement
}

func (n *bvhTreeNode) String() string {
	return "{box:" + fmt.Sprintf(
		"%v",
		n.box,
	) + "left:" + fmt.Sprintf(
		"%v",
		n.left,
	) + ", right:" + fmt.Sprintf(
		"%v",
		n.right,
	) + "}"
}

type bvhTreeLeaf struct {
	box  bvhBoundingBox
	objs []Object
}

func (l *bvhTreeLeaf) String() string {
	return "{box:" + fmt.Sprintf("%v", l.box) + ", len(objs):" + fmt.Sprintf("%v", len(l.objs)) + "}"
}

type bvhBoundingBox extremes

func ConstructBvhTree(objs []Object) BvhTree {
	return BvhTree{
		root: constructElement(objs),
	}
}

func constructElement(objs []Object) bvhTreeElement {
	box := calculateBoundingBox(objs)

	if len(objs) <= 4 {
		return &bvhTreeLeaf{
			box:  box,
			objs: objs,
		}
	}

	switch {
	case box.xDiff() >= box.yDiff() && box.xDiff() >= box.zDiff():
		sort.Slice(objs, func(i, j int) bool {
			return objs[i].extremes().maxX < objs[j].extremes().maxX
		})
	case box.yDiff() >= box.xDiff() && box.yDiff() >= box.zDiff():
		sort.Slice(objs, func(i, j int) bool {
			return objs[i].extremes().maxY < objs[j].extremes().maxY
		})
	default:
		sort.Slice(objs, func(i, j int) bool {
			return objs[i].extremes().maxZ < objs[j].extremes().maxZ
		})
	}

	elementsLeft := len(objs) - len(objs)/2

	objsLeft := objs[0:elementsLeft]
	objsRight := objs[elementsLeft:]

	return &bvhTreeNode{
		box:   box,
		left:  constructElement(objsLeft),
		right: constructElement(objsRight),
	}
}

func calculateBoundingBox(objs []Object) bvhBoundingBox {
	if len(objs) == 0 {
		return bvhBoundingBox{}
	}

	box := objs[0].extremes()

	for i := 1; i < len(objs); i++ {
		box = merge(box, objs[i].extremes())
	}

	return bvhBoundingBox(box)
}

func (t BvhTree) GetRelevantObjects(ray Ray) []Object {
	return t.root.getRelevantObjects(ray)
}

func (n *bvhTreeNode) getRelevantObjects(ray Ray) []Object {
	if !n.box.intersects(ray) {
		return []Object{}
	}

	objsLeft := n.left.getRelevantObjects(ray)
	objsRight := n.right.getRelevantObjects(ray)

	ret := append(objsLeft[:len(objsLeft):len(objsLeft)], objsRight...)

	return ret
}

func (l *bvhTreeLeaf) getRelevantObjects(ray Ray) []Object {
	if !l.box.intersects(ray) {
		return []Object{}
	}

	return l.objs
}

// source: https://tavianator.com/2011/ray_box.html
func (b bvhBoundingBox) intersects(ray Ray) bool {
	inverseRayDir := Vector{
		X: 1 / ray.Direction.Normalize().X,
		Y: 1 / ray.Direction.Normalize().Y,
		Z: 1 / ray.Direction.Normalize().Z,
	}

	tx1 := (b.minX - ray.Origin.X) * inverseRayDir.X
	tx2 := (b.maxX - ray.Origin.X) * inverseRayDir.X

	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)

	ty1 := (b.minY - ray.Origin.Y) * inverseRayDir.Y
	ty2 := (b.maxY - ray.Origin.Y) * inverseRayDir.Y

	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))

	tz1 := (b.minZ - ray.Origin.Z) * inverseRayDir.Z
	tz2 := (b.maxZ - ray.Origin.Z) * inverseRayDir.Z

	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))

	return tmax >= math.Max(0, tmin)
}

func (b bvhBoundingBox) xDiff() float64 {
	return b.maxX - b.minX
}

func (b bvhBoundingBox) yDiff() float64 {
	return b.maxY - b.minY
}

func (b bvhBoundingBox) zDiff() float64 {
	return b.maxZ - b.minZ
}
