package flatsphere

import "math"

// A package of functionality for converting from spherical locations to planar coordinates, or from
// planar coordinates into spherical locations. Contains information specifying some characteristics
// of the planar space mapped to by the projection functions, which can differ between projections.
type Projection interface {
	// Convert a location on the sphere (in radians) to a coordinate on the plane.
	Project(lat float64, lon float64) (x float64, y float64)
	// Convert a coordinate on the plane to a location in radians on the sphere.
	Inverse(x float64, y float64) (lat float64, lon float64)
	// Retrieve the planar bounds of the projection.
	PlanarBounds() Bounds
}

// Compute both area distortion and angular distortion at a particular location on the sphere, for the given projection.
func DistortionAt(proj Projection, latitude float64, longitude float64) (area float64, angular float64) {
	nudge := 1e-8
	nudgedLat := latitude + nudge
	// step to the side to avoid interruptions
	pCx, pCy := proj.Project(nudgedLat, longitude)
	// consider a point slightly east
	pEx, pEy := proj.Project(nudgedLat, longitude+nudge/math.Cos(nudgedLat))
	// consider a point slightly north
	pNx, pNy := proj.Project(nudgedLat+nudge, longitude)

	deltaA := (pEx-pCx)*(pNy-pCy) - (pEy-pCy)*(pNx-pCx)
	areaDistortion := math.Log(math.Abs(deltaA / (nudge * nudge)))
	if math.Abs(areaDistortion) > 25.0 {
		areaDistortion = math.NaN()
	}

	s1ps2 := math.Hypot((pEx-pCx)+(pNy-pCy), (pEy-pCy)-(pNx-pCx))
	s1ms2 := math.Hypot((pEx-pCx)-(pNy-pCy), (pEy-pCy)+(pNx-pCx))
	angularDistortion := math.Abs(math.Log(math.Abs((s1ps2 - s1ms2) / (s1ps2 + s1ms2))))
	if angularDistortion > 25.0 {
		angularDistortion = math.NaN()
	}

	return areaDistortion, angularDistortion
}

// Compute the area distortion of a projection at a particular location.
func AreaDistortionAt(proj Projection, latitude float64, longitude float64) float64 {
	nudge := 1e-8
	nudgedLat := latitude + nudge
	// step to the side to avoid interruptions
	pCx, pCy := proj.Project(nudgedLat, longitude)
	// consider a point slightly east
	pEx, pEy := proj.Project(nudgedLat, longitude+nudge/math.Cos(nudgedLat))
	// consider a point slightly north
	pNx, pNy := proj.Project(nudgedLat+nudge, longitude)

	deltaA := (pEx-pCx)*(pNy-pCy) - (pEy-pCy)*(pNx-pCx)
	areaDistortion := math.Log(math.Abs(deltaA / (nudge * nudge)))
	if math.Abs(areaDistortion) > 25.0 {
		areaDistortion = math.NaN()
	}
	return areaDistortion
}

// Compute the angular distortion of a projection at a particular location.
func AngularDistortionAt(proj Projection, latitude float64, longitude float64) float64 {
	nudge := 1e-8
	nudgedLat := latitude + nudge
	// step to the side to avoid interruptions
	pCx, pCy := proj.Project(nudgedLat, longitude)
	// consider a point slightly east
	pEx, pEy := proj.Project(nudgedLat, longitude+nudge/math.Cos(nudgedLat))
	// consider a point slightly north
	pNx, pNy := proj.Project(nudgedLat+nudge, longitude)

	s1ps2 := math.Hypot((pEx-pCx)+(pNy-pCy), (pEy-pCy)-(pNx-pCx))
	s1ms2 := math.Hypot((pEx-pCx)-(pNy-pCy), (pEy-pCy)+(pNx-pCx))
	angularDistortion := math.Abs(math.Log(math.Abs((s1ps2 - s1ms2) / (s1ps2 + s1ms2))))
	if angularDistortion > 25.0 {
		angularDistortion = math.NaN()
	}
	return angularDistortion
}
