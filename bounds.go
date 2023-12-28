package flatsphere

type Bounds struct {
	XMin float64
	XMax float64
	YMin float64
	YMax float64
}

func (b Bounds) Width() float64 {
	return b.XMax - b.XMin
}

func (b Bounds) Height() float64 {
	return b.YMax - b.YMin
}

func (b Bounds) AspectRatio() float64 {
	return b.Width() / b.Height()
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
	return Bounds{}
}
