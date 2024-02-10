# flatsphere

[![GoReportCard](https://goreportcard.com/badge/owlpinetech/flatsphere)](https://goreportcard.com/report/github.com/owlpinetech/flatsphere)
![test](https://github.com/owlpinetech/flatsphere/actions/workflows/go.yml/badge.svg)

A [Go](https://go.dev/) library for converting between spherical and planar coordinates.

## Prerequisites

- **[Go](https://go.dev/)**: requires [Go 1.21](https://go.dev/doc/devel/release#go1.21.0) or higher.

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

**Invertability**: all projections in this library have a well defined inverse function. However, only some of them reasonably satisfy the invertability property `x = f_1(f(x))` due to floating-point error in computations, edge-cases resulting in infinities or NaNs, or ambiguity at some type of critical point (commonly the poles or prime meridian). The table below checks of *invertible* for all projections which satisfy the `x = f_1(f(x))` property for all valid spherical locations **except the poles** to a reasonably fine degree of precision, referred to here as **everywhere floating-point invertible**. Oblique transforms of floating-point invertible standard projections do not necessarily share that property.

Efforts are ongoing to improve the coverage of this property to more projections where possible.

|Projection|Everywhere Floating-point Invertible|
|----------|----------|
|Mercator|:white_check_mark:|
|Plate carrée|:white_check_mark:|
|Equirectangular|:white_check_mark:|
|Lambert cylindrical|:white_check_mark:|
|Behrmann|:white_check_mark:|
|Gall orthographic|:white_check_mark:|
|Hobo-Dyer|:white_check_mark:|
|Gall stereographic|:white_check_mark:|
|Miller|:white_check_mark:|
|Central|:white_check_mark:|
|Sinusoidal|:white_check_mark:|
|HEALPix|:white_check_mark:|
|Mollweide| |
|Homolosine| |
|Eckert IV| |
|Stereographic| |
|Polar| |
|Lambert azimuthal| |
|Gnomonic| |
|Orthographic| |
|Robinson|:white_check_mark:|
|Natural Earth|:white_check_mark:|
|Cassini|:white_check_mark:|
|Aitoff| |
|Hammer| |
|Lagrange|:white_check_mark:|

## Credits

Inspired by the [Map-Projections](https://github.com/jkunimune/Map-Projections) library made by [@jkunimune](https://github.com/jkunimune). Right now this library is essentially a Go-flavored port of his original Java, minus a few things to keep it focused on library use cases.
