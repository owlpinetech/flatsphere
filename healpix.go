package flatsphere

import (
	"math"
)

// An equal area projection combining Lambert cylindrical with interrupted Collignon for the polar regions.
// https://en.wikipedia.org/wiki/HEALPix
// See also: "Mapping on the HEALPix Grid", https://arxiv.org/pdf/astro-ph/0412607.pdf
type HEALPixStandard struct{}

func NewHEALPixStandard() HEALPixStandard {
	return HEALPixStandard{}
}

func (h HEALPixStandard) Project(lat float64, lon float64) (x float64, y float64) {
	z := math.Sin(lat)
	if math.Abs(z) <= 2.0/3.0 {
		// equatiorial region
		return lon, 3 * (math.Pi / 8) * z
	} else {
		// polar region
		sigma := math.Sqrt(3 * (1 - math.Abs(z)))
		y := (math.Pi / 4) * (2 - sigma)
		y = math.Copysign(y, lat)

		facetX := (math.Pi / 4) * (2*math.Floor(2+(2*lon)/math.Pi) - 3)
		x := facetX + sigma*(lon-facetX)
		return x, y
	}
}

func (h HEALPixStandard) Inverse(x float64, y float64) (float64, float64) {
	absY := math.Abs(y)
	if absY <= math.Pi/4 {
		// equatorial region
		z := (8 / (3 * math.Pi)) * y
		lat := math.Asin(z)
		return lat, x
	} else if absY < math.Pi/2 {
		// polar region
		sigma := 2 - math.Abs(4*y)/math.Pi
		z := 1 - (sigma * sigma / 3)
		lat := math.Asin(z)
		lat = math.Copysign(lat, y)

		facetX := (math.Pi / 4) * (2*math.Floor(2+(2*x)/math.Pi) - 3)
		lon := facetX + (x-facetX)/sigma
		return lat, lon
	} else {
		return math.Copysign(math.Pi/2, y), x
	}
}

func (h HEALPixStandard) PlanarBounds() Bounds {
	return Bounds{
		XMin: -math.Pi,
		XMax: math.Pi,
		YMin: -math.Pi / 2,
		YMax: math.Pi / 2,
	}
}
