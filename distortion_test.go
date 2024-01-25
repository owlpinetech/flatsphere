package flatsphere

import "testing"

func Test_DistortionNull(t *testing.T) {
	proj := NewMercator()
	eqArea, eqAngular := DistortionAt(proj, 0, 0)

	if !withinTolerance(eqArea, 0.0, 0.000001) || !withinTolerance(eqAngular, 0.0, 0.000001) {
		t.Errorf("expected 0, 0 distortion at mercator equator, got %f, %f", eqArea, eqAngular)
	}

	eqArea = AreaDistortionAt(proj, 0, 0)
	if !withinTolerance(eqArea, 0.0, 0.000001) {
		t.Errorf("expected 0 area distortion at mercator equator, got %f", eqArea)
	}

	eqAngular = AngularDistortionAt(proj, 0, 0)
	if !withinTolerance(eqAngular, 0.0, 0.000001) {
		t.Errorf("expected 0 angular distortion at mercator equator, got %f", eqAngular)
	}
}
