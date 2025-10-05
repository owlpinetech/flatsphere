package flatsphere

import (
	"math"
)

type ObliqueAspect struct {
	poleLat   float64
	poleLon   float64
	poleTheta float64

	sinPoleLat float64 // cached precompute of sin(poleLat)
	cosPoleLat float64 // cached precompute of cos(poleLat)
}

func NewObliqueAspect(poleLat float64, poleLon float64, poleTheta float64) ObliqueAspect {
	return ObliqueAspect{
		poleLat:    poleLat,
		poleLon:    poleLon,
		poleTheta:  poleTheta,
		sinPoleLat: math.Sin(poleLat),
		cosPoleLat: math.Cos(poleLat),
	}
}

// Applies the pole shift and rotation of the Oblique projection transform to the given
// input latitude and longitude points, so that the returned latitude/longitude are able
// to be used for the non-transformed 'original' projection.
func (o ObliqueAspect) TransformFromOblique(latitude float64, longitude float64) (float64, float64) {
	var newLat, newLon float64

	poleRelCos := math.Cos(o.poleLon - longitude)

	// relative latitude
	if o.poleLat == math.Pi/2 {
		newLat = latitude
	} else {
		preAsin := o.sinPoleLat*math.Sin(latitude) + o.cosPoleLat*math.Cos(latitude)*poleRelCos
		if preAsin > 1 && preAsin < 1+1e-9 {
			preAsin = 1
		}
		newLat = math.Asin(preAsin)
	}

	// relative longitude
	switch o.poleLat {
	case math.Pi / 2:
		newLon = longitude - o.poleLon
	case -math.Pi / 2:
		newLon = o.poleLon - longitude - math.Pi
	default:
		numer := o.cosPoleLat*math.Sin(latitude) - o.sinPoleLat*math.Cos(latitude)*poleRelCos
		denom := math.Cos(newLat)
		newLon = math.Acos(numer/denom) - math.Pi

		if math.IsNaN(newLon) {
			if (poleRelCos >= 0 && latitude < o.poleLat) || (poleRelCos < 0 && latitude < -o.poleLat) {
				newLon = 0
			} else {
				newLon = -math.Pi
			}
		} else if math.Sin(longitude-o.poleLon) > 0 {
			newLon = -newLon
		}
	}

	// apply rotate
	newLon = newLon - o.poleTheta

	// constrain and kill roundoff error
	if math.Abs(newLon) > math.Pi {
		newLon = coerceAngle(newLon)
	}
	if newLon >= math.Pi-1e-7 {
		newLon = -math.Pi
	}

	return newLat, newLon
}

func coerceAngle(angle float64) float64 {
	x := angle + math.Pi
	y := 2 * math.Pi
	floorMod := x - math.Floor(x/y)*y
	return floorMod - math.Pi
}

// Given a latitude/longitude in the non-transformed 'original' projection space, applies
// the pole shift and rotation of the Oblique projection so that the returned latitude/longitude
// are in the Oblique projection space.
func (o ObliqueAspect) TransformToOblique(latitude float64, longitude float64) (float64, float64) {
	rotateLon := longitude + o.poleTheta

	preAsin := o.sinPoleLat*math.Sin(latitude) - o.cosPoleLat*math.Cos(latitude)*math.Cos(rotateLon)
	if preAsin > 1 && preAsin < 1+1e-9 {
		preAsin = 1
	}
	if preAsin < -1 && preAsin > -1-1e-9 {
		preAsin = -1
	}
	newLat := math.Asin(preAsin)
	var newLon float64
	inner := math.Sin(latitude)/o.cosPoleLat/math.Cos(newLat) - math.Tan(o.poleLat)*math.Tan(newLat)
	if o.poleLat == math.Pi/2 {
		newLon = rotateLon + o.poleLon
	} else if o.poleLat == -math.Pi/2 {
		newLon = -rotateLon + o.poleLon + math.Pi
	} else if math.Abs(inner) > 1 {
		if (rotateLon == 0 && latitude < -o.poleLat) || (rotateLon != 0 && latitude < o.poleLat) {
			newLon = o.poleLon + math.Pi
		} else {
			newLon = o.poleLon
		}
	} else if math.Sin(rotateLon) > 0 {
		newLon = o.poleLon + math.Acos(inner)
	} else {
		newLon = o.poleLon - math.Acos(inner)
	}

	if math.Abs(newLon) > math.Pi {
		newLon = coerceAngle(newLon)
	}
	return newLat, newLon
}

type ObliqueProjection struct {
	orig Projection
	ObliqueAspect
}

func NewObliqueProjection(original Projection, poleLat float64, poleLon float64, poleTheta float64) ObliqueProjection {
	return ObliqueProjection{
		original,
		NewObliqueAspect(poleLat, poleLon, poleTheta),
	}
}

func (o ObliqueProjection) Project(latitude float64, longitude float64) (float64, float64) {
	return o.orig.Project(o.TransformFromOblique(latitude, longitude))
}

func (o ObliqueProjection) Inverse(x float64, y float64) (float64, float64) {
	return o.TransformToOblique(o.orig.Inverse(x, y))
}

func (o ObliqueProjection) PlanarBounds() Bounds {
	return o.orig.PlanarBounds()
}
