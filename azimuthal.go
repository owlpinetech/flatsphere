package flatsphere

import (
	"math"
)

// An ancient azimuthal conformal projection (on spheres only, not ellipsoids) which is often used
// for rendering planets to maintain shape of craters. Diverges as latitude approaches Pi/2.
// https://en.wikipedia.org/wiki/Stereographic_map_projection
type Stereographic struct{}

func NewStereographic() Stereographic {
	return Stereographic{}
}

func (s Stereographic) Project(latitude float64, longitude float64) (float64, float64) {
	r := 1 / math.Tan(latitude/2+math.Pi/4)
	return r * math.Sin(longitude), -r * math.Cos(longitude)
}

func (s Stereographic) Inverse(x float64, y float64) (float64, float64) {
	return math.Pi/2 - 2*math.Atan(math.Hypot(x, y)), math.Atan2(x, -y)
}

func (s Stereographic) PlanarBounds() Bounds {
	return NewBounds(math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1))
}

// An ancient equidistant azimuthal projection.
// https://en.wikipedia.org/wiki/Azimuthal_equidistant_projection
type Polar struct{}

func NewPolar() Polar {
	return Polar{}
}

func (p Polar) Project(latitude float64, longitude float64) (float64, float64) {
	r := math.Pi/2 - latitude
	return r * math.Sin(longitude), -r * math.Cos(longitude)
}

func (p Polar) Inverse(x float64, y float64) (float64, float64) {
	return math.Pi/2 - math.Hypot(x, y), math.Atan2(x, -y)
}

func (p Polar) PlanarBounds() Bounds {
	return NewCircleBounds(math.Pi)
}

// An equal-area azimuthal projection.
// https://en.wikipedia.org/wiki/Lambert_azimuthal_equal-area_projection
type LambertAzimuthal struct{}

func NewLambertAzimuthal() LambertAzimuthal {
	return LambertAzimuthal{}
}

func (p LambertAzimuthal) Project(latitude float64, longitude float64) (float64, float64) {
	r := math.Cos((math.Pi/2 + latitude) / 2)
	return r * math.Sin(longitude), -r * math.Cos(longitude)
}

func (p LambertAzimuthal) Inverse(x float64, y float64) (float64, float64) {
	r := math.Hypot(x, y)
	rLat, rLon := math.Asin(1-2*r*r), math.Atan2(x, -y)
	return rLat, rLon
}

func (l LambertAzimuthal) PlanarBounds() Bounds {
	return NewCircleBounds(1)
}

// An ancient projection in which all great cirlces are straight lines. Many use cases
// but rapidly distorts the further away from the center of the projection.
// https://en.wikipedia.org/wiki/Gnomonic_projection
type Gnomonic struct{}

func NewGnomonic() Gnomonic {
	return Gnomonic{}
}

func (g Gnomonic) Project(latitude float64, longitude float64) (float64, float64) {
	r := math.Tan(math.Pi/2 - latitude)
	return r * math.Sin(longitude), -r * math.Cos(longitude)
}

func (g Gnomonic) Inverse(x float64, y float64) (float64, float64) {
	return math.Pi/2 - math.Atan(math.Hypot(x, y)), math.Atan2(x, -y)
}

func (g Gnomonic) PlanarBounds() Bounds {
	return NewRectangleBounds(4, 4)
}

// A projection of a hemisphere of a sphere as if viewed from an infinite distance away.
// https://en.wikipedia.org/wiki/Orthographic_map_projection
type Orthographic struct{}

func NewOrthographic() Orthographic {
	return Orthographic{}
}

func (o Orthographic) Project(latitude float64, longitude float64) (float64, float64) {
	return math.Cos(latitude) * math.Sin(longitude), -math.Cos(latitude) * math.Cos(longitude)
}

func (o Orthographic) Inverse(x float64, y float64) (float64, float64) {
	return math.Acos(math.Hypot(x, y)), math.Atan2(x, -y)
}

func (o Orthographic) PlanarBounds() Bounds {
	return NewCircleBounds(1)
}

// A projection of that mimics the actual appearance of the earth from a fixed viewing distance.
// https://en.wikipedia.org/wiki/General_Perspective_projection
type VerticalPerspective struct {
	D float64 // Distance from the center of the sphere to the viewpoint, in multiples of the sphere radius.
}

func NewVerticalPerspective(d float64) VerticalPerspective {
	if d == 1 {
		panic("D cannot be 1 in VerticalPerspective projection")
	}
	return VerticalPerspective{D: d}
}

func (p VerticalPerspective) Project(latitude float64, longitude float64) (float64, float64) {
	if math.IsInf(p.D, 1) {
		return Orthographic{}.Project(latitude, longitude)
	}
	cosC := math.Cos(latitude) * math.Cos(longitude)
	k := (p.D - 1) / (p.D - cosC)

	x := k * math.Cos(latitude) * math.Sin(longitude)
	y := k * math.Sin(latitude)
	return x, y
}

func (p VerticalPerspective) Inverse(x float64, y float64) (float64, float64) {
	if math.IsInf(p.D, 1) {
		return Orthographic{}.Inverse(x, y)
	}
	r := math.Hypot(x, y)
	if r == 0 {
		return 0, 0
	}
	denom := (p.D-1)/r + r/(p.D-1)
	numer := p.D - math.Sqrt(1-((r*r*(p.D+1))/(p.D-1)))
	c := math.Asin(numer / denom)
	if r > (p.D-1)/p.D {
		c = math.Pi - c
	}
	lat := math.Asin((y * math.Sin(c)) / r)
	lon := math.Atan2(x*math.Sin(c), r*math.Cos(c))
	return lat, lon
}

func (p VerticalPerspective) PlanarBounds() Bounds {
	if math.IsInf(p.D, 1) {
		return Orthographic{}.PlanarBounds()
	}
	return NewCircleBounds(math.Sqrt((p.D - 1) / (p.D + 1)))
}

type ObliqueVerticalPerspective struct {
	CameraLat float64
	CameraLon float64
	D         float64
}

func NewObliqueVerticalPerspective(camLat float64, camLon float64, d float64) ObliqueVerticalPerspective {
	return ObliqueVerticalPerspective{CameraLat: camLat, CameraLon: camLon, D: d}
}

func (o ObliqueVerticalPerspective) Project(latitude float64, longitude float64) (float64, float64) {
	if math.IsInf(o.D, 1) {
		return NewObliqueProjection(NewOrthographic(), o.CameraLat, o.CameraLon, 0).Project(latitude, longitude)
	}
	cosC := math.Sin(o.CameraLat)*math.Sin(latitude) + math.Cos(o.CameraLat)*math.Cos(latitude)*math.Cos(longitude-o.CameraLon)
	k := (o.D - 1) / (o.D - cosC)

	x := k * math.Cos(latitude) * math.Sin(longitude-o.CameraLon)
	y := k * (math.Cos(o.CameraLat)*math.Sin(latitude) - math.Sin(o.CameraLat)*math.Cos(latitude)*math.Cos(longitude-o.CameraLon))
	return x, y
}

func (o ObliqueVerticalPerspective) Inverse(x float64, y float64) (float64, float64) {
	if math.IsInf(o.D, 1) {
		return NewObliqueProjection(NewOrthographic(), o.CameraLat, o.CameraLon, 0).Inverse(x, y)
	}
	r := math.Hypot(x, y)
	if r == 0 {
		return o.CameraLat, o.CameraLon
	}
	denom := (o.D-1)/r + r/(o.D-1)
	numer := o.D - math.Sqrt(1-((r*r*(o.D+1))/(o.D-1)))
	c := math.Asin(numer / denom)
	if r > (o.D-1)/o.D {
		c = math.Pi - c
	}
	lat := math.Asin(math.Cos(c)*math.Sin(o.CameraLat) + (y*math.Sin(c)*math.Cos(o.CameraLat))/r)
	lon := o.CameraLon + math.Atan2(x*math.Sin(c), r*math.Cos(c)*math.Cos(o.CameraLat)-y*math.Sin(c)*math.Sin(o.CameraLat))
	return lat, lon
}

func (p ObliqueVerticalPerspective) PlanarBounds() Bounds {
	if math.IsInf(p.D, 1) {
		return Orthographic{}.PlanarBounds()
	}
	return NewCircleBounds(math.Sqrt((p.D - 1) / (p.D + 1)))
}
