package flatsphere

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
