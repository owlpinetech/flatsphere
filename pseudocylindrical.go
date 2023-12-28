package flatsphere

import "math"

// An equal-area projection representing the poles as points.
// https://en.wikipedia.org/wiki/Sinusoidal_projection
type Sinusoidal struct{}

func NewSinusoidal() Sinusoidal {
	return Sinusoidal{}
}

func (s Sinusoidal) Project(lat float64, lon float64) (x float64, y float64) {
	return math.Cos(lat) * lon, lat
}

func (s Sinusoidal) Inverse(x float64, y float64) (lat float64, lon float64) {
	return y, x / math.Cos(y)
}

func (s Sinusoidal) Bounds() Bounds {
	return NewRectangleBounds(2*math.Pi, math.Pi)
}
