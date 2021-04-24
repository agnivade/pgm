package pgm

import (
	"math"
)

type hull []point

type corner int

const (
	upperRight corner = iota + 1
	upperLeft
	lowerLeft
	lowerRight
)

func cross(a, b, o point) float64 {
	return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
}

// Uses the monotone chain algorithm.
// https://en.wikibooks.org/wiki/Algorithm_Implementation/Geometry/Convex_hull/Monotone_chain
func buildHull(pts []point) hull {
	// points will already be sorted.
	var lower, upper []point

	// Build lower hull
	for i := 0; i < len(pts); i++ {
		for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], pts[i]) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, pts[i])
	}

	// Build upper hull
	for i := len(pts) - 1; i >= 0; i-- {
		for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], pts[i]) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, pts[i])
	}

	lower = lower[:len(lower)-1]
	upper = upper[:len(upper)-1]
	return append(lower, upper...)
}

func (h hull) getSmallestRectangle() rectangle {
	rectangles := []rectangle{}
	i := newCaliper(h, h.getIndex(upperRight), 90)
	j := newCaliper(h, h.getIndex(upperLeft), 180)
	k := newCaliper(h, h.getIndex(lowerLeft), 270)
	l := newCaliper(h, h.getIndex(lowerRight), 0)

	for l.currentAngle < 90 {
		rectangles = append(rectangles, newRectangle(
			l.getIntersection(i),
			i.getIntersection(j),
			j.getIntersection(k),
			k.getIntersection(l),
		))

		smallestTheta := getSmallestTheta(i, j, k, l)

		i.rotateBy(smallestTheta)
		j.rotateBy(smallestTheta)
		k.rotateBy(smallestTheta)
		l.rotateBy(smallestTheta)
	}

	index := 0
	area := math.MaxFloat64

	for i, r := range rectangles {
		tmp := r.getArea()
		if tmp < area {
			area = tmp
			index = i
		}
	}

	return rectangles[index]
}

func (h hull) getIndex(c corner) int {
	index := 0
	pt := h[index]

	for i := 1; i < len(h); i++ {
		tmp := h[i]
		change := false

		switch c {
		case upperRight:
			change = (tmp.x > pt.x || (tmp.x == pt.x && tmp.y > pt.y))
		case upperLeft:
			change = (tmp.y > pt.y || (tmp.y == pt.y && tmp.x < pt.x))
		case lowerLeft:
			change = (tmp.x < pt.x || (tmp.x == pt.x && tmp.y < pt.y))
		case lowerRight:
			change = (tmp.y < pt.y || (tmp.y == pt.y && tmp.x > pt.x))
		}

		if change {
			index = i
			pt = tmp
		}
	}

	return index
}
