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
