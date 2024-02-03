package flatsphere

import "math"

// A compromise azimuthal projection stretched into an elliptical shape.
// https://en.wikipedia.org/wiki/Aitoff_projection
type Aitoff struct{}

func NewAitoff() Aitoff {
	return Aitoff{}
}

func (ai Aitoff) Project(lat float64, lon float64) (float64, float64) {
	a := math.Acos(math.Cos(lat) * math.Cos(lon/2))
	if a == 0 {
		return 0, 0
	}
	x := 2 * math.Cos(lat) * math.Sin(lon/2) * a / math.Sin(a)
	y := math.Sin(lat) * a / math.Sin(a)
	return x, y
}

func (ai Aitoff) Inverse(x float64, y float64) (float64, float64) {
	interLat, interLon := NewPolar().Inverse(x/2, y)
	transLat, transLon := NewObliqueAspect(0, 0, 0).TransformToOblique(interLat, interLon)
	return transLat, transLon * 2
}

func (a Aitoff) PlanarBounds() Bounds {
	return NewEllipseBounds(math.Pi, math.Pi/2)
}

// An elliptical equal-area projection.
// https://en.wikipedia.org/wiki/Hammer_projection
type Hammer struct{}

func NewHammer() Hammer {
	return Hammer{}
}

func (h Hammer) Project(lat float64, lon float64) (float64, float64) {
	z := math.Sqrt(1 + math.Cos(lat)*math.Cos(lon/2))
	x := 2 * math.Cos(lat) * math.Sin(lon/2) / z
	y := math.Sin(lat) / z
	return x, y
}

func (h Hammer) Inverse(x float64, y float64) (float64, float64) {
	z := math.Sqrt(1 - x*x/8 - y*y/2)
	shift := 0.0
	if math.Hypot(x/2, y) > 1 {
		shift = 2 * math.Pi
		if math.Signbit(x) {
			shift = -shift
		}
	}
	lat := math.Asin(z * y * math.Sqrt2)
	lon := 2*math.Atan(math.Sqrt(0.5)*z*x/(2*z*z-1)) + shift
	return lat, lon
}

func (h Hammer) PlanarBounds() Bounds {
	return NewEllipseBounds(2, 1)
}
