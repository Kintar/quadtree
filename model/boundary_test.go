package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBoundSquare(t *testing.T) {
	b := NewBoundSquare(50, 50, 100)
	expected := BoundingSquare{
		cx:     50,
		cy:     50,
		extent: 50,
		size:   100,
	}
	assert.Equal(t, b, expected)
}

func TestNewBoundSquareFromCorner(t *testing.T) {
	b := NewBoundSquareFromCorner(0, 0, 100)
	expected := BoundingSquare{
		cx:     50,
		cy:     50,
		extent: 50,
		size:   100,
	}
	assert.Equal(t, b, expected)
}

func TestBoundingSquare_Contains(t *testing.T) {
	b := NewBoundSquare(50, 50, 100)
	assert.True(t, b.Contains(25, 10))
	assert.False(t, b.Contains(-5, 50))
}

func TestBoundingSquare_SubAreas(t *testing.T) {
	b := NewBoundSquare(150, 50, 100)
	assert.EqualValues(t, 50, b.extent)
	assert.EqualValues(t, 100, b.size)
	subAreas := Subdivide(b)
	expectedSubs := [4]BoundingSquare{
		{125, 25, 25, 50},
		{175, 25, 25, 50},
		{125, 75, 25, 50},
		{175, 75, 25, 50},
	}

	assert.Equal(t, expectedSubs, subAreas)
}

func TestBoundingSquare_Intersects(t *testing.T) {
	b1 := NewBoundSquare(50, 50, 100)
	assert.True(t, b1.Intersects(NewBoundSquare(25, 25, 15)))
}

func TestBoundingCircle_Intersects(t *testing.T) {
	c1 := NewBoundingCircle(10, 10, 5)
	c2 := NewBoundingCircle(7, 7, 2)
	assert.True(t, c1.Intersects(c2))

	c2 = NewBoundingCircle(17, 10, 8)
	assert.True(t, c1.Intersects(c2))
}

func TestIntersectCircleAndSquare(t *testing.T) {
	b := NewBoundSquare(150, 150, 100)
	c := NewBoundingCircle(10, 10, 20)
	// Case one non-intersection due to center distances being greater than extent+radius
	assert.False(t, b.IntersectsCircle(c))

	// Case two : if either component is within radius+extent, but the first case did not fail
	c = NewBoundingCircle(90, 101, 20)
	assert.True(t, c.IntersectsSquare(b))

	// Case two : again, but different component
	c = NewBoundingCircle(101, 90, 20)
	assert.True(t, b.IntersectsCircle(c))

	// Case three : none of the above are true, but one corner lies within the circle
	c = NewBoundingCircle(50, 50, 72)
	assert.True(t, c.IntersectsSquare(b))
}

func TestBoundingCircle_Contains(t *testing.T) {
	c := NewBoundingCircle(50, 50, 20)
	assert.True(t, c.Contains(40, 40))
	assert.False(t, c.Contains(80, 50))
}
