package flatsphere

import "math"

// Try to find a root of the given target function. Returns NaN if the process fails.
func newtonsMethod(
	initialGuess float64,
	targetFunc func(float64) float64,
	derivative func(float64) float64,
	tolerance float64,
	epsilon float64,
	maxIterations int,
) float64 {
	currentGuess := initialGuess
	for i := 0; i < maxIterations; i++ {
		// get function/derivative values for current guess
		y := targetFunc(currentGuess)
		yPrime := derivative(currentGuess)

		// if denominator is too small error out
		if math.Abs(yPrime) < epsilon {
			break
		}

		nextGuess := currentGuess - y/yPrime
		if math.Abs(nextGuess-currentGuess) < tolerance {
			return nextGuess
		}
		currentGuess = nextGuess
	}
	return math.NaN()
}

// Given a table of x,y values representing points on a function, Aitken interpolation approximates
// the y value of the described function at a given x value, estimating using the points in the table
// between xStart and xEnd.
func aitkenInterpolation(xs []float64, ys []float64, xStart int, xEnd int, atX float64) float64 {
	samples := xEnd - xStart
	approximations := make([][]float64, samples)

	approximations[0] = ys[xStart:xEnd] // fill in the zero row

	// iterate recursively, forming a triangular matrix
	for i := 1; i < samples; i++ {
		approximations[i] = make([]float64, samples)
		for j := i; j < samples; j++ {
			determ := determinant(approximations[i-1][i-1], approximations[i-1][j], xs[xStart+i-1]-atX, xs[xStart+j]-atX)
			denom := xs[xStart+j] - xs[xStart+i-1]
			approximations[i][j] = determ / denom
		}
	}

	// the last value is the final interpolated value
	return approximations[samples-1][samples-1]
}

func determinant(a, b, c, d float64) float64 {
	return a*d - b*c
}
