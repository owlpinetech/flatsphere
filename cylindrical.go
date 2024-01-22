package flatsphere

import (
	"math"
)

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

func (m Mercator) PlanarBounds() Bounds {
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

func (p PlateCarree) PlanarBounds() Bounds {
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

func (e Equirectangular) PlanarBounds() Bounds {
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
	return lon, math.Sin(lat) * l.PlanarBounds().YMax
}

func (l CylindricalEqualArea) Inverse(x float64, y float64) (lat float64, lon float64) {
	return math.Asin(y / l.PlanarBounds().YMax), x
}

func (l CylindricalEqualArea) PlanarBounds() Bounds {
	return NewRectangleBounds(2*math.Pi, 2/l.Stretch)
}

// An equal area projection with least distortion near the equator.
// https://en.wikipedia.org/wiki/Lambert_cylindrical_equal-area_projection
type LambertCylindrical struct {
	CylindricalEqualArea
}

func NewLambertCylindrical() LambertCylindrical {
	return LambertCylindrical{NewCylindricalEqualArea(0)}
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

// A compromise cylindrical projection that tries to minimize distortion as much as possible.
// https://en.wikipedia.org/wiki/Gall_stereographic_projection
type GallStereographic struct{}

func NewGallStereographic() GallStereographic {
	return GallStereographic{}
}

func (g GallStereographic) Project(lat float64, lon float64) (x float64, y float64) {
	return lon, math.Tan(lat/2) * (1 + math.Sqrt(2))
}

func (g GallStereographic) Inverse(x float64, y float64) (lat float64, lon float64) {
	return 2 * math.Atan(y/(1+math.Sqrt(2))), x
}

func (g GallStereographic) PlanarBounds() Bounds {
	return NewRectangleBounds(2*math.Pi, 1.5*math.Pi)
}

// A compromise cylindrical projection intended to resemble Mercator with less distortion at the poles.
// https://en.wikipedia.org/wiki/Miller_cylindrical_projection
type Miller struct{}

func NewMiller() Miller {
	return Miller{}
}

func (m Miller) Project(lat float64, lon float64) (x float64, y float64) {
	return lon, math.Log(math.Tan(math.Pi/4+0.8*lat/2)) / 0.8
}

func (m Miller) Inverse(x float64, y float64) (lat float64, lon float64) {
	return math.Atan(math.Sinh(y*0.8)) / 0.8, x
}

func (m Miller) PlanarBounds() Bounds {
	return NewRectangleBounds(2*math.Pi, 2.5*math.Log(math.Tan(9*math.Pi/20)))
}

// A compromise cylindrical projection with prominent use in panoramic photography, but very distorted
// for mapping purposes.
// https://en.wikipedia.org/wiki/Central_cylindrical_projection
type Central struct{}

func NewCentral() Central {
	return Central{}
}

func (c Central) Project(lat float64, lon float64) (x float64, y float64) {
	return lon, math.Tan(lat)
}

func (c Central) Inverse(x float64, y float64) (lat float64, lon float64) {
	return math.Atan(y), x
}

func (c Central) PlanarBounds() Bounds {
	return NewRectangleBounds(2*math.Pi, 2*math.Pi)
}

// A transverse version of the Plate–Carée projection, implemented directly for efficiency.
// https://en.wikipedia.org/wiki/Cassini_projection
type Cassini struct{}

func NewCassini() Cassini {
	return Cassini{}
}

func (c Cassini) Project(lat float64, lon float64) (float64, float64) {
	x := math.Asin(math.Cos(lat) * math.Sin(lon))
	y := math.Atan2(math.Sin(lat), math.Cos(lat)*math.Cos(lon))
	return x, y
}

func (c Cassini) Inverse(x float64, y float64) (float64, float64) {
	lat := math.Asin(math.Cos(x) * math.Sin(y))
	lon := math.Atan2(math.Sin(x), math.Cos(x)*math.Cos(y))
	return lat, lon
}

func (c Cassini) PlanarBounds() Bounds {
	return NewRectangleBounds(math.Pi, 2*math.Pi)
}
