package pgm

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"testing"
	"time"
	"fmt"
)

func BenchmarkIndex(b *testing.B) {
	// read file
	// use TestGendata to generate corpus.
	f, err := os.Open("testdata/sorted.txt")
	if err != nil {
		b.Error(err)
	}

	data := make([]float64, 0, 1024)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		res, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			b.Error(err)
		}
		data = append(data, res)
	}

	ind := NewIndex(data, 128)
	for _, level := range ind.levels {
		fmt.Println(level)
	}

	now := time.Now()
	pos, err := ind.Search(1471031908028)
	if err != nil {
		b.Error(err)
	}
	b.Log(int(math.Round(data[pos.Lo])), int(math.Round(data[pos.Hi])), time.Since(now))

	// pos, err = ind.Search(1471031908020)
	// if err != nil {
	// 	b.Error(err)
	// }
	// b.Log(int(math.Round(data[pos.Lo])), int(math.Round(data[pos.Hi])))

	// pos, err = ind.Search(1471031908030)
	// if err != nil {
	// 	b.Error(err)
	// }
	// b.Log(int(math.Round(data[pos.Lo])), int(math.Round(data[pos.Hi])))
}
