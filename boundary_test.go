package quadtree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoundary_SubAreas(t *testing.T) {
	b := boundary{
		100, 0,
		200, 100,
	}
	subAreas := b.SubAreas()
	expectedSubs := [4]boundary{
		{100, 0, 150, 50},
		{150, 0, 200, 50},
		{100, 50, 150, 100},
		{150, 50, 200, 100},
	}

	assert.Equal(t, expectedSubs, subAreas)
}

func TestBoundary_Intersects(t *testing.T) {
	b1 := boundary{0, 0, 100, 100}
	assert.True(t, b1.Intersects(boundary{10, 10, 30, 30}))

}
