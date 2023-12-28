package flatsphere

import (
	"math"
	"testing"
)

func FuzzMercatorProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewMercator())
}

func FuzzPlateCarreeProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewPlateCarree())
}

func FuzzEquirectangularPositiveProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewEquirectangular(45*math.Pi/180))
}

func FuzzEquirectangularNegativeProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewEquirectangular(-45*math.Pi/180))
}

func FuzzLambertCylindrical(f *testing.F) {
	projectInverseFuzz(f, NewCylindricalEqualArea(0))
}

func FuzzBehrmann(f *testing.F) {
	projectInverseFuzz(f, NewBehrmann())
}

func FuzzGallOrthographic(f *testing.F) {
	projectInverseFuzz(f, NewGallOrthographic())
}

func FuzzHoboDyer(f *testing.F) {
	projectInverseFuzz(f, NewHoboDyer())
}

func FuzzGallStereographic(f *testing.F) {
	projectInverseFuzz(f, NewGallStereographic())
}

func FuzzMiller(f *testing.F) {
	projectInverseFuzz(f, NewMiller())
}

func FuzzCentral(f *testing.F) {
	projectInverseFuzz(f, NewCentral())
}

func FuzzSinusoidal(f *testing.F) {
	projectInverseFuzz(f, NewSinusoidal())
}

func FuzzHEALPix(f *testing.F) {
	projectInverseFuzz(f, NewHEALPix())
}

func withinTolerance(n1, n2, tolerance float64) bool {
	if n1 == n2 {
		return true
	}
	diff := math.Abs(n1 - n2)
	if n2 == 0 {
		return diff < tolerance
	}
	return (diff / math.Abs(n2)) < tolerance
}

func projectInverseFuzz(f *testing.F, proj Projection) {
	f.Add(0.0, 0.0)
	f.Add(math.Pi/2, math.Pi/4)
	f.Add(-math.Pi/2, -math.Pi/4)
	f.Fuzz(func(t *testing.T, lat float64, lon float64) {
		x, y := proj.Project(lat, lon)
		rlat, rlon := proj.Inverse(x, y)
		if !withinTolerance(lat, rlat, 0.00001) || !withinTolerance(lon, rlon, 0.00001) {
			t.Errorf("expected %f,%f, got %f,%f", lat, lon, rlat, rlon)
		}
	})
}
