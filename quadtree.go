// Package quadtree implements a quadtree for storing arbitrary data.
package quadtree

import (
	"github.com/kintar/quadtree/model"
)

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
	boundary   model.BoundingSquare
	depth      int
	data       []*TreeLeaf[T]
	parent     *QuadTree[T]
	nodes      [4]*QuadTree[T]
}

func NewQuadTree[T any](cx, cy, size float64) *QuadTree[T] {
	bounds := model.NewBoundSquare(cx, cy, size)

	return &QuadTree[T]{
		capacity: 8,
		boundary: bounds,
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
	subBounds := model.Subdivide(q.boundary)
	// Iterate over the new bounds and create trees
	for idx, bound := range subBounds {
		x, y := bound.Location2D()
		node := NewQuadTree[T](x, y, bound.Size())
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
	return q.internalInsert(&TreeLeaf[T]{
		data,
		&point{x, y},
	})
}

func (q *QuadTree[T]) internalInsert(leaf *TreeLeaf[T]) bool {
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
		if node.internalInsert(leaf) {
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

	// This lets us determine how many splits have been made to the quadtree, which in turn will tell us the
	// minimum radius around this search point that will intersect at least four nodes. Those nodes MIGHT be within
	// the same division quadrant, or they might not, so in order to actually search them, we have to run a
	// FindWithinCircle centered on the search point with a radius that encompasses the center of the node one depth
	// above our search target.
	searchNode := leafNode
	if leafNode.parent != nil {
		searchNode = leafNode.parent
	}

	// TODO: I'm sure this is not efficient, but it's not a problem yet, so no optimization until it is!
	var results []*TreeLeaf[T]
	// Step up one parent level at a time until we either find the correct number of items or run out of tree
	for len(results) < count {
		// If we're searching a node with no parent, we just want to search the entire node, not some specific radius
		radius := searchNode.boundary.Size()
		// But if we're searching a node with parents, we want to limit the search radius
		if searchNode != leafNode {
			rx, ry := searchNode.boundary.Location2D()
			radius = distance(x, y, rx, ry)
		}

		results = q.FindWithinCircle(model.NewBoundingCircle(x, y, radius))
		if leafNode.parent == nil {
			break
		}
		leafNode = leafNode.parent
	}

	// If we have too many, sort by distance from search target and snip off the excess
	if len(results) > count {
		results = sortDataByDistance(x, y, results)[:count]
	}

	return results
}

func (q *QuadTree[T]) FindWithinSquare(bb model.BoundingSquare) []*TreeLeaf[T] {
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
				result = append(result, node.FindWithinSquare(bb)...)
			}
		}
	}

	return result
}

func (q *QuadTree[T]) FindWithinCircle(c model.BoundingCircle) []*TreeLeaf[T] {
	if !q.boundary.IntersectsCircle(c) {
		return nil
	}
	result := make([]*TreeLeaf[T], 0, q.branchSize)
	if len(q.data) > 0 {
		for _, data := range q.data {
			if c.Contains(data.x, data.y) {
				result = append(result, data)
			}
		}
	} else {
		for _, node := range q.nodes {
			if node != nil {
				result = append(result, node.FindWithinCircle(c)...)
			}
		}
	}

	x, y := c.Location2D()
	return sortDataByDistance[T](x, y, result)
}

func (q *QuadTree[T]) AllChildren() []T {
	var response []T
	for _, node := range q.collectChildren() {
		response = append(response, node.Content)
	}
	return response
}

type TreeVisitorFunc[T any] func(*TreeLeaf[T]) error

// VisitWithinSquare applies a visitor function to every node entry within the given bounding square.
// If the function returns an error, this function immediately returns the error
func (q *QuadTree[T]) VisitWithinSquare(b model.BoundingSquare, visitorFunc TreeVisitorFunc[T]) error {
	var err error
	for _, child := range q.FindWithinSquare(b) {
		if err = visitorFunc(child); err != nil {
			return err
		}
	}
	return nil
}

// VisitWithinCircle applies a visitor function to every node entry within the given bounding circle.
// If the function returns an error, this function immediately returns the error
func (q *QuadTree[T]) VisitWithinCircle(c model.BoundingCircle, visitorFunc TreeVisitorFunc[T]) error {
	var err error
	for _, child := range q.FindWithinCircle(c) {
		if err = visitorFunc(child); err != nil {
			return err
		}
	}
	return nil
}

func (q *QuadTree[T]) Center() (float64, float64) {
	return q.boundary.Location2D()
}

func (q *QuadTree[T]) Size() float64 {
	return q.boundary.Size()
}
