// Package quadtree implements a quadtree for storing arbitrary data.
package quadtree

import "sort"

type TreeData[T any] struct {
	data T
	*point
}

type QuadTree[T any] struct {
	branchSize int
	capacity   int
	boundary   *axisAlignedBoundingBox
	depth      int
	data       []*TreeData[T]
	parent     *QuadTree[T]
	nodes      [4]*QuadTree[T]
}

func NewQuadTree[T any](x1, y1, x2, y2 float64) *QuadTree[T] {
	return &QuadTree[T]{
		capacity: 8,
		boundary: &axisAlignedBoundingBox{
			min: &point{x1, y1},
			max: &point{x2, y2},
		},
	}
}

func (q *QuadTree[T]) divide() {
	subBounds := q.boundary.SubAreas()
	// Iterate over the new bounds and create trees
	for idx, bound := range subBounds {
		node := NewQuadTree[T](bound.min.x, bound.min.y, bound.max.x, bound.max.y)
		node.parent = q
		q.nodes[idx] = node
		// Iterate over our non-nil parent data and see if they belong in this node
		for i, datum := range q.data {
			if datum != nil && node.boundary.Contains(datum.x, datum.y) {
				node.data = append(node.data, datum)
				// Since we assigned it to a sub-node, no need to consider it again
				q.data[i] = nil
			}
		}
	}

	q.data = nil
}

func (q *QuadTree[T]) Len() int {
	return q.branchSize
}

func (q *QuadTree[T]) Insert(x, y float64, data T) {
	// Is the data within our boundary? If not, just ignore it.
	if !q.boundary.Contains(x, y) {
		return
	}

	// Are we a leaf node?
	if q.nodes[0] == nil {
		q.branchSize++
		// Yes, but do we have capacity?
		if len(q.data) == q.capacity {
			// We are at capacity, divide
			q.divide()
		} else {
			// we are NOT at capacity. Slot it in
			q.data = append(q.data, &TreeData[T]{
				data:  data,
				point: &point{x, y},
			})
		}
	}

	// If we got here, we are NOT a leaf node. Ask our kids to store the data
	for i := 0; i < 4; i++ {
		q.nodes[i].Insert(x, y, data)
	}
}

func (q *QuadTree[T]) leafContaining(x, y float64) *QuadTree[T] {
	if q.boundary.Contains(x, y) {
		if q.nodes[0] == nil {
			// We're a leaf node, it's us!
			return q
		}
		// We're not a leaf, so see if it's in any of our kids
		for _, node := range q.nodes {
			result := node.leafContaining(x, y)
			if result != nil {
				return result
			}
		}
	}

	return nil
}

func (q *QuadTree[T]) collectChildren() []*TreeData[T] {
	data := make([]*TreeData[T], 0, q.branchSize)
	data = append(data, q.data...)
	if q.nodes[0] != nil {
		for _, node := range q.nodes {
			data = append(node.collectChildren())
		}
	}
	return data
}

// FindNearest returns up to count elements, sorted by distance to (x,y)
func (q *QuadTree[T]) FindNearest(x, y float64, count int) []*TreeData[T] {
	leafNode := q.leafContaining(x, y)
	if leafNode == nil {
		return nil
	}
	// walk up the tree until we find a node with at least count children, or we run out of tree
	for leafNode.branchSize <= count && leafNode.parent != nil {
		leafNode = leafNode.parent
	}
	// Sort 'em
	data := dataByDistance[T]{
		points: leafNode.collectChildren(),
		origin: &point{x, y},
	}
	sort.Sort(data)

	// if we got too many, slice 'em up
	if len(data.points) > count {
		data.points = data.points[:count]
	}
	
	return data.points
}
