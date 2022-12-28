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
	qt := NewQuadTree[string](100, 100, 200)
	for y := 0.0; y < 200; y += 10 {
		for x := 0.0; x < 200; x += 10 {
			qt.Insert(x, y, fmt.Sprintf("%.0f, %.0f", x, y))
		}
	}
	resultLeafs := qt.FindNearest(100, 100, 9)
	results := make([]string, 9)
	for i, rl := range resultLeafs {
		results[i] = rl.Content
	}
	expected := []string{
		"100, 100",
		"110, 100",
		"100, 110",
		"110, 110",
		"120, 100",
		"120, 110",
		"100, 120",
		"110, 120",
		"120, 120",
	}
	assert.EqualValues(t, expected, results)
}

func TestQuadTree_FindWithin(t *testing.T) {
	qt := NewQuadTree[string](100, 100, 200)
	for y := 0.0; y < 200; y += 10 {
		for x := 0.0; x < 200; x += 10 {
			qt.Insert(x, y, fmt.Sprintf("%.0f, %.0f", x, y))
		}
	}
	resultLeafs := qt.FindWithin(model.NewBoundSquare(20, 20, 20))
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
