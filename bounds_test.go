package flatsphere

import (
	"testing"
)

func TestBoundsProperties(t *testing.T) {
	testCases := []struct {
		name   string
		bounds Bounds
		width  float64
		height float64
	}{
		{"Circle", NewCircleBounds(1.0), 2.0, 2.0},
		{"Ellipse", NewEllipseBounds(2.0, 3.0), 4.0, 6.0},
		{"Rectangle", NewRectangleBounds(5.0, 1.0), 5.0, 1.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !withinTolerance(tc.bounds.Width(), tc.width, 0.000001) {
				t.Errorf("expected width %f, got %f", tc.width, tc.bounds.Width())
			}
			if !withinTolerance(tc.bounds.Height(), tc.height, 0.000001) {
				t.Errorf("expected height %f, got %f", tc.height, tc.bounds.Height())
			}
		})
	}
}

func TestWithin(t *testing.T) {
	testCases := []struct {
		name   string
		bounds Bounds
		xloc   float64
		yloc   float64
		within bool
	}{
		{"Origin", NewRectangleBounds(2.0, 2.0), 0.0, 0.0, true},
		{"OffsetInside", NewRectangleBounds(2.0, 2.0), 0.5, 0.5, true},
		{"NegativeOffsetInside", NewRectangleBounds(2.0, 2.0), -0.5, -0.5, true},
		{"Edge", NewRectangleBounds(2.0, 2.0), 1.0, 0.0, true},
		{"Corner", NewRectangleBounds(2.0, 2.0), 1.0, 1.0, true},
		{"NegativeCorner", NewRectangleBounds(2.0, 2.0), -1.0, -1.0, true},
		{"Outside", NewRectangleBounds(2.0, 2.0), 3.0, 3.0, false},
		{"NegativeOutside", NewRectangleBounds(2.0, 2.0), -3.0, -3.0, false},
		{"XAxisOutside", NewRectangleBounds(2.0, 2.0), 3.0, 0.0, false},
		{"YAxisOutside", NewRectangleBounds(2.0, 2.0), 0.0, 3.0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.bounds.Within(tc.xloc, tc.yloc) != tc.within {
				t.Errorf("expected %v, got %v", tc.within, tc.bounds.Within(tc.xloc, tc.yloc))
			}
		})
	}
}
