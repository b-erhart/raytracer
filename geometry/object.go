package geometry

type Object interface {
	HitBy(ray *Ray) bool
}
