package pgm

// Index contains the PGM index.
// It is just a hierarchy of segments.
type Index struct {
	levels [][]Segment
}

// Segment contains the slope and intercept of the line segment
// and the corresponding key.
type Segment struct {
	key       float64
	slope     float64
	intercept float64
}

type point struct {
	x float64
	y float64
}

// NewIndex returns a new static PGM index.
func NewIndex(input []float64, epsilon float64) *Index {
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
		hull := buildHull(temp)
		// fmt.Println("hull", hull)
		r := hull.getSmallestRectangle()
		h := r.getHeight()
		// fmt.Println("height", h)

		if h > 2*epsilon {
			slope, intercept := r.slopeAndIntercept()
			// Add to model
			model = append(model, Segment{key: temp[0].y, slope: slope, intercept: intercept})
			// Empty convex hull
			temp = temp[:0]
			i--
		}
	}
	if len(temp) > 0 {
		r := buildHull(temp).getSmallestRectangle()
		slope, intercept := r.slopeAndIntercept()
		model = append(model, Segment{key: temp[0].y, slope: slope, intercept: intercept})
	}
	return model
}

