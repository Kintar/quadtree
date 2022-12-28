package model

import "math"

// BoundingSquare defines a square region in a 2d plane
type BoundingSquare struct {
	// (cx, cy) is the center of the bounding box
	cx float64
	cy float64
	// extent is the distance this box extends along each axis from the center point
	extent float64
	// size is the total length of the sides of the square
	size float64
}

// Contains returns true if the specified point lies within the bounding box
func (bb *BoundingSquare) Contains(x, y float64) bool {
	return x >= bb.cx-bb.extent && x < bb.cx+bb.extent && y >= bb.cy-bb.extent && y < bb.cy+bb.extent
}

// Intersects returns true if the two boundary squares overlap
func (bb *BoundingSquare) Intersects(b2 BoundingSquare) bool {
	maxDistance := bb.size + b2.size
	return math.Abs(bb.cx-b2.cx)*2 <= maxDistance &&
		math.Abs(bb.cy-b2.cy)*2 <= maxDistance
}

func (bb *BoundingSquare) IntersectsCircle(c BoundingCircle) bool {
	return c.IntersectsSquare(*bb)
}

// NewBoundSquare creates a new bounding square centered at (cx,cy) with a side length equal to size
func NewBoundSquare(cx, cy, size float64) BoundingSquare {
	size = math.Abs(size)
	return BoundingSquare{
		cx:     cx,
		cy:     cy,
		extent: size / 2,
		size:   size,
	}
}

// NewBoundSquareFromCorner create a new bounding square with the one corner at x,y and opposite corner
// at x+size,y+size
func NewBoundSquareFromCorner(x, y, size float64) BoundingSquare {
	size = math.Abs(size)
	extent := size / 2
	return BoundingSquare{
		cx:     x + extent,
		cy:     y + extent,
		extent: extent,
		size:   size,
	}
}

type BoundingCircle struct {
	cx         float64
	cy         float64
	radius     float64
	radSquared float64
}

// NewBoundingCircle creates a new bounding circle centered at (cx,cy) with a given radius
func NewBoundingCircle(cx, cy, radius float64) BoundingCircle {
	return BoundingCircle{
		cx:         cx,
		cy:         cy,
		radius:     radius,
		radSquared: radius * radius,
	}
}

func distSquared(x1, y1, x2, y2 float64) float64 {
	xx := x1 - x2
	yy := y1 - y2
	return xx*xx + yy*yy
}

// Contains returns true if the point (x,y) lies within the bounding circle
func (c *BoundingCircle) Contains(x, y float64) bool {
	return distSquared(x, y, c.cx, c.cy) < c.radSquared
}

func (c *BoundingCircle) Intersects(c2 BoundingCircle) bool {
	totalRad := c.radius + c2.radius
	return math.Abs(c.cx-c2.cx) < totalRad && math.Abs(c.cy-c2.cy) < totalRad
}

// IntersectsSquare returns true if this circle intersects the bounding square
func (c *BoundingCircle) IntersectsSquare(b BoundingSquare) bool {
	// Distance between the center coordinates of both boundaries
	cdX, cdY := math.Abs(c.cx-b.cx), math.Abs(c.cy-b.cy)

	// If either component of the distance is greater than the combined radius and extent of the bounds, they cannot
	// intersect
	if cdX > b.extent+c.radius || cdY > b.extent+c.radius {
		return false
	}

	// Since neither component is too far for an intersection, then if either component is closer than the extent of
	// the bounding square, then they intersect
	if cdX < b.extent || cdY < b.extent {
		return true
	}

	// Last case. All four corners of the square lie the same distance from its center, so using absolute values allows
	// us to use a single check to tell if any corner of the bounding box lies within the circle
	centerDistanceSquared := distSquared(cdX, cdY, b.extent, b.extent)
	bbdX, bbdY := math.Abs(b.cx-b.extent), math.Abs(b.cy-b.extent)
	return centerDistanceSquared < bbdX*bbdX+bbdY*bbdY
}

func Subdivide(bb BoundingSquare) [4]BoundingSquare {
	halfExtent := bb.extent / 2
	return [4]BoundingSquare{
		NewBoundSquare(bb.cx-halfExtent, bb.cy-halfExtent, bb.extent),
		NewBoundSquare(bb.cx+halfExtent, bb.cy-halfExtent, bb.extent),
		NewBoundSquare(bb.cx-halfExtent, bb.cy+halfExtent, bb.extent),
		NewBoundSquare(bb.cx+halfExtent, bb.cy+halfExtent, bb.extent),
	}
}
