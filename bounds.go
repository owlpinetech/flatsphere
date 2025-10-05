package flatsphere

type Bounds interface {
	Width() float64
	Height() float64
	Within(x float64, y float64) bool
}

// The aspect ratio of the bounds (Width / Height)
func AspectRatio(b Bounds) float64 {
	return b.Width() / b.Height()
}

// Represents a rectangular region in arbitrary units, where spherical positions are mapped to the plane. Valid
// planar coordinates are within the rectangle defined by (XMin, YMin) and (XMax, YMax).
type RectangleBounds struct {
	XMin float64
	XMax float64
	YMin float64
	YMax float64
}

// Construct a bounding area of the given width and height centered on the origin.
func NewRectangleBounds(width float64, height float64) RectangleBounds {
	return RectangleBounds{
		XMin: -width / 2,
		XMax: width / 2,
		YMin: -height / 2,
		YMax: height / 2,
	}
}

// The width of the bounds (XMax - XMin)
func (b RectangleBounds) Width() float64 {
	return b.XMax - b.XMin
}

// The height of the bounds (YMax - YMin)
func (b RectangleBounds) Height() float64 {
	return b.YMax - b.YMin
}

// Determines whether the given point is inside the bounding rectangle.
func (b RectangleBounds) Within(x float64, y float64) bool {
	return x >= b.XMin &&
		x <= b.XMax &&
		y >= b.YMin &&
		y <= b.YMax
}

// Represents a circular region in arbitrary units, where spherical positions are mapped to the plane. Valid
// planar coordinates are within the circle of the given radius centered on the origin.
type CircleBounds struct {
	Radius float64
}

// Construct a bounding area containing the circle described by the given
// radius, centered on the origin.
func NewCircleBounds(radius float64) CircleBounds {
	return CircleBounds{Radius: radius}
}

func (c CircleBounds) Width() float64 {
	return 2 * c.Radius
}

func (c CircleBounds) Height() float64 {
	return 2 * c.Radius
}

func (c CircleBounds) Within(x float64, y float64) bool {
	return x*x+y*y <= c.Radius*c.Radius
}

type EllipseBounds struct {
	SemiaxisX float64
	SemiaxisY float64
}

// Construct a bounding area containing the ellipse described by the given
// semiaxes, centered on the origin.
func NewEllipseBounds(semiaxisX float64, semiaxisY float64) EllipseBounds {
	return EllipseBounds{SemiaxisX: semiaxisX, SemiaxisY: semiaxisY}
}

func (e EllipseBounds) Width() float64 {
	return 2 * e.SemiaxisX
}

func (e EllipseBounds) Height() float64 {
	return 2 * e.SemiaxisY
}

func (e EllipseBounds) Within(x float64, y float64) bool {
	return (x*x)/(e.SemiaxisX*e.SemiaxisX)+(y*y)/(e.SemiaxisY*e.SemiaxisY) <= 1
}
