package pretty

import (
	"fmt"
	"math"
	"testing"

	"github.com/KushamiNeko/GoFun/Utility/test"
)

func TestRGB32ToHSV32(t *testing.T) {

	cases := []struct {
		input  []float64
		output []float64
	}{
		{
			input: []float64{
				0.0 / math.MaxUint8,
				255.0 / math.MaxUint8,
				255.0 / math.MaxUint8,
			},
			output: []float64{
				0.5,
				1.0,
				1.0,
			},
		},
		{
			input: []float64{
				255.0 / math.MaxUint8,
				255.0 / math.MaxUint8,
				255.0 / math.MaxUint8,
			},
			output: []float64{
				0.0,
				0.0,
				1.0,
			},
		},
		{
			input: []float64{
				0.0 / math.MaxUint8,
				0.0 / math.MaxUint8,
				0.0 / math.MaxUint8,
			},
			output: []float64{
				0.0,
				0.0,
				0.0,
			},
		},
		{
			input: []float64{
				255.0 / math.MaxUint8,
				0.0 / math.MaxUint8,
				0.0 / math.MaxUint8,
			},
			output: []float64{
				0.0,
				1.0,
				1.0,
			},
		},
		{
			input: []float64{
				0.0 / math.MaxUint8,
				255.0 / math.MaxUint8,
				0.0 / math.MaxUint8,
			},
			output: []float64{
				0.333333333,
				1.0,
				1.0,
			},
		},
		{
			input: []float64{
				255.0 / math.MaxUint8,
				255.0 / math.MaxUint8,
				0.0 / math.MaxUint8,
			},
			output: []float64{
				0.1666666667,
				1.0,
				1.0,
			},
		},
		{
			input: []float64{
				15.0 / math.MaxUint8,
				153.0 / math.MaxUint8,
				153.0 / math.MaxUint8,
			},
			output: []float64{
				0.5,
				0.9,
				0.6,
			},
		},
		{
			input: []float64{
				27.0 / math.MaxUint8,
				89.0 / math.MaxUint8,
				89.0 / math.MaxUint8,
			},
			output: []float64{
				0.5,
				0.7,
				0.35,
			},
		},
	}

	const precision = 2

	for i, c := range cases {
		t.Run(fmt.Sprintf("case@%d", i), func(t *testing.T) {
			h, s, v := RGB32ToHSV32(c.input[0], c.input[1], c.input[2])
			if math.Round(h*math.Pow10(precision)) != math.Round(c.output[0]*math.Pow10(precision)) {
				t.Errorf("expect: %.4f, get: %.4f", c.output[0], h)
			}
			if math.Round(s*math.Pow10(precision)) != math.Round(c.output[1]*math.Pow10(precision)) {
				t.Errorf("expect: %.4f, get: %.4f", c.output[1], s)
			}
			if math.Round(v*math.Pow10(precision)) != math.Round(c.output[2]*math.Pow10(precision)) {
				t.Errorf("expect: %.4f, get: %.4f", c.output[2], v)
			}
		})
	}
}

func TestValidateHexColor(t *testing.T) {

	cases := []struct {
		input       string
		shouldPanic bool
	}{
		{
			input:       "000000",
			shouldPanic: false,
		},
		{
			input:       "#000000",
			shouldPanic: false,
		},
		{
			input:       "ffffff",
			shouldPanic: false,
		},
		{
			input:       "#ffffff",
			shouldPanic: false,
		},
		{
			input:       "#ff00ff",
			shouldPanic: false,
		},
		{
			input:       "#ffffff ",
			shouldPanic: true,
		},
		{
			input:       " #ffffff",
			shouldPanic: true,
		},
		{
			input:       "00000",
			shouldPanic: true,
		},
		{
			input:       " 000000",
			shouldPanic: true,
		},
		{
			input:       "#0000000",
			shouldPanic: true,
		},
		{
			input:       "0000000",
			shouldPanic: true,
		},
		{
			input:       "-1-1-1",
			shouldPanic: true,
		},
		{
			input:       "gggggg",
			shouldPanic: true,
		},
		{
			input:       "ffff0",
			shouldPanic: true,
		},
		{
			input:       "#ffff0",
			shouldPanic: true,
		},
		{
			input:       "@000000",
			shouldPanic: true,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case@%d", i), func(t *testing.T) {
			if c.shouldPanic {
				test.ShouldPanic(t, func() {
					HexToRGB8(c.input)
				})
			} else {
				HexToRGB8(c.input)
			}
		})
	}
}

func TestHexToRGB8(t *testing.T) {
	cases := []struct {
		input  string
		output []uint8
	}{
		{
			input: "000000",
			output: []uint8{
				0,
				0,
				0,
			},
		},
		{
			input: "ffffff",
			output: []uint8{
				255,
				255,
				255,
			},
		},
		{
			input: "bf6007",
			output: []uint8{
				191,
				96,
				7,
			},
		},
		{
			input: "#55934b",
			output: []uint8{
				85,
				147,
				75,
			},
		},
		{
			input: "#041421",
			output: []uint8{
				4,
				20,
				33,
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case@%d", i), func(t *testing.T) {
			r, g, b := HexToRGB8(c.input)

			if r != c.output[0] {
				t.Errorf("R should be %d but get %d\n", c.output[0], r)
			}

			if g != c.output[1] {
				t.Errorf("G should be %d but get %d\n", c.output[1], r)
			}

			if b != c.output[2] {
				t.Errorf("B should be %d but get %d\n", c.output[2], r)
			}
		})
	}

}
