# flatsphere

[![GoReportCard](https://goreportcard.com/badge/owlpinetech/flatsphere)](https://goreportcard.com/report/github.com/owlpinetech/flatsphere)
![test](https://github.com/owlpinetech/flatsphere/actions/workflows/go.yml/badge.svg)

A [Go](https://go.dev/) library for converting between spherical and planar coordinates.

## Prerequisites

- **[Go](https://go.dev/)**: any one of the **three latest major** [releases](https://go.dev/doc/devel/release).

## Install

    go get github.com/owlpinetech/flatsphere

## Usage

#### Projecting/Inverse

Go from latitude/longitude pairs to points on the x/y plane, or the reverse.

    mercator := flatsphere.NewMercator() // or some other projection
    x, y := mercator.Project(lat, lon)

    rlat, rlon := mercator.Inverse(sampleX, sampleY)

#### Bounds of Projection Plane

Determine the domain of the projection inverse function, to know which valid x/y values can be supplied. Helpful when iterating over the projected space.

    mercator := flatsphere.NewMercator()
    bounds := mercator.PlanarBounds()
    for x := bounds.XMin; x <= bounds.XMax; x += bounds.Width() / xSteps {
        for y := bounds.YMin; y <= bounds.YMax; y += bounds.Height() / ySteps {
            lat, lon := mercator.Inverse(x, y)
        }
    }

#### Reprojecting

Convert planar points in one projection into another projection.

    origProj := flatsphere.NewMercator()
    newProj := flatsphere.NewLambert()
    lat, lon := origProj.Inverse(origX, origY)
    newX, newY := newProj.Project(lat, lon)

#### Oblique Projections

Easily create variants of existing projections with different center points and rotations around the center point.

    mercator := flatsphere.NewMercator() // or some other projection
    transverseMercator := flatsphere.NewOblique(mercator, 0, math.Pi/2, -math.Pi/2)
    x, y := transverseMercator.Project(lat, lon)

#### Distortion

Determine how representative of reality a projection is at a point on the sphere.

    mercator := flatsphere.NewMercator() // or some other projection
    areaDistortion, angularDistortion = proj.DistortionAt(lat, lon)

## Projections

A list of the predefined projections supported by the package.

**Invertability**: all projections in this library have a well defined inverse function. However, only some of them reasonably satisfy the invertability property `x = f_1(f(x))` due to floating-point error in computations, edge-cases resulting in infinities or NaNs, or ambiguity at some type of critical point (commonly the poles or prime meridian). The table below checks of *invertible* for all projections which satisfy the `x = f_1(f(x))` property for all valid spherical locations to a reasonably fine degree of precision, referred to here as **everywhere floating-point invertible**. Oblique transforms of floating-point invertible standard projections do not necessarily share that property.

Efforts are ongoing to improve the coverage of this property to more projections where possible.

|Projection|Everywhere Floating-point Invertible|
|----------|----------|
|Mercator|:white-check-mark:|
|Plate carrée|:white-check-mark:|
|Equirectangular|:white-check-mark:|
|Lambert cylindrical|:white-check-mark:|
|Behrmann|:white-check-mark:|
|Gall orthographic|:white-check-mark:|
|Hobo-Dyer|:white-check-mark:|
|Gall stereographic|:white-check-mark:|
|Miller|:white-check-mark:|
|Central|:white-check-mark:|
|Sinusoidal|:white-check-mark:|
|HEALPix| |
|Mollweide| |
|Stereographic| |
|Polar| |
|Lambert azimuthal| |
|Gnomonic| |
|Orthographic| |

## Credits

Inspired by the [Map-Projections](https://github.com/jkunimune/Map-Projections) library made by [@jkunimune](https://github.com/jkunimune). Right now this library is essentially a Go-flavored port of his original Java, minus a few things to keep it focused on library use cases.
