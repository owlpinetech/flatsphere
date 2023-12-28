package flatsphere

import (
	"math"
)

// An equal area projection combining Lambert cylindrical with interrupted Collignon for the polar regions.
// https://en.wikipedia.org/wiki/HEALPix
type HEALPix struct{}

func NewHEALPix() HEALPix {
	return HEALPix{}
}

func (h HEALPix) Project(lat float64, lon float64) (x float64, y float64) {
	colatitude := lat + math.Pi/2
	z := math.Cos(colatitude)
	if math.Abs(z) <= 2.0/3.0 {
		// equatiorial region
		return lon, 3 * (math.Pi / 8) * z
	} else {
		// polar region
		facetX := math.Mod(lon, math.Pi/2)
		sigma := 2 - math.Sqrt(3*(1-math.Abs(z)))
		if z < 0 {
			sigma = -sigma
		}
		y := (math.Pi / 4) * sigma
		x := lon - (math.Abs(sigma)-1)*(facetX-math.Pi/4)
		return x, y
	}
}

func (h HEALPix) Inverse(x float64, y float64) (lat float64, lon float64) {
	absY := math.Abs(y)
	//if absY >= math.Pi/2 {
	//	panic(fmt.Sprintf("flatsphere: domain error in projection coordinate y dimension, %v too big", y))
	//}

	if absY <= math.Pi/4 {
		// equatorial region
		z := (8 / (3 * math.Pi)) * y
		colat := math.Acos(z)
		return math.Pi/2 - colat, x
	} else {
		// polar region
		tt := math.Mod(x, math.Pi/2)
		lng := x - ((absY-math.Pi/4)/(absY-math.Pi/2))*(tt-math.Pi/4)
		zz := 2 - 4*absY/math.Pi
		z := (1 - 1.0/3.0*(zz*zz)) * (y / absY)
		colat := math.Acos(z)
		return math.Pi/2 - colat, lng
	}
}

func (h HEALPix) Bounds() Bounds {
	return Bounds{
		XMin: 0,
		XMax: 2 * math.Pi,
		YMin: -math.Pi / 2,
		YMax: math.Pi / 2,
	}
}
