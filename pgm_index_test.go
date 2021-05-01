package pgm

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestPGMIndex(t *testing.T) {
	testCases := []struct {
		data      []float64
		segments  int // segments in the last level
		epsilon   int
		searchMem []struct {
			k       float64
			present bool
		}
		searchPred []float64
	}{
		{
			data:     []float64{2, 12, 15, 18, 23, 24, 29, 31, 34, 36, 38, 48, 55, 59, 60, 71, 73, 74, 76, 88, 95, 102, 115, 122, 123, 124, 158, 159, 161, 164, 165, 187, 189, 190},
			segments: 5,
			epsilon:  1,
			searchMem: []struct {
				k       float64
				present bool
			}{
				{k: 31, present: true},
				{k: 32, present: false},
				{k: 80, present: false},
				{k: 95, present: true},
				{k: 190, present: true},
				{k: 200, present: false},
			},
			searchPred: []float64{1, 16, 15, 40, 48, 100, 200},
		},
		{
			data:     []float64{21, 24, 46, 50, 52, 108, 109, 141, 147, 152, 178, 185, 275, 282, 310, 324, 332, 373, 380, 415, 433, 442, 452, 471, 476, 496},
			segments: 3,
			epsilon:  1,
		},
		{
			data:     []float64{1, 2, 13, 36, 37, 57, 69, 107, 140, 176, 215, 229, 246, 260, 288, 324, 337, 341, 381, 390, 409, 411, 416, 442, 444, 453, 476, 497},
			segments: 3,
			epsilon:  1,
		},
		{
			data:     []float64{11, 28, 119, 131, 167, 345, 348, 362, 369, 439},
			segments: 2,
			epsilon:  1,
		},
	}

	for _, tc := range testCases {
		sort.Float64s(tc.data)

		ind := NewIndex(tc.data, tc.epsilon)
		// for _, level := range ind.levels {
		// 	fmt.Println(level)
		// }

		if len(ind.levels[0]) != tc.segments {
			t.Errorf("incorrect number of segments. Got: %d, Want: %d", len(ind.levels[0]), tc.segments)
		}
		verifyIndex(t, ind, tc.data, tc.epsilon)

		if tc.searchMem != nil {
			for _, mem := range tc.searchMem {
				pos, err := ind.Search(mem.k)
				if err != nil {
					t.Errorf("error received: %v", err)
				}

				found := false
				for _, d := range tc.data[pos.Lo : pos.Hi+1] {
					if d == mem.k {
						found = true
						break
					}
				}

				if found != mem.present {
					t.Errorf("incorrect membership result for %f. Got %t, Want: %t", mem.k, found, mem.present)
				}
			}
		}

		if tc.searchPred != nil {
			for _, pred := range tc.searchPred {
				pos, err := ind.Search(pred)
				if err != nil {
					t.Errorf("error received: %v", err)
				}
				t.Log(pred, tc.data[pos.Lo], tc.data[pos.Hi])

				if pos.Lo >= pos.Hi {
					t.Errorf("lo %d is greater than hi %d", pos.Lo, pos.Hi)
				}
				// pos.Lo <= k <= pos.Hi
				if tc.data[pos.Lo] > pred && pos.Lo != 0 {
					t.Errorf("lo %f is greater than k %f", tc.data[pos.Lo], pred)
				}
				if tc.data[pos.Hi] < pred && pos.Hi != len(tc.data)-1 {
					t.Errorf("hi %f is lesser than k %f", tc.data[pos.Hi], pred)
				}
			}
		}
	}
}

func verifyIndex(t *testing.T, ind *Index, input []float64, epsilon int) {
	verifySegment := func(t *testing.T, set []float64, start, end int, s Segment) {
		t.Helper()
		// Iterate all points in the segment and verify they are within e.
		for i := start; i < end; i++ {
			// (y -c)/m - key = err
			err := ((set[i] - s.intercept) / s.slope) - float64(i)
			if err > float64(2*epsilon) {
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
