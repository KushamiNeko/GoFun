package indicator

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

func rolling(data []float64, n int, cb func(val float64, roll []float64) float64) []float64 {

	f := make([]float64, len(data))

	for i, v := range data {
		if i < n-1 {
			f[i] = math.NaN()
			continue
		}

		roll := make([]float64, n)
		copy(roll, data[i-(n-1):i+1])
		//for k, v := range data[i-(n-1) : i+1] {
		//roll[k] = v
		//}

		f[i] = cb(v, roll)
	}

	return f
}

func SimpleMovingAverge(data []float64, n int) []float64 {

	ma := rolling(data, n, func(val float64, roll []float64) float64 {

		m := stat.Mean(roll, nil)
		if m <= 0 {
			return math.NaN()
		} else {
			return m
		}
	})

	return ma
}

func BollingerBand(data []float64, n int, mul float64) []float64 {

	bb := rolling(data, n, func(val float64, roll []float64) float64 {
		ma, std := stat.MeanStdDev(roll, nil)

		b := ma + (std * mul)
		if b <= 0 {
			return math.NaN()
		} else {
			return b
		}
	})

	return bb
}

func CommodityChannelIndex(highs, lows, closes []float64, n int) []float64 {
	tp := TypicalPrices(highs, lows, closes)

	cci := rolling(tp, n, func(val float64, roll []float64) float64 {
		ma, std := stat.MeanStdDev(roll, nil)

		c := (val - ma) / (0.015 * std)
		if c <= 0 {
			return math.NaN()
		} else {
			return c
		}
	})

	return cci
}

func TypicalPrices(highs, lows, closes []float64) []float64 {
	if len(highs) != len(lows) || len(lows) != len(closes) {
		panic("high low and close should have the same size")
	}

	tp := make([]float64, len(closes))

	for i := 0; i < len(closes); i++ {
		tp[i] = (highs[i] + lows[i] + closes[i]) / 3.0
	}

	return tp
}
