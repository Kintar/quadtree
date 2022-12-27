package quadtree

type point struct {
	x, y float64
}

type axisAlignedBoundingBox struct {
	min, max *point
}

// SubAreas returns four new axisAlignedBoundingBox pointers which represent the four sub-areas
// within this box
func (b *axisAlignedBoundingBox) SubAreas() [4]*axisAlignedBoundingBox {
	halfX := b.max.x / 2
	halfY := b.max.y / 2
	topCenter := &point{x: halfX, y: b.min.y}
	midLeft := &point{x: b.min.x, y: halfY}
	midCenter := &point{x: halfX, y: halfY}
	midRight := &point{x: b.max.x, y: halfY}
	bottomCenter := &point{x: halfX, y: b.max.y}

	return [4]*axisAlignedBoundingBox{
		{min: b.min, max: midCenter},
		{min: topCenter, max: midRight},
		{min: midLeft, max: bottomCenter},
		{min: midCenter, max: b.max},
	}
}

func (b *axisAlignedBoundingBox) Contains(x, y float64) bool {
	return x >= b.min.x && x < b.max.x && y >= b.min.y && y < b.max.y
}

func (b *axisAlignedBoundingBox) Intersects(aabb *axisAlignedBoundingBox) bool {
	return b.Contains(aabb.min.x, aabb.min.y) ||
		b.Contains(aabb.max.x, aabb.max.y) ||
		b.Contains(aabb.min.x, aabb.max.y) ||
		b.Contains(aabb.max.x, aabb.min.y)
}
