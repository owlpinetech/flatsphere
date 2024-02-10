package flatsphere

import (
	"math"
)

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
	preAsin := z * y * math.Sqrt2
	if preAsin > 1 && preAsin < 1+1e-9 {
		preAsin = 1
	}
	if preAsin < -1 && preAsin > -1-1e-9 {
		preAsin = -1
	}
	lat := math.Asin(preAsin)
	lon := 2*math.Atan(math.Sqrt(0.5)*z*x/(2*z*z-1)) + shift
	return lat, lon
}

func (h Hammer) PlanarBounds() Bounds {
	return NewEllipseBounds(2, 1)
}

// A circular conformal projection of the whole earth.
type Lagrange struct{}

func NewLagrange() Lagrange {
	return Lagrange{}
}

func (l Lagrange) Project(lat float64, lon float64) (float64, float64) {
	p := (1 + math.Sin(lat)) / (1 - math.Sin(lat))
	v := math.Pow(p, 0.25)
	c := (v+1/v)/2 + math.Cos(lon/2)
	if math.IsInf(c, 0) {
		if math.Signbit(lat) {
			return 0, -1
		} else {
			return 0, 1
		}
	}
	x := math.Sin(lon/2) / c
	y := (v - 1/v) / (2 * c)
	return x, y
}

func (l Lagrange) Inverse(x float64, y float64) (float64, float64) {
	r2 := x*x + y*y
	th := 2 * y / (1 + r2)
	if th == 1 {
		if math.Signbit(y) {
			return -math.Pi / 2, 0
		} else {
			return math.Pi / 2, 0
		}
	}
	t := math.Pow((1+th)/(1-th), 2)
	lat := math.Asin((t - 1) / (t + 1))
	lon := 2 * math.Atan2(2*x, 1-r2)
	return lat, lon
}

func (l Lagrange) PlanarBounds() Bounds {
	return NewCircleBounds(1)
}
