package flatsphere

import (
	"math"
	"testing"
)

func FuzzObliqueNoOp(f *testing.F) {
	obliqEq := NewObliqueProjection(NewPlateCarree(), math.Pi/2, 0, 0)
	cassini := NewPlateCarree()

	f.Add(0.0, 0.0)
	f.Add(math.Pi/4, math.Pi/4)
	f.Add(math.Pi/2, 0.0)
	f.Add(0.0, math.Pi/2)
	f.Fuzz(func(t *testing.T, lat float64, lon float64) {
		if math.Abs(lat) > math.Pi/2 {
			lat = math.Mod(lat, math.Pi/2)
		}
		if math.Abs(lon) > math.Pi {
			lon = math.Mod(lon, math.Pi)
		}

		xo, yo := obliqEq.Project(lat, lon)
		xc, yc := cassini.Project(lat, lon)

		if !withinTolerance(xo, xc, 0.000001) || !withinTolerance(yo, yc, 0.000001) {
			t.Errorf("expected %e,%e, but got %e,%e", xc, yc, xo, yo)
		}
	})
}

func rotatePoint(x, y float64, rad float64) (float64, float64) {
	rotX := x*math.Cos(rad) - y*math.Sin(rad)
	rotY := x*math.Sin(rad) + y*math.Cos(rad)
	return rotX, rotY
}

func FuzzObliquePlateCarreeVsCassiniProject(f *testing.F) {
	obliqEq := NewObliqueProjection(NewPlateCarree(), 0, math.Pi/2, -math.Pi/2)
	cassini := NewCassini()

	f.Add(0.0, 0.0)
	f.Add(math.Pi/4, math.Pi/4)
	f.Add(math.Pi/2, 0.0)
	f.Add(0.0, math.Pi/2)
	f.Fuzz(func(t *testing.T, lat float64, lon float64) {
		if math.Abs(lat) > math.Pi/2 {
			lat = math.Mod(lat, math.Pi/2)
		}
		if math.Abs(lon) > math.Pi {
			lon = math.Mod(lon, math.Pi)
		}

		xo, yo := obliqEq.Project(lat, lon)
		xc, yc := cassini.Project(lat, lon)
		xr, yr := rotatePoint(xc, yc, math.Pi/2)

		if !withinTolerance(xo, xr, 0.000001) || !withinTolerance(yo, yr, 0.000001) {
			t.Errorf("expected %e,%e, but got %e,%e", xr, yr, xo, yo)
		}
	})
}
