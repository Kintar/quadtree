package quadtree

import (
	"fmt"
	"github.com/kintar/quadtree/model"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestQuadTree_Insert_BasicFunctionality(t *testing.T) {
	qt := NewQuadTree[int](50, 50, 100)
	qt.Insert(50, 50, 1)
	assert.EqualValues(t, 1, len(qt.data))
}

func TestQuadTree_Insert_DividesAtCapacity(t *testing.T) {
	qt := NewQuadTree[int](50, 50, 100)
	for i := 0; i < qt.capacity+2; i++ {
		qt.Insert(50, float64(i), i)
	}

	assert.False(t, qt.nodes[0] == nil)
	assert.EqualValues(t, 0, len(qt.data), "parent node data is not being cleared")
}

func TestQuadTree_Len(t *testing.T) {
	qt := NewQuadTree[int](50, 50, 100)
	const count = 10
	for i := 0; i < count; i++ {
		qt.Insert(51, float64(i), i)
	}
	assert.EqualValues(t, count, qt.Len())
}

func TestQuadTree_collectChildren(t *testing.T) {
	qt := NewQuadTree[int](50, 50, 100)
	const count = 10
	for i := 0; i < count; i++ {
		qt.Insert(51, float64(i), i)
	}
	children := qt.collectChildren()
	assert.EqualValues(t, count, len(children))
}

func TestQuadTree_Insert_HighCapacity(t *testing.T) {
	defer func() {
		p := recover()
		if p != nil {
			t.Fatalf("panic caught: %v", p)
		}
	}()

	qt := NewQuadTree[string](50, 50, 100)
	for y := 0.0; y < 200; y += 10 {
		for x := 0.0; x < 200; x += 10 {
			qt.Insert(x, y, fmt.Sprintf("%.0f, %.0f", x, y))
		}
	}
}

func TestQuadTree_FindNearest(t *testing.T) {
	qt := NewQuadTree[string](10, 10, 20)
	// The tree is centered on 10, 10, which means after one subdivision the quadrants will be
	// [0,0-10,10), [10,0-20,10), [0,10-10,20), [10,10-20,20)
	//
	// Knowing this, we can cluster four very close points at the center of the tree, but in separate quadrants, then
	// add enough points to the top-left of the tree to cause a second subdivision in the first quadrant. This means
	// that a naive search originating at the point closest to origin within the largest top-left quadrant will only
	// return items which all lay within the largest top-left quadrant, even though there are closer points in the other
	// three top-level quadrants
	qt.Insert(9, 9, "one")
	qt.Insert(11, 9, "two")
	qt.Insert(9, 11, "three")
	qt.Insert(11, 11, "four")

	expected := []string{
		"one", "two", "three", "four",
	}
	sort.Strings(expected)

	// If we fetch now, we'll get our expected results
	resultLeafs := qt.FindNearest(9, 9, 4)
	var results []string
	for _, l := range resultLeafs {
		results = append(results, l.Content)
	}
	sort.Strings(results)
	assert.Equal(t, expected, results)

	// Now if we add eight more points to the upper left quadrant, it will subdivide into:
	// [0.0-5,5), [5,0-10,5), [0,5-5,10), [5,5-10,10)
	qt.Insert(5, 5, "five")
	qt.Insert(6, 6, "six")
	qt.Insert(4, 4, "seven")
	qt.Insert(3, 3, "eight")
	qt.Insert(3, 5, "nine")
	qt.Insert(2, 2, "ten")
	qt.Insert(6, 4, "eleven")
	qt.Insert(7, 3, "twelve")
	qt.Insert(8, 2, "thirteen")

	// spatially, a query for the four points closest to 9,9 should still return the points clustered at 10,10
	resultLeafs = qt.FindNearest(9, 9, 4)
	results = nil
	for _, l := range resultLeafs {
		results = append(results, l.Content)
	}
	sort.Strings(results)
	assert.Equal(t, expected, results)

}

func TestQuadTree_FindWithinSquare(t *testing.T) {
	qt := NewQuadTree[string](100, 100, 200)
	for y := 0.0; y < 200; y += 10 {
		for x := 0.0; x < 200; x += 10 {
			qt.Insert(x, y, fmt.Sprintf("%.0f, %.0f", x, y))
		}
	}
	resultLeafs := qt.FindWithinSquare(model.NewBoundSquare(20, 20, 20))
	results := make([]string, len(resultLeafs))
	for i, rl := range resultLeafs {
		results[i] = rl.Content
	}
	expected := []string{
		"10, 10",
		"20, 10",
		"10, 20",
		"20, 20",
	}

	sort.Strings(expected)
	sort.Strings(results)
	assert.EqualValues(t, expected, results)

}

func TestQuadTree_FindWithinCircle(t *testing.T) {
	qt := NewQuadTree[string](100, 100, 200)
	qt.Insert(100, 100, "100, 100")
	qt.Insert(10, 20, "10, 20")
	qt.Insert(10, 10, "10, 10")
	qt.Insert(20, 10, "20, 10")

	resultLeafs := qt.FindWithinCircle(model.NewBoundingCircle(20, 20, 20))
	results := make([]string, len(resultLeafs))
	for i, rl := range resultLeafs {
		results[i] = rl.Content
	}
	expected := []string{
		"10, 20",
		"20, 10",
		"10, 10",
	}

	sort.Strings(expected)
	sort.Strings(results)
	assert.EqualValues(t, expected, results)

}
