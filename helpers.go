package quadtree

import (
	"math"
	"sort"
)

func pointDistanceSquared(p1, p2 *point) float64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return dx*dx + dy*dy
}

func pointDistance(p1, p2 *point) float64 {
	return math.Sqrt(pointDistanceSquared(p1, p2))
}

func distanceSquared(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}

func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(distanceSquared(x1, y1, x2, y2))
}

func sortDataByDistance[T any](x, y float64, leafs []*TreeLeaf[T]) []*TreeLeaf[T] {
	pts := dataByDistance[T]{&point{x, y}, leafs}
	sort.Sort(pts)
	return pts.points
}

type dataByDistance[T any] struct {
	origin *point
	points []*TreeLeaf[T]
}

func (p dataByDistance[T]) Len() int {
	return len(p.points)
}

func (p dataByDistance[T]) Less(i, j int) bool {
	return pointDistanceSquared(p.origin, p.points[i].point) < pointDistanceSquared(p.origin, p.points[j].point)
}

func (p dataByDistance[T]) Swap(i, j int) {
	p.points[i], p.points[j] = p.points[j], p.points[i]
}
