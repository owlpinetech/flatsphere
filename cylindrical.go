package flatsphere

import "math"

// An early standard cylindrical projection useful for navigation.
// https://en.wikipedia.org/wiki/Mercator_projection
type Mercator struct{}

func NewMercator() Mercator {
	return Mercator{}
}

func (m Mercator) Project(lat float64, lon float64) (x float64, y float64) {
	return lon, math.Log(math.Tan(math.Pi/4 + lat/2))
}

func (m Mercator) Inverse(x float64, y float64) (lat float64, lon float64) {
	return math.Atan(math.Sinh(y)), x
}

func (m Mercator) Bounds() Bounds {
	return NewRectangleBounds(2*math.Pi, 2*math.Pi)
}
