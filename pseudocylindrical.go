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

func (s Sinusoidal) PlanarBounds() Bounds {
	return NewRectangleBounds(2*math.Pi, math.Pi)
}

// An equal-area pseudocylindrical map commonly used for maps of the celestial sphere.
// https://en.wikipedia.org/wiki/Mollweide_projection
type Mollweide struct{}

func NewMollweide() Mollweide {
	return Mollweide{}
}

func (m Mollweide) Project(lat float64, lon float64) (x float64, y float64) {
	f := func(t float64) float64 { return 2*t + math.Sin(2*t) }
	d := func(t float64) float64 { return 2 + 2*math.Cos(2*t) }
	theta := newtonsMethod(math.Pi*math.Sin(lat), f, d, 1e-9, 1e-15, 125)
	if math.IsNaN(theta) {
		theta = math.Copysign(math.Pi/2, lat)
	}
	return lon / math.Pi * 2 * math.Cos(theta), math.Sin(theta)
}

func (m Mollweide) Inverse(x float64, y float64) (lat float64, lon float64) {
	theta := math.Asin(y)
	lat = math.Asin((2*theta + math.Sin(2*theta)) / math.Pi)
	lon = x / math.Cos(theta) * math.Pi / 2
	return lat, lon
}

func (m Mollweide) PlanarBounds() Bounds {
	return NewRectangleBounds(2*math.Pi, math.Pi)
}

// An equal-area projection combining Sinusoidal and Mollweide at different hemispheres.
// While most presentations of this projection use an interrupted form, this type is an
// uninterrupted version of Homolosine.
// https://en.wikipedia.org/wiki/Goode_homolosine_projection
type Homolosine struct {
	m     Mollweide
	s     Sinusoidal
	phiH  float64
	scale float64
	yH    float64
}

func NewHomolosine() Homolosine {
	m := NewMollweide()
	_, y := m.Project(0.71098, 0)
	return Homolosine{
		m,
		NewSinusoidal(),
		0.71098,
		math.Sqrt2,
		y * math.Sqrt2,
	}
}

func (h Homolosine) Project(lat float64, lon float64) (x float64, y float64) {
	if math.Abs(lat) <= h.phiH {
		return h.s.Project(lat, lon)
	} else {
		x, y := h.m.Project(lat, lon)
		if lat > 0 {
			return x * h.scale, y*h.scale + h.phiH - h.yH
		} else {
			return x * h.scale, y*h.scale - h.phiH + h.yH
		}
	}
}

func (h Homolosine) Inverse(x float64, y float64) (lat float64, lon float64) {
	if math.Abs(y) <= h.phiH {
		return h.s.Inverse(x, y)
	} else if y > 0 {
		return h.m.Inverse(x/h.scale, (y-h.phiH+h.yH)/h.scale)
	} else {
		return h.m.Inverse(x/h.scale, (y+h.phiH-h.yH)/h.scale)
	}
}

func (h Homolosine) PlanarBounds() Bounds {
	return NewRectangleBounds(2*math.Pi, math.Pi)
}

// An equal-area pseudocylindrical projection, in which the polar lines are half the size of the equator.
// https://en.wikipedia.org/wiki/Eckert_IV_projection
type EckertIV struct{}

func NewEckertIV() EckertIV {
	return EckertIV{}
}

func (e EckertIV) Project(latitude float64, longitude float64) (float64, float64) {
	theta := newtonsMethod(latitude,
		func(t float64) float64 { return t + math.Sin(2*t)/2 + 2*math.Sin(t) },
		func(t float64) float64 { return 1 + math.Cos(2*t) + 2*math.Cos(t) },
		1e-4, 1e-15, 125)
	return longitude / math.Pi * (1 + math.Cos(theta)), math.Sin(theta)
}

func (e EckertIV) Inverse(x float64, y float64) (float64, float64) {
	theta := math.Asin(y)
	latNumer := theta + math.Sin(2*theta)/2 + 2*math.Sin(theta)
	lat := math.Asin(latNumer / (2 + math.Pi/2))
	lon := x / (1 + math.Cos(theta)) * math.Pi
	return lat, lon
}

func (e EckertIV) PlanarBounds() Bounds {
	return NewRectangleBounds(4, 2)
}

// An equal-area pseudocylindrical projection.
// https://en.wikipedia.org/wiki/Equal_Earth_projection
type EqualEarth struct{}

var (
	eeB      float64 = math.Sqrt(3) / 2
	eeYscale float64 = equalEarthPoly(math.Pi/3) / (math.Pi / 3)
)

func NewEqualEarth() EqualEarth {
	return EqualEarth{}
}

func equalEarthPoly(x float64) float64 {
	return 0.003796*math.Pow(x, 9) + 0.000893*math.Pow(x, 7) - 0.081106*math.Pow(x, 3) + 1.340264*x
}

func equalEarthDeriv(x float64) float64 {
	return 9*0.003796*math.Pow(x, 8) + 7*0.000893*math.Pow(x, 6) - 3*0.081106*math.Pow(x, 2) + 1.340264
}

func (e EqualEarth) Project(latitude float64, longitude float64) (float64, float64) {
	theta := math.Asin(math.Sqrt(3) / 2 * math.Sin(latitude))
	return math.Cos(theta) / eeB / equalEarthDeriv(theta) * longitude, equalEarthPoly(theta)
}

func (e EqualEarth) Inverse(x float64, y float64) (float64, float64) {
	theta := newtonsMethod(y/eeYscale, equalEarthPoly, equalEarthDeriv, 1e-6, 1e-15, 125)
	return math.Asin(math.Sin(theta) / eeB), x * eeB / math.Cos(theta) * equalEarthDeriv(theta)
}

func (e EqualEarth) PlanarBounds() Bounds {
	return NewRectangleBounds(2*math.Pi, math.Pi)
}
