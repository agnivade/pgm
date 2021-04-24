package pgm

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestPGMIndex(t *testing.T) {
	input := []struct {
		data     []float64
		segments int // segments in the last level
		epsilon  float64
	}{
		{
			data:     []float64{2, 12, 15, 18, 23, 24, 29, 31, 34, 36, 38, 48, 55, 59, 60, 71, 73, 74, 76, 88, 95, 102, 115, 122, 123, 124, 158, 159, 161, 164, 165, 187, 189, 190},
			segments: 5,
			epsilon:  1.0,
		},
		{
			data:     []float64{21, 24, 46, 50, 52, 108, 109, 141, 147, 152, 178, 185, 275, 282, 310, 324, 332, 373, 380, 415, 433, 442, 452, 471, 476, 496},
			segments: 3,
			epsilon:  1.0,
		},
		{
			data:     []float64{1, 2, 13, 36, 37, 57, 69, 107, 140, 176, 215, 229, 246, 260, 288, 324, 337, 341, 381, 390, 409, 411, 416, 442, 444, 453, 476, 497},
			segments: 3,
			epsilon:  1.0,
		},
		{
			data:     []float64{11, 28, 119, 131, 167, 345, 348, 362, 369, 439},
			segments: 2,
			epsilon:  1.0,
		},
	}

	for _, item := range input {
		sort.Float64s(item.data)

		ind := NewIndex(item.data, item.epsilon)
		for _, level := range ind.levels {
			fmt.Println(level)
		}

		if len(ind.levels[0]) != item.segments {
			t.Errorf("incorrect number of segments. Got: %d, Want: %d", len(ind.levels[0]), item.segments)
		}
		verifyIndex(t, ind, item.data, item.epsilon)
	}
}

func verifyIndex(t *testing.T, ind *Index, input []float64, epsilon float64) {
	var verifySegment = func(t *testing.T, set []float64, start, end int, s Segment) {
		t.Helper()
		// Iterate all points in the segment and verify they are within e.
		for i := start; i < end; i++ {
			// (y -c)/m - key = err
			err := ((set[i] - s.intercept) / s.slope) - float64(i)
			if err > 2*epsilon {
				t.Errorf("error threshold exceeded, x: %d, y:%f", i, set[i])
			}
		}
	}

	// Verify each level
	for i := 0; i < len(ind.levels); i++ {
		level := ind.levels[i]
		var set []float64
		if i == 0 {
			set = input
		} else {
			set = set[:0] // reset
			for _, seg := range ind.levels[i-1] {
				set = append(set, seg.key)
			}
		}
		// Find set of keys for each segment
		current := 0
		for j := 0; j < len(level); j++ {
			s := level[j]
			// Check if last segment or not
			if j+1 == len(level) {
				verifySegment(t, set, current, len(set), s)
			} else {
				nextKey := level[j+1].key
				start := current
				for set[current] != nextKey {
					current++
				}
				verifySegment(t, set, start, current, s)
			}
		}
	}
}

func TestGendata(t *testing.T) {
	t.Skip("only to generate corpus")
	var input []int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		input = append(input, rand.Intn(50))
	}
	sort.Ints(input)
	input = removeDups(input)
	fmt.Println(input)
}

func removeDups(in []int) []int {
	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] == in[i] {
			continue
		}
		j++
		in[j] = in[i]
	}
	return in[:j+1]
}
