package flatsphere

import (
	"math"
	"slices"
)

// A class of pseudocylindrical projections defined by a table of ratio values at particular latitutdes.
// Ratio values between the latitudes with defined entries are computed via interpolation during projection/inverse.
// The polynomial order parameter is use to control the interpolation expensiveness/accuracy, and should be a positive
// even number. The y scale parameter is used to control the scaling of latitude projection relative to width, allowing
// the parallel distance ratios to be input in normalized (-1, 1) range (i.e. as an actual ratio).
type TabularProjection struct {
	halfPolynomialOrder int
	yScale              float64
	latitudes           []float64
	parallelLengthRatio []float64
	parallelDistRatio   []float64
}

// Create a new pseudocylindrical projection from a table of values. The tables should be sorted by the latitude
// entry, such that the least latitude, and its corresponding ratio entries, are the first row in the table.
// Polynomial order should be a positive even integer, bigger = more accurate but more compute. yScale should be
// a positive float in the range (0,1).
func NewTabularProjection(
	latitudes []float64,
	parallelLengthRatios []float64,
	parallelDistanceRatios []float64,
	polynomialOrder int,
	yScale float64,
) TabularProjection {
	return TabularProjection{polynomialOrder / 2, yScale, latitudes, parallelLengthRatios, parallelDistanceRatios}
}

var (
	robinsonNaturalEarthLatitudes []float64 = []float64{-90, -85, -80, -75, -70, -65, -60, -55, -50, -45, -40, -35, -30, -25, -20, -15, -10, -5, 0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90}
	robinsonLengthRatios          []float64 = []float64{0.5322, 0.5722, 0.6213, 0.6732, 0.7186, 0.7597, 0.7986, 0.8350, 0.8679, 0.8962, 0.9216, 0.9427, 0.9600, 0.9730, 0.9822, 0.9900, 0.9954, 0.9986, 1.0000, 0.9986, 0.9954, 0.9900, 0.9822, 0.9730, 0.9600, 0.9427, 0.9216, 0.8962, 0.8679, 0.8350, 0.7986, 0.7597, 0.7186, 0.6732, 0.6213, 0.5722, 0.5322}
	robinsonDistRatios            []float64 = []float64{-1.0000, -0.9761, -0.9394, -0.8936, -0.8435, -0.7903, -0.7346, -0.6769, -0.6176, -0.5571, -0.4958, -0.4340, -0.3720, -0.3100, -0.2480, -0.1860, -0.1240, -0.0620, 0.0000, 0.0620, 0.1240, 0.1860, 0.2480, 0.3100, 0.3720, 0.4340, 0.4958, 0.5571, 0.6176, 0.6769, 0.7346, 0.7903, 0.8435, 0.8936, 0.9394, 0.9761, 1.0000}
	naturalEarthLengthRatios      []float64 = []float64{0.5630, 0.6270, 0.6754, 0.7160, 0.7525, 0.7874, 0.8196, 0.8492, 0.8763, 0.9006, 0.9222, 0.9409, 0.9570, 0.9703, 0.9811, 0.9894, 0.9953, 0.9988, 1.0000, 0.9988, 0.9953, 0.9894, 0.9811, 0.9703, 0.9570, 0.9409, 0.9222, 0.9006, 0.8763, 0.8492, 0.8196, 0.7874, 0.7525, 0.7160, 0.6754, 0.6270, 0.5630}
	naturalEarthDistRatios        []float64 = []float64{-1.0000, -0.9761, -0.9394, -0.8936, -0.8435, -0.7903, -0.7346, -0.6769, -0.6176, -0.5571, -0.4958, -0.4340, -0.3720, -0.3100, -0.2480, -0.1860, -0.1240, -0.0620, 0.0000, 0.0620, 0.1240, 0.1860, 0.2480, 0.3100, 0.3720, 0.4340, 0.4958, 0.5571, 0.6176, 0.6769, 0.7346, 0.7903, 0.8435, 0.8936, 0.9394, 0.9761, 1.0000}
)

// Create a new Robinson projection, a well-known instance of a pseudocylindrical projection defined by a table of values.
// https://en.wikipedia.org/wiki/Robinson_projection
func NewRobinson() TabularProjection {
	return NewTabularProjection(robinsonNaturalEarthLatitudes, robinsonLengthRatios, robinsonDistRatios, 4, 0.5072)
}

// Create a new Natural Earth projection, a well-known instance of a pseudocylindrical projection defined by a table of values.
// https://en.wikipedia.org/wiki/Natural_Earth_projection
func NewNaturalEarth() TabularProjection {
	return NewTabularProjection(robinsonNaturalEarthLatitudes, naturalEarthLengthRatios, naturalEarthDistRatios, 4, 0.520)
}

func (t TabularProjection) Project(lat float64, lon float64) (float64, float64) {
	latDegree := lat * math.Pi / 180
	xInterp := t.interpolate(latDegree, t.latitudes, t.parallelLengthRatio)
	yInterp := t.interpolate(latDegree, t.latitudes, t.parallelDistRatio)
	return lon / math.Pi * xInterp, t.yScale * yInterp
}

func (t TabularProjection) Inverse(x float64, y float64) (float64, float64) {
	yNorm := y / t.yScale
	latInterp := t.interpolate(yNorm, t.parallelDistRatio, t.latitudes)
	lonInterp := t.interpolate(yNorm, t.parallelDistRatio, t.parallelLengthRatio)
	return latInterp * 180 / math.Pi, math.Pi * x / lonInterp
}

func (t TabularProjection) interpolate(at float64, xs []float64, ys []float64) float64 {
	ind, _ := slices.BinarySearch(t.latitudes, at)
	return aitkenInterpolation(xs, ys, max(ind-t.halfPolynomialOrder, 0), min(ind+t.halfPolynomialOrder, len(xs)), at)
}

func (t TabularProjection) PlanarBounds() Bounds {
	return NewRectangleBounds(2, t.yScale*2)
}
