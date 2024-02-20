package flatsphere

import (
	"fmt"
	"math"
	"testing"
)

type projectTestCase struct {
	lat     float64
	lon     float64
	expectX float64
	expectY float64
}

type inverseTestCase struct {
	x         float64
	y         float64
	expectLat float64
	expectLon float64
}

func checkProject(t *testing.T, name string, proj Projection, cases []projectTestCase) {
	for ind, tc := range cases {
		t.Run(fmt.Sprintf("%s%d", name, ind), func(t *testing.T) {
			x, y := proj.Project(tc.lat, tc.lon)
			if !withinTolerance(x, tc.expectX, 0.000001) || !withinTolerance(y, tc.expectY, 0.000001) {
				t.Errorf("projection at %f,%f expected %e,%e, but got %e,%e", tc.lat, tc.lon, tc.expectX, tc.expectY, x, y)
			}
		})
	}
}

func checkInverse(t *testing.T, name string, proj Projection, cases []inverseTestCase) {
	for ind, tc := range cases {
		t.Run(fmt.Sprintf("%s%d", name, ind), func(t *testing.T) {
			lat, lon := proj.Inverse(tc.x, tc.y)
			if !withinTolerance(lat, tc.expectLat, 0.000001) || !withinTolerance(lon, tc.expectLon, 0.000001) {
				t.Errorf("inverse at %f,%f expected %e,%e, but got %e,%e", tc.x, tc.y, tc.expectLat, tc.expectLon, lat, lon)
			}
		})
	}
}

func TestMercatorProjectSanity(t *testing.T) {
	checkProject(t, "mercator", NewMercator(), []projectTestCase{
		{0, 0, 0, 0},
		{math.Pi / 4, math.Pi / 2, math.Pi / 2, math.Log(math.Tan(math.Pi/4) + 1/math.Cos(math.Pi/4))},
		{-math.Pi / 4, math.Pi / 2, math.Pi / 2, math.Log(math.Tan(-math.Pi/4) + 1/math.Cos(-math.Pi/4))},
		{math.Pi / 4, -math.Pi / 2, -math.Pi / 2, math.Log(math.Tan(math.Pi/4) + 1/math.Cos(math.Pi/4))},
	})
}

func TestMercatorInverseSanity(t *testing.T) {
	checkInverse(t, "invMercator", NewMercator(), []inverseTestCase{
		{0, 0, 0, 0},
		{math.Pi / 2, math.Log(math.Tan(math.Pi/4) + 1/math.Cos(math.Pi/4)), math.Pi / 4, math.Pi / 2},
		{math.Pi / 2, math.Log(math.Tan(-math.Pi/4) + 1/math.Cos(-math.Pi/4)), -math.Pi / 4, math.Pi / 2},
		{-math.Pi / 2, math.Log(math.Tan(math.Pi/4) + 1/math.Cos(math.Pi/4)), math.Pi / 4, -math.Pi / 2},
	})
}

func TestCassiniProjectSanity(t *testing.T) {
	checkProject(t, "cassini", NewCassini(), []projectTestCase{
		{0, 0, 0, 0},
		{0, math.Pi, 0, math.Pi},
		{0, math.Pi / 2, math.Pi / 2, 0},
		{math.Pi / 4, 0, 0, math.Pi / 4},
		{-math.Pi / 4, 0, 0, -math.Pi / 4},
		{math.Pi / 2, 0, 0, math.Pi / 2},
		{-math.Pi / 2, 0, 0, -math.Pi / 2},
		{math.Pi / 2, -math.Pi, 0, math.Pi / 2},
		{-math.Pi / 2, math.Pi, 0, -math.Pi / 2},
	})
}

func TestCassiniInverseSanity(t *testing.T) {
	checkInverse(t, "invCassini", NewCassini(), []inverseTestCase{
		{0, 0, 0, 0},
		{0, math.Pi, 0, math.Pi},
		{math.Pi / 2, 0, 0, math.Pi / 2},
		{0, math.Pi / 4, math.Pi / 4, 0},
		{0, -math.Pi / 4, -math.Pi / 4, 0},
		{0, math.Pi / 2, math.Pi / 2, 0},
		{0, -math.Pi / 2, -math.Pi / 2, 0},
		{math.Pi / 2, math.Pi, 0, math.Pi / 2},
	})
}

func TestEisenlohrProjectSanity(t *testing.T) {
	checkProject(t, "eisenlohr", NewEisenlohr(), []projectTestCase{
		{0, 0, 0, 0},
		{math.Pi / 2, 0, 0, 1 - math.Pi/4},
		{-math.Pi / 2, 0, 0, -(1 - math.Pi/4)},
		{0, math.Pi, math.Sqrt2 + math.Log(math.Sqrt2-1), 0},
		{0, -math.Pi, -(math.Sqrt2 + math.Log(math.Sqrt2-1)), 0},
	})
}

/*func TestTransverseMercatorProjectSanity(t *testing.T) {
	xFrom := func(lat, lon float64) float64 {
		return math.Log((1+math.Sin(lon)*math.Cos(lat))/(1-math.Sin(lon)*math.Cos(lat))) / 2
	}
	yFrom := func(lat, lon float64) float64 {
		return math.Atan(math.Tan(lat) / math.Cos(lon))
	}
	checkProject(t, "transverseMercator", NewObliqueProjection(NewMercator(), 0, math.Pi/2, -math.Pi/2), []projectTestCase{
		{0, 0, 0, 0},
		{math.Pi / 2, 0, 0, math.Pi / 2},
		{math.Pi / 4, 0, 0, math.Pi / 4},
		{-math.Pi / 2, 0, 0, -math.Pi / 2},
		{math.Pi / 4, math.Pi / 2, xFrom(math.Pi/4, math.Pi/2), yFrom(math.Pi/4, math.Pi/2)},
		{-math.Pi / 4, math.Pi / 2, math.Pi / 2, -math.Log(math.Tan(-math.Pi/4) + 1/math.Cos(-math.Pi/4))},
		{math.Pi / 4, -math.Pi / 2, -math.Pi / 2, -math.Log(math.Tan(math.Pi/4) + 1/math.Cos(math.Pi/4))},
	})
}

func TestTransverseMercatorInverseSanity(t *testing.T) {
	checkInverse(t, "invTransverseMercator", NewObliqueProjection(NewMercator(), 0, math.Pi/2, -math.Pi/2), []inverseTestCase{
		{0, 0, 0, 0},
		{-math.Pi / 2, math.Log(math.Tan(math.Pi/4) + 1/math.Cos(math.Pi/4)), math.Pi / 4, math.Pi / 2},
		{math.Pi / 2, -math.Log(math.Tan(-math.Pi/4) + 1/math.Cos(-math.Pi/4)), -math.Pi / 4, math.Pi / 2},
		{-math.Pi / 2, -math.Log(math.Tan(math.Pi/4) + 1/math.Cos(math.Pi/4)), math.Pi / 4, -math.Pi / 2},
	})
}*/
