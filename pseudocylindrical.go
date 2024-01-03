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
