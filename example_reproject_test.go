package flatsphere

import (
	"fmt"
	"math"
)

// Example_Reproject demonstrates taking a position in one projected plane and converting it into the analogous position
// in a different projection.
func Example_reproject() {
	// determine the original projection of the data
	original := NewMercator()
	// decide on a new desired projection of the data
	target := NewCassini()

	// take a point within the PlanarBounds() of the original projection
	ox, oy := math.Pi, math.Pi
	// get a latitude and longitude on the sphere from the original projected point
	lat, lon := original.Inverse(ox, oy)
	// get a projected point in the desired projection
	cx, cy := target.Project(lat, lon)

	fmt.Printf("original x, y: %f, %f\n", ox, oy)
	fmt.Printf("target x, y: %f, %f\n", cx, cy)

	// Output: original x, y: 3.141593, 3.141593
	// target x, y: 0.000000, 1.657170
}
