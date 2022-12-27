// Package quadtree implements a quadtree for storing arbitrary data.
package quadtree

type TreeLeaf[T any] struct {
	Content T
	*point
}

func (td *TreeLeaf[T]) Location() (float64, float64) {
	return td.x, td.y
}

type QuadTree[T any] struct {
	branchSize int
	capacity   int
	boundary   boundary
	depth      int
	data       []*TreeLeaf[T]
	parent     *QuadTree[T]
	nodes      [4]*QuadTree[T]
}

func NewQuadTree[T any](x1, y1, x2, y2 float64) *QuadTree[T] {
	bound := boundary{
		x1, y1, x2, y2,
	}
	bound.align()

	return &QuadTree[T]{
		capacity: 8,
		boundary: bound,
	}
}

// NodeCapacity updates the capacity of the tree to the new value.
// NodeCapacity can be called on an individual sub-branch of the tree, and will only affect that branch and
// its children.
// Setting a node's capacity lower than its current size will not cause the leaf to divide until a new insertion is made
func (q *QuadTree[T]) NodeCapacity(capacity int) {
	q.capacity = capacity
	if q.nodes[0] != nil {
		for _, node := range q.nodes {
			node.NodeCapacity(capacity)
		}
	}
}

func (q *QuadTree[T]) divide() {
	subBounds := q.boundary.SubAreas()
	// Iterate over the new bounds and create trees
	for idx, bound := range subBounds {
		node := NewQuadTree[T](bound.minX, bound.minY, bound.maxX, bound.maxY)
		node.capacity = q.capacity
		node.parent = q
		node.depth = q.depth + 1
		q.nodes[idx] = node
		for _, data := range q.data {
			if node.boundary.Contains(data.x, data.y) {
				node.data = append(node.data, data)
				node.branchSize++
			}
		}
	}

	q.data = nil
}

func (q *QuadTree[T]) Len() int {
	return q.branchSize
}

func (q *QuadTree[T]) Insert(x, y float64, data T) bool {
	return q.insert(&TreeLeaf[T]{
		data,
		&point{x, y},
	})
}

func (q *QuadTree[T]) insert(leaf *TreeLeaf[T]) bool {
	// Is the Data within our boundary? If not, just ignore it.
	if !q.boundary.Contains(leaf.x, leaf.y) {
		return false
	}

	// If we're a leaf node...
	if q.nodes[0] == nil {
		// Are we maxed out?
		if len(q.data) >= q.capacity {
			// we are, split
			q.divide()
		} else {
			// we're not, add to our own list and return
			q.data = append(q.data, leaf)
			q.branchSize++
			return true
		}
	}

	// If we got here, we have children to do this task for us!
	for _, node := range q.nodes {
		if node.insert(leaf) {
			q.branchSize++
			return true
		}
	}

	// This return should NEVER be hit,
	return false
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

func (q *QuadTree[T]) collectChildren() []*TreeLeaf[T] {
	data := make([]*TreeLeaf[T], 0, q.branchSize)
	data = append(data, q.data...)
	if q.nodes[0] != nil {
		for _, node := range q.nodes {
			data = append(data, node.collectChildren()...)
		}
	}
	return data
}

// FindNearest returns up to count elements, sorted by distance to (x,y)
func (q *QuadTree[T]) FindNearest(x, y float64, count int) []*TreeLeaf[T] {
	leafNode := q.leafContaining(x, y)
	if leafNode == nil {
		return nil
	}
	// walk up the tree until we find a node with at least count children, or we run out of tree
	for leafNode.branchSize <= count && leafNode.parent != nil {
		leafNode = leafNode.parent
	}
	// Sort 'em
	data := sortDataByDistance(&point{x, y}, leafNode.collectChildren())

	// if we got too many, slice 'em up
	if len(data) > count {
		data = data[:count]
	}

	return data
}

// FindWithin returns all data within the tree contained by the specified bounding box
func (q *QuadTree[T]) FindWithin(x1, y1, x2, y2 float64) []*TreeLeaf[T] {
	boundary := boundary{x1, y1, x2, y2}
	boundary.align()
	return q.findWithin(boundary)
}

func (q *QuadTree[T]) findWithin(bb boundary) []*TreeLeaf[T] {
	if !q.boundary.Intersects(bb) {
		return nil
	}
	result := make([]*TreeLeaf[T], 0, q.branchSize)
	if len(q.data) > 0 {
		for _, data := range q.data {
			if bb.Contains(data.x, data.y) {
				result = append(result, data)
			}
		}
	} else {
		for _, node := range q.nodes {
			if node != nil {
				result = append(result, node.findWithin(bb)...)
			}
		}
	}

	return result
}
