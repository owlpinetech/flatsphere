package flatsphere

import (
	"math"
	"testing"
)

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
