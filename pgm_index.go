package pgm

import (
	"math"
)

type Index struct {
	levels [][]Segment
}

type Segment struct {
	key       float64
	slope     float64
	intercept float64
}

type point struct {
	x float64
	y float64
}

type rectangle [4]point

type corner int

const (
	upperRight corner = iota + 1
	upperLeft
	lowerLeft
	lowerRight
)

func buildPGMIndex(input []float64, epsilon float64) *Index {
	index := &Index{
		levels: make([][]Segment, 0),
	}
	keys := make([]float64, len(input))
	copy(keys, input)
	// repeat
	// M = Build-PLA-model(keys, ε)
	// levels[i] = M; i = i + 1
	// m = Size(M)
	// keys = [M[0].key, . . . , M[m − 1].key]
	// until m = 1
	for {
		model := buildPLAModel(keys, epsilon) // returns a single level
		index.levels = append(index.levels, model)

		// Re-create keys for the next level.
		keys = keys[:0] // wiping
		for _, seg := range model {
			keys = append(keys, seg.key)
		}

		if len(model) == 1 {
			break
		}
	}

	return index
}

func buildPLAModel(keys []float64, epsilon float64) []Segment {
	model := []Segment{}
	temp := []point{}
	for i := 0; i < len(keys); i++ {
		temp = append(temp, point{x: float64(i), y: keys[i]})
		if len(temp) < 3 {
			continue
		}
		hull := convexHull(temp)
		// fmt.Println("hull", hull)
		r := getSmallestRectangle(hull)
		h := getHeight(r)
		// fmt.Println("height", h)

		if h > 2*epsilon {
			slope, intercept := getSlopeAndIntercept(r)
			// Add to model
			model = append(model, Segment{key: temp[0].y, slope: slope, intercept: intercept})
			// Empty convex hull
			temp = temp[:0]
			i--
		}
	}
	if len(temp) > 0 {
		r := getSmallestRectangle(convexHull(temp))
		slope, intercept := getSlopeAndIntercept(r)
		model = append(model, Segment{key: temp[0].y, slope: slope, intercept: intercept})
	}
	return model
}

func getSmallestRectangle(hull []point) rectangle {
	rectangles := []rectangle{}
	i := newCaliper(hull, getIndex(hull, upperRight), 90)
	j := newCaliper(hull, getIndex(hull, upperLeft), 180)
	k := newCaliper(hull, getIndex(hull, lowerLeft), 270)
	l := newCaliper(hull, getIndex(hull, lowerRight), 0)

	for l.currentAngle < 90 {
		rectangles = append(rectangles, rectangle{
			l.getIntersection(i),
			i.getIntersection(j),
			j.getIntersection(k),
			k.getIntersection(l),
		})

		smallestTheta := getSmallestTheta(i, j, k, l)

		i.rotateBy(smallestTheta)
		j.rotateBy(smallestTheta)
		k.rotateBy(smallestTheta)
		l.rotateBy(smallestTheta)
	}

	index := 0
	area := math.MaxFloat64

	for i, r := range rectangles {
		tmp := getArea(r)
		if tmp < area {
			area = tmp
			index = i
		}
	}

	return rectangles[index]
}

func getArea(r rectangle) float64 {
	deltaXAB := r[0].x - r[1].x
	deltaYAB := r[0].y - r[1].y
	deltaXBC := r[1].x - r[2].x
	deltaYBC := r[1].y - r[2].y

	lengthAB := math.Sqrt((deltaXAB * deltaXAB) + (deltaYAB * deltaYAB))
	lengthBC := math.Sqrt((deltaXBC * deltaXBC) + (deltaYBC * deltaYBC))

	return lengthAB * lengthBC
}

func getHeight(r rectangle) float64 {
	deltaXAB := r[0].x - r[1].x
	deltaYAB := r[0].y - r[1].y

	lengthAB := math.Sqrt((deltaXAB * deltaXAB) + (deltaYAB * deltaYAB))
	return lengthAB
}

func getSlopeAndIntercept(r rectangle) (slope, intercept float64) {
	// We take the midpoint of two sides and draw a line between them.
	p1 := point{x: (r[0].x + r[1].x) / 2, y: (r[0].y + r[1].y) / 2}
	p2 := point{x: (r[2].x + r[3].x) / 2, y: (r[2].y + r[3].y) / 2}

	slope = (p2.y - p1.y) / (p2.x - p1.x)
	intercept = p2.y - slope*p2.x
	return
}

func getSmallestTheta(i, j, k, l *Caliper) float64 {
	thetaI := i.getDeltaAngleNextPoint()
	thetaJ := j.getDeltaAngleNextPoint()
	thetaK := k.getDeltaAngleNextPoint()
	thetaL := l.getDeltaAngleNextPoint()

	if thetaI <= thetaJ && thetaI <= thetaK && thetaI <= thetaL {
		return thetaI
	} else if thetaJ <= thetaK && thetaJ <= thetaL {
		return thetaJ
	} else if thetaK <= thetaL {
		return thetaK
	} else {
		return thetaL
	}
}

func getIndex(hull []point, c corner) int {
	index := 0
	pt := hull[index]

	for i := 1; i < len(hull); i++ {
		tmp := hull[i]
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

func cross(a, b, o point) float64 {
	return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
}

// points will already be sorted.
// Uses the monotone chain algorithm.
// https://en.wikibooks.org/wiki/Algorithm_Implementation/Geometry/Convex_hull/Monotone_chain
func convexHull(pts []point) []point {
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
