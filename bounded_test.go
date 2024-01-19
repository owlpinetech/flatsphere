package flatsphere

import (
	"math"
	"testing"
)

func FuzzMercatorProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewMercator())
}

func FuzzPlateCarreeProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewPlateCarree())
}

func FuzzEquirectangularPositiveProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewEquirectangular(45*math.Pi/180))
}

func FuzzEquirectangularNegativeProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewEquirectangular(-45*math.Pi/180))
}

func FuzzLambertCylindricalProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewCylindricalEqualArea(0))
}

func FuzzBehrmannProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewBehrmann())
}

func FuzzGallOrthographicProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewGallOrthographic())
}

func FuzzHoboDyerProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewHoboDyer())
}

func FuzzGallStereographicProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewGallStereographic())
}

func FuzzMillerProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewMiller())
}

func FuzzCentralProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewCentral())
}

func FuzzSinusoidalProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewSinusoidal())
}

func FuzzHEALPixStandardProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewHEALPixStandard())
}

func FuzzMollweideProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewMollweide())
}

func FuzzHomolosineProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewHomolosine())
}

func FuzzEckertIVProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewEckertIV())
}

func FuzzStereographicProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewStereographic())
}

func FuzzPolarProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewPolar())
}

func FuzzLambertAzimuthalProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewLambertAzimuthal())
}

//func FuzzTransverseMercatorProjectBounded(f *testing.F) {
//	projectionBoundedFuzz(f, NewObliqueProjection(NewMercator(), 0, math.Pi/2, -math.Pi/2))
//}

//func FuzzGnomonicProjectBounded(f *testing.F) {
//	projectionBoundedFuzz(f, NewGnomonic())
//}

func FuzzOrthographicProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewOrthographic())
}

func FuzzRobinsonProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewRobinsonProjection())
}

func FuzzNaturalEarthProjectBounded(f *testing.F) {
	projectionBoundedFuzz(f, NewNaturalEarthProjection())
}

func projectionBoundedFuzz(f *testing.F, proj Projection) {
	f.Add(0.0, 0.0)
	f.Add(0.0, math.Pi)
	f.Add(math.Pi/2, math.Pi/4)
	f.Add(math.Pi/2, 0.0)
	f.Add(-math.Pi/2, -math.Pi/4)
	f.Add(math.Pi/2, math.Pi)
	f.Fuzz(func(t *testing.T, lat float64, lon float64) {
		lat = math.Mod(lat, math.Pi/2)
		lon = math.Mod(lon, math.Pi)
		x, y := proj.Project(lat, lon)
		if !proj.PlanarBounds().Within(x, y) {
			t.Errorf("projected point for %e,%e (equal to %e,%e) was outside of the planar bounds for %v", lat, lon, x, y, proj)
		}
	})
}
