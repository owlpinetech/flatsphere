# flatsphere

[![GoReportCard](https://goreportcard.com/badge/owlpinetech/flatsphere)](https://goreportcard.com/report/github.com/owlpinetech/flatsphere)

A [Go][] library for converting between spherical and planar coordinates.

## Prerequisites

- **[Go][]**: any one of the **three latest major** [releases][go-releases].

## Install

    go get github.com/owlpinetech/flatsphere

## Usage

    // converting from one projection to another
    origProj := flatsphere.NewMercator()
    newProj := flatsphere.NewLambert()
    lat, lon := origProj.Inverse(origX, origY)
    newX, newY := newProj.Project(lat, lon)

## Credits

Inspired by the [Map-Projections](https://github.com/jkunimune/Map-Projections) library made by [@jkunimune](https://github.com/jkunimune). Right now this library is essentially a Go-flavored port of his original Java, minus a few things to keep it focused on library use cases.
