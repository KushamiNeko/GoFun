package indicator

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/stat"
)

func TestRollingMean(t *testing.T) {

	tables := []map[string][]float64{
		map[string][]float64{
			"data":   []float64{1, 2, 3, 4, 5},
			"n":      []float64{2},
			"mul":    []float64{1},
			"result": []float64{math.NaN(), 1.5, 2.5, 3.5, 4.5},
		},
		map[string][]float64{
			"data":   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			"n":      []float64{4},
			"mul":    []float64{1},
			"result": []float64{math.NaN(), math.NaN(), math.NaN(), 2.5, 3.5, 4.5, 5.5, 6.5, 7.5},
		},
		map[string][]float64{
			"data":   []float64{1, 2, 3, 4, 5},
			"n":      []float64{2},
			"mul":    []float64{4},
			"result": []float64{math.NaN(), 6, 10, 14, 18},
		},
	}

	for _, table := range tables {
		nd := rolling(table["data"], int(table["n"][0]), func(val float64, roll []float64) float64 {
			return stat.Mean(roll, nil) * table["mul"][0]
		})

		for i, v := range nd {
			ex := table["result"]

			if math.IsNaN(v) {
				if !math.IsNaN(ex[i]) {
					t.Errorf("expect %f but get %f", ex[i], v)
				}
			} else {
				if v != ex[i] {
					t.Errorf("expect %f but get %f", ex[i], v)
				}
			}

		}
	}

}
