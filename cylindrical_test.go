package flatsphere

import (
	"math"
	"testing"
)

func TestEqualAreaStretch(t *testing.T) {
	testCases := []struct {
		name    string
		lat     float64
		stretch float64
	}{
		{"Lambert", 0, 1},
		{"GallOrtho", math.Pi / 4, 1.0 / 2.0},
		{"NegGallOrtho", -math.Pi / 4, 1.0 / 2.0},
		{"NorthPole", math.Pi / 2, 0},
		{"WrappedAround", math.Pi, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			proj := NewCylindricalEqualArea(tc.lat)
			if !withinTolerance(proj.Stretch, tc.stretch, 0.000001) {
				t.Errorf("expected streth factor %f, got %f", tc.stretch, proj.Stretch)
			}
		})
	}
}

func FuzzEqualAreaCtorValidity(f *testing.F) {
	f.Add(0.0)
	f.Add(math.Pi / 4)
	f.Add(-math.Pi / 4)
	f.Add(math.Pi / 2)
	f.Add(math.Pi)
	f.Fuzz(func(t *testing.T, lat float64) {
		proj := NewCylindricalEqualArea(lat)
		if proj.Stretch > 1 {
			t.Errorf("managed to create a cylindrical equal area projection with invalid streth factor %f", proj.Stretch)
		}
	})
}
