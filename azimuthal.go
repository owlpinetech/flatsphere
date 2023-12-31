package flatsphere

import (
	"fmt"
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
	return NewRectangleBounds(4, 4)
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
	if math.IsNaN(rLat) {
		fmt.Printf("nan lat: r = %e, x = %e, y = %e, presin = %e, valid = %v\n", r, x, y, 1-2*r*r, 1-2*r*r < -1)
	}
	return rLat, rLon
}

func (l LambertAzimuthal) PlanarBounds() Bounds {
	return NewCircleBounds(1)
}
