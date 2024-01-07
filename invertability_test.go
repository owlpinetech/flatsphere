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

//func FuzzHEALPixStandard(f *testing.F) {
//	projectInverseFuzz(f, NewHEALPixStandard())
//}

//func FuzzMollweide(f *testing.F) {
//	projectInverseFuzz(f, NewMollweide())
//}

//func FuzzStereographic(f *testing.F) {
//	projectInverseFuzz(f, NewStereographic())
//}

//func FuzzPolar(f *testing.F) {
//	projectInverseFuzz(f, NewPolar())
//}

//func FuzzLambertAzimuthal(f *testing.F) {
//	projectInverseFuzz(f, NewLambertAzimuthal())
//}

func withinTolerance(n1, n2, tolerance float64) bool {
	if n1 == n2 {
		return true
	}
	diff := math.Abs(n1 - n2)
	return diff < tolerance
}

func projectInverseFuzz(f *testing.F, proj Projection) {
	f.Add(0.0, 0.0)
	f.Add(0.0, math.Pi)
	f.Add(math.Pi/2, math.Pi/4)
	f.Add(math.Pi/2, 0.0)
	f.Add(-math.Pi/2, -math.Pi/4)
	f.Add(math.Pi/2, math.Pi)
	f.Add(66.0, 0.0)
	f.Fuzz(func(t *testing.T, lat float64, lon float64) {
		if math.Abs(lat) > math.Pi/2 {
			lat = math.Mod(lat, math.Pi/2)
		}
		if math.Abs(lon) > math.Pi {
			lon = math.Mod(lon, math.Pi)
		}
		x, y := proj.Project(lat, lon)
		rlat, rlon := proj.Inverse(x, y)
		if !withinTolerance(lat, rlat, 0.00001) || !withinTolerance(lon, rlon, 0.00001) {
			t.Errorf("expected %e,%e, got %e,%e", lat, lon, rlat, rlon)
		}
	})
}
