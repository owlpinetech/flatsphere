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

// A special case of the equirectangular projection which allows for easy conversion between
// pixel coordinates and locations on the sphere. The scale is less distorted the closer toward
// the equator a spherical position is.
// https://en.wikipedia.org/wiki/Equirectangular_projection
type PlateCarree struct{}

func NewPlateCarree() PlateCarree {
	return PlateCarree{}
}

func (p PlateCarree) Project(lat float64, lon float64) (x float64, y float64) {
	return lon, lat
}

func (p PlateCarree) Inverse(x float64, y float64) (lat float64, lon float64) {
	return y, x
}

func (p PlateCarree) Bounds() Bounds {
	return NewRectangleBounds(2*math.Pi, math.Pi)
}

// A linear mapping of latitude and longitude to x and y, with the given latitude of focus
// where results will be undistorted.
// https://en.wikipedia.org/wiki/Equirectangular_projection
type Equirectangular struct {
	// The latitude (in radians) at which the scale is true (undistorted) for this projection.
	Parallel float64
}

// Construct a new equirectangular projection with less distortion at the given parellel (in radians).
func NewEquirectangular(parallel float64) Equirectangular {
	return Equirectangular{Parallel: parallel}
}

func (e Equirectangular) Project(lat float64, lon float64) (x float64, y float64) {
	return lon, lat / math.Cos(e.Parallel)
}

func (e Equirectangular) Inverse(x float64, y float64) (lat float64, lon float64) {
	return y * math.Cos(e.Parallel), x
}

func (e Equirectangular) Bounds() Bounds {
	return NewRectangleBounds(2*math.Pi, math.Pi/math.Cos(e.Parallel))
}

// A generalized form of equal-area cylindrical projection.
// https://en.wikipedia.org/wiki/Cylindrical_equal-area_projection
type CylindricalEqualArea struct {
	// The stretch factor of the resulting projection based on the desired parallel of least distortion.
	// If Stretch == 1, then the parallel of least distortion is the equator. If Stretch > 1, then there
	// are no parallels where the horizontal scale matches the vertical scale.
	Stretch float64
}

// Construct a new cylindrical equal-area projection with least distortion around the given latitude in radians.
func NewCylindricalEqualArea(parallel float64) CylindricalEqualArea {
	stretch := math.Pow(math.Cos(parallel), 2)
	return CylindricalEqualArea{Stretch: stretch}
}

// The northern latitude (in radians) at which the map has the least distortion.
func (l CylindricalEqualArea) Parallel() float64 {
	return math.Acos(math.Sqrt(l.Stretch))
}

func (l CylindricalEqualArea) Project(lat float64, lon float64) (x float64, y float64) {
	return lon, math.Sin(lat) * l.Bounds().YMax
}

func (l CylindricalEqualArea) Inverse(x float64, y float64) (lat float64, lon float64) {
	return math.Asin(y / l.Bounds().YMax), x
}

func (l CylindricalEqualArea) Bounds() Bounds {
	return NewRectangleBounds(2*math.Pi, 2/l.Stretch)
}

// An equal area projection with least distortion near the equator.
// https://en.wikipedia.org/wiki/Lambert_cylindrical_equal-area_projection
type Lambert struct {
	CylindricalEqualArea
}

func NewLambert() Lambert {
	return Lambert{NewCylindricalEqualArea(0)}
}

// An equal area projection with least distortion at 30 degrees latitude.
// https://en.wikipedia.org/wiki/Behrmann_projection
type Behrmann struct {
	CylindricalEqualArea
}

func NewBehrmann() Behrmann {
	return Behrmann{NewCylindricalEqualArea(math.Pi / 6)}
}

// An equal area projection with least distortion at 45 degrees latitude.
// https://en.wikipedia.org/wiki/Gall%E2%80%93Peters_projection
type GallOrthographic struct {
	CylindricalEqualArea
}

func NewGallOrthographic() GallOrthographic {
	return GallOrthographic{NewCylindricalEqualArea(45 * math.Pi / 180)}
}

// An equal area projection with least distortion at 37.5 degrees latitude.
// https://en.wikipedia.org/wiki/Hobo%E2%80%93Dyer_projection
type HoboDyer struct {
	CylindricalEqualArea
}

func NewHoboDyer() HoboDyer {
	return HoboDyer{NewCylindricalEqualArea(37.5 * math.Pi / 180)}
}
