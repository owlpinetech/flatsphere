package flatsphere

// Represents a rectangular region in arbitrary units, where spherical positions are mapped to the plane.
type Bounds struct {
	XMin float64
	XMax float64
	YMin float64
	YMax float64
}

// The width of the bounds (XMax - XMin)
func (b Bounds) Width() float64 {
	return b.XMax - b.XMin
}

// The height of the bounds (YMax - YMin)
func (b Bounds) Height() float64 {
	return b.YMax - b.YMin
}

// The aspect ratio of the bounds (Width / Height)
func (b Bounds) AspectRatio() float64 {
	return b.Width() / b.Height()
}

// Determines whether the given point is inside the bounding rectangle.
func (b Bounds) Within(x float64, y float64) bool {
	return x >= b.XMin &&
		x <= b.XMax &&
		y >= b.YMin &&
		y <= b.YMax
}

// Construct a new bounding area from the minimum and maximum along each of the two axes.
// The center point will be a position halfway between the minimum and maximum.
func NewBounds(xmin, ymin, xmax, ymax float64) Bounds {
	return Bounds{
		XMin: xmin,
		YMin: ymin,
		XMax: xmax,
		YMax: ymax,
	}
}

// Construct a bounding area containing the circle described by the given
// radius, centered on the origin.
func NewCircleBounds(radius float64) Bounds {
	return NewEllipseBounds(radius, radius)
}

// Construct a bounding area containing the ellipse described by the given
// semiaxes, centered on the origin.
func NewEllipseBounds(semiaxisX float64, semiaxisY float64) Bounds {
	return Bounds{
		XMin: -semiaxisX,
		XMax: semiaxisX,
		YMin: -semiaxisY,
		YMax: semiaxisY,
	}
}

// Construct a bounding area of the given width and height centered on the origin.
func NewRectangleBounds(width float64, height float64) Bounds {
	return Bounds{
		XMin: -width / 2,
		XMax: width / 2,
		YMin: -height / 2,
		YMax: height / 2,
	}
}
