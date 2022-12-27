package quadtree

import (
	"math"
)

type point struct {
	x, y float64
}

type boundary struct {
	minX float64
	minY float64
	maxX float64
	maxY float64
}

func (b *boundary) align() {
	if b.minY > b.maxY {
		b.minY, b.maxY = b.maxY, b.minY
	}
	if b.minX > b.maxX {
		b.minX, b.maxX = b.maxX, b.minX
	}
}

// SubAreas returns four new boundary pointers which represent the four sub-areas
// within this box
func (b *boundary) SubAreas() [4]boundary {
	halfX := b.minX + (b.maxX-b.minX)/2
	halfY := b.minY + (b.maxY-b.minY)/2

	return [4]boundary{
		{b.minX, b.minY, halfX, halfY},
		{halfX, b.minY, b.maxX, halfY},
		{b.minX, halfY, halfX, b.maxY},
		{halfX, halfY, b.maxX, b.maxY},
	}
}

func (b *boundary) Contains(x, y float64) bool {
	return x >= b.minX && x < b.maxX && y >= b.minY && y < b.maxY
}

func (b *boundary) Intersects(aabb boundary) bool {
	b1Width := b.maxX - b.minX
	b2Width := aabb.maxX - aabb.minX

	b1CenterX := b.maxX - (b1Width / 2)
	b2CenterX := aabb.maxX - (b2Width / 2)

	xDistance := math.Abs(b1CenterX-b2CenterX) * 2
	if xDistance > b1Width+b2Width {
		return false
	}

	b1Height := b.maxY - b.minY
	b2Height := aabb.maxY - aabb.minY
	b1CenterY := b.maxY - (b1Height / 2)
	b2CenterY := aabb.maxY - (b2Height / 2)

	yDistance := math.Abs(b1CenterY-b2CenterY) * 2
	if yDistance > b1Height+b2Height {
		return false
	}

	return true
}
