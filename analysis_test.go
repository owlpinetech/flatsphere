package flatsphere

import (
	"math"
	"testing"
)

// Verify that our approximation method applied to the square root function compares
// to the built-in square root function nicely.
func FuzzNewtonsMethodSqrt(f *testing.F) {
	f.Add(612.0)
	f.Add(4.0)
	f.Add(100.0)
	f.Add(2.0)
	f.Add(123.0)
	f.Skip(0.0)
	f.Fuzz(func(t *testing.T, sqrd float64) {
		f := func(x float64) float64 { return x*x - sqrd }
		d := func(x float64) float64 { return 2 * x }
		approx := newtonsMethod(sqrd/2, f, d, 1e-6, 1e-12, 100)
		native := math.Sqrt(sqrd)
		if math.IsNaN(approx) && !math.IsNaN(native) {
			t.Errorf("expected %e, got %e", native, approx)
		} else if !math.IsNaN(approx) && !withinTolerance(approx, native, 0.00001) {
			t.Errorf("expected %e, got %e", native, approx)
		}
	})
}

func FuzzAitkenInterpolation(f *testing.F) {
	xs := []float64{0.0, 1.0, 2.0, 3.0}
	ys := []float64{5.0, 2.0, -5.0, -10.0}
	poly := func(x float64) float64 { return 5.0 + x - 5.0*x*x + x*x*x }
	f.Add(0.0)
	f.Add(1.0)
	f.Add(2.0)
	f.Add(-1.0)
	f.Fuzz(func(t *testing.T, x float64) {
		x = math.Mod(x, xs[len(xs)-1])
		approx := aitkenInterpolation(xs, ys, 0, len(xs), x)
		polyRes := poly(x)
		if !withinTolerance(approx, polyRes, 0.0001) {
			t.Errorf("expected %f, got %f", polyRes, approx)
		}
	})
}
