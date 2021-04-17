package main

import (
	"fmt"
	"sort"
)

type Segment struct {
	key       float64
	slope     float64
	intercept float64
}

type point struct {
	x float64
	y float64
}

func buildPGMIndex(input []float64, epsilon float64) {
	levels := make([][]Segment, 0)
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
		levels = append(levels, model)

		// Re-create keys for the next level.
		keys = keys[:0] // wiping
		for _, seg := range model {
			keys = append(keys, seg.key)
		}

		fmt.Println("len model", len(model))

		if len(model) == 1 {
			break
		}
	}

	for _, level := range levels {
		fmt.Println(level)
	}
}

func buildPLAModel(keys []float64, epsilon float64) []Segment {
	model := []Segment{}
	temp := []point{}
	for i, k := range keys {
		temp = append(temp, point{x: float64(i), y: k})
		if len(temp) < 2 {
			continue
		}
		hull := convexHull(temp)
		fmt.Println(hull)
		// q.AddPoint(s2.PointFromCoords(float64(i), k, 0))
		// fmt.Println(q.ConvexHull().Vertices())
		// rect := q.ConvexHull().RectBound()
		// lo := s2.PointFromLatLng(rect.Lo())
		// hi := s2.PointFromLatLng(rect.Hi())

		// // Calculate the height of the rectangle
		// height := 9.0
		// // height := r.Size().Y
		// // fmt.Println("inner", k, height)
		// if height > 2*epsilon {
		// 	slope := (hi.Y - lo.Y) / (hi.X - lo.X)
		// 	intercept := hi.Y - slope*hi.X
		// 	// Add to model
		// 	model = append(model, Segment{key: k, slope: slope, intercept: intercept})
		// 	// Empty convex hull
		// 	q = s2.NewConvexHullQuery()
		// }
	}
	return model
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

func main() {
	input := []float64{2, 12, 15, 18, 23, 24, 29, 31, 34, 36, 38, 48, 55, 59, 60, 71, 73, 74, 76, 88, 95, 102, 115, 122, 123, 124, 158, 159, 161, 164, 165, 187, 189, 190}
	sort.Float64s(input)

	buildPGMIndex(input, 1.0)
}
