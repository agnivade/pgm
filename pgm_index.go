package pgm

import (
	"errors"
	"math"
)

// Index contains the PGM index.
// It is just a hierarchy of segments.
type Index struct {
	levels  [][]Segment
	epsilon int
	length  int // length of the array.
}

// Segment contains the slope and intercept of the line segment
// and the corresponding key.
type Segment struct {
	key       float64
	slope     float64
	intercept float64
}

type ApproxPos struct {
	Pos int
	Lo  int
	Hi  int
}

type point struct {
	x float64
	y float64
}

// NewIndex returns a new static PGM index.
func NewIndex(input []float64, epsilon int) *Index {
	index := &Index{
		levels:  make([][]Segment, 0),
		epsilon: epsilon,
		length:  len(input),
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

func (ind *Index) Search(k float64) (ApproxPos, error) {
	if len(ind.levels) < 1 || len(ind.levels[len(ind.levels)-1]) != 1 {
		return ApproxPos{}, errors.New("invalid index")
	}
	root := ind.levels[len(ind.levels)-1][0]
	pos := computePos(k, root)

	for i := len(ind.levels) - 2; i >= 0; i-- {
		lo := max(pos-ind.epsilon, 0)
		hi := min(pos+ind.epsilon, len(ind.levels[i])-1)

		var s, t Segment
		// The rightmost segment s' in levels[i][lo, hi] such that s'.key ≤ k
		for j, item := range ind.levels[i][lo : hi+1] {
			if item.key <= k {
				continue
			}
			s = ind.levels[i][max(lo+j-1, 0)]
			t = item
		}
		if s.key == 0 {
			s = ind.levels[i][hi]
			t = s
		}
		pos = min(computePos(k, s), computePos(t.key, t))
	}
	// We pad the results by one to account
	// for the inaccuracies created by rounding off.
	// This needs to be fixed.
	lo := max(pos-ind.epsilon-1, 0)
	hi := min(pos+ind.epsilon+1, ind.length-1)
	return ApproxPos{Lo: lo, Hi: hi, Pos: pos}, nil
}

func computePos(k float64, s Segment) int {
	// TODO: The rounding off introduces inaccuracies.
	// We need to use integers all the way.
	return int(math.Round((k - s.intercept) / s.slope))
}

func buildPLAModel(keys []float64, epsilon int) []Segment {
	model := []Segment{}
	temp := []point{}
	for i := 0; i < len(keys); i++ {
		temp = append(temp, point{x: float64(i), y: keys[i]})
		if len(temp) < 3 {
			continue
		}
		hull := buildHull(temp)
		r := hull.getSmallestRectangle()
		h := r.getHeight()

		if h > float64(2*epsilon) {
			slope, intercept := r.slopeAndIntercept()
			// Add to model
			model = append(model, Segment{key: temp[0].y, slope: slope, intercept: intercept})
			// Empty convex hull
			temp = temp[:0]
			i--
		}
	}
	if len(temp) == 1 {
		model = append(model, Segment{key: temp[0].y, slope: 1, intercept: temp[0].y - temp[0].x})
	} else if len(temp) > 1 {
		r := buildHull(temp).getSmallestRectangle()
		slope, intercept := r.slopeAndIntercept()
		model = append(model, Segment{key: temp[0].y, slope: slope, intercept: intercept})
	}
	return model
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
