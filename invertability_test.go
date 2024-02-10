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

func FuzzLambertCylindricalProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewCylindricalEqualArea(0))
}

func FuzzBehrmannProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewBehrmann())
}

func FuzzGallOrthographicProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewGallOrthographic())
}

func FuzzHoboDyerProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewHoboDyer())
}

func FuzzGallStereographicProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewGallStereographic())
}

func FuzzMillerProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewMiller())
}

func FuzzCentralProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewCentral())
}

func FuzzSinusoidalProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewSinusoidal())
}

func FuzzHEALPixStandardProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewHEALPixStandard())
}

//func FuzzMollweideProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewMollweide())
//}

//func FuzzHomolosineProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewHomolosine())
//}

//func FuzzEckertIVProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewEckertIV())
//}

//func FuzzStereographicProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewStereographic())
//}

//func FuzzPolarProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewPolar())
//}

//func FuzzLambertAzimuthalProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewLambertAzimuthal())
//}

//func FuzzTransverseMercatorProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewObliqueProjection(NewMercator(), 0, math.Pi/2, -math.Pi/2))
//}

func FuzzRobinsonProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewRobinson())
}

func FuzzNaturalEarthProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewNaturalEarth())
}

//func FuzzEqualEarthProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewEqualEarth())
//}

func FuzzCassiniProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewCassini())
}

//func FuzzAitoffProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewAitoff())
//}

//func FuzzHammerProjectInverse(f *testing.F) {
//	projectInverseFuzz(f, NewHammer())
//}

func FuzzLagrangeProjectInverse(f *testing.F) {
	projectInverseFuzz(f, NewLagrange())
}

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
		lat = math.Mod(lat, math.Pi/2)
		lon = math.Mod(lon, math.Pi)
		x, y := proj.Project(lat, lon)
		rlat, rlon := proj.Inverse(x, y)

		if withinTolerance(lat, math.Pi/2, 0.0000001) || withinTolerance(lat, -math.Pi/2, 0.0000001) {
			if !withinTolerance(lat, rlat, 0.000001) {
				t.Errorf("expected lat %e, but got %e from %e, %e", lat, rlat, x, y)
			}
		} else {
			if !withinTolerance(lat, rlat, 0.00001) || !withinTolerance(lon, rlon, 0.00001) {
				t.Errorf("expected %e,%e, got %e,%e", lat, lon, rlat, rlon)
			}
		}
	})
}

/*func FuzzObliqueTransformInverse(f *testing.F) {
	f.Add(0.0, 0.0, 0.0, math.Pi/4, -math.Pi/4)
	f.Fuzz(func(t *testing.T, lat float64, lon float64, poleLat float64, poleLon float64, poleTheta float64) {
		if math.Abs(lat) > math.Pi/2 {
			lat = math.Mod(lat, math.Pi/2)
		}
		if math.Abs(lon) > math.Pi {
			lon = math.Mod(lon, math.Pi)
		}
		if math.Abs(poleLat) > math.Pi/2 {
			poleLat = math.Mod(poleLat, math.Pi/2)
		}
		if math.Abs(poleLon) > math.Pi {
			poleLon = math.Mod(poleLon, math.Pi)
		}
		if math.Abs(poleTheta) > math.Pi {
			poleTheta = math.Mod(poleTheta, math.Pi)
		}
		oblique := NewObliqueProjection(nil, poleLat, poleLon, poleTheta)
		x, y := oblique.TransformToOblique(lat, lon)
		rlat, rlon := oblique.TransformFromOblique(x, y)
		if !withinTolerance(lat, rlat, 0.00001) || !withinTolerance(lon, rlon, 0.00001) {
			t.Errorf("expected %e,%e, got %e,%e", lat, lon, rlat, rlon)
		}
	})
}*/
