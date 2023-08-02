package geometry

type Sphere struct {
	Center Vector
	Radius float64
}

func (s *Sphere) HitBy(ray *Ray) bool {
	return Cross(ray.Direction.Normalize(), Sub(&s.Center, ray.Origin)).Length() < s.Radius
}
