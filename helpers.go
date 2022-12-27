package quadtree

import "sort"

func distanceSquared(p1, p2 *point) float64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return dx*dx + dy*dy
}

func sortDataByDistance[T any](origin *point, leafs []*TreeData[T]) []*TreeData[T] {
	pts := dataByDistance[T]{origin, leafs}
	sort.Sort(pts)
	return pts.points
}

type dataByDistance[T any] struct {
	origin *point
	points []*TreeData[T]
}

func (p dataByDistance[T]) Len() int {
	return len(p.points)
}

func (p dataByDistance[T]) Less(i, j int) bool {
	return distanceSquared(p.origin, p.points[i].point) < distanceSquared(p.origin, p.points[i].point)
}

func (p dataByDistance[T]) Swap(i, j int) {
	p.points[i], p.points[j] = p.points[j], p.points[i]
}
