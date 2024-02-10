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
