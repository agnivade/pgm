package pgm

import (
	"math"
)

type rectangle [4]point

func newRectangle(i, j, k, l point) rectangle {
	return rectangle{i, j, k, l}
}

func (r *rectangle) getArea() float64 {
	deltaXAB := r[0].x - r[1].x
	deltaYAB := r[0].y - r[1].y
	deltaXBC := r[1].x - r[2].x
	deltaYBC := r[1].y - r[2].y

	lengthAB := math.Sqrt((deltaXAB * deltaXAB) + (deltaYAB * deltaYAB))
	lengthBC := math.Sqrt((deltaXBC * deltaXBC) + (deltaYBC * deltaYBC))

	return lengthAB * lengthBC
}

func (r *rectangle) getHeight() float64 {
	deltaXAB := r[0].x - r[1].x
	deltaYAB := r[0].y - r[1].y

	lengthAB := math.Sqrt((deltaXAB * deltaXAB) + (deltaYAB * deltaYAB))
	return lengthAB
}

func (r *rectangle) slopeAndIntercept() (slope, intercept float64) {
	// We take the midpoint of two sides and draw a line between them.
	p1 := point{x: (r[0].x + r[1].x) / 2, y: (r[0].y + r[1].y) / 2}
	p2 := point{x: (r[2].x + r[3].x) / 2, y: (r[2].y + r[3].y) / 2}

	slope = (p2.y - p1.y) / (p2.x - p1.x)
	intercept = p2.y - slope*p2.x
	return
}
