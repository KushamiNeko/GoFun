package data

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/KushamiNeko/GoFun/Fun/Chart/futures"
)

func currentSymbol(cct continuousContract, current *futures.Contract, symbol string) *futures.Contract {
	if current == nil || current.Symbol() != symbol {
		cts, err := cct.readContract(symbol)
		if err != nil {
			panic(err)
		}

		current = futures.NewContract(
			symbol,
			cts,
			futures.BarchartSymbolFormat,
		)
	}

	return current
}

func TestContinuousContractFirstOfMonthPanamaCanal(t *testing.T) {
	t.Parallel()

	from, err := time.Parse(`20060102`, "20190101")
	if err != nil {
		panic(err)
	}

	to, err := time.Parse(`20060102`, "20191231")
	if err != nil {
		panic(err)
	}

	cct := continuousContract{}

	ts, err := cct.read(
		from,
		to,
		"es",
		futures.FinancialContractMonths,
		futures.BarchartSymbolFormat,
		futures.FirstOfMonth,
		futures.PanamaCanal,
	)
	if err != nil {
		panic(err)
	}

	var (
		current    *futures.Contract
		adjustment float64
	)

	cols := []string{
		"open",
		"high",
		"low",
		"close",
		"volume",
		"openinterest",
	}

	loc := time.Now().Location()
	rollingTimes := []time.Time{
		time.Date(
			2019,
			time.March,
			1,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.June,
			3,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.September,
			3,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.December,
			2,
			0,
			0,
			0,
			0,
			loc,
		),
	}

	tsl := len(ts.FullTimes())

	for i := tsl - 1; i >= 0; i-- {
		dt := ts.FullTimes()[i]

		switch {
		case dt.Before(rollingTimes[0]):
			current = currentSymbol(cct, current, "esh19")
			adjustment = 10.5

		case (dt.Equal(rollingTimes[0]) || dt.After(rollingTimes[0])) && dt.Before(rollingTimes[1]):
			current = currentSymbol(cct, current, "esm19")
			adjustment = 5.25

		case (dt.Equal(rollingTimes[1]) || dt.After(rollingTimes[1])) && dt.Before(rollingTimes[2]):
			current = currentSymbol(cct, current, "esu19")
			adjustment = 2.25

		case (dt.Equal(rollingTimes[2]) || dt.After(rollingTimes[2])) && dt.Before(rollingTimes[3]):
			current = currentSymbol(cct, current, "esz19")
			adjustment = 1.75

		case dt.Equal(rollingTimes[3]) || dt.After(rollingTimes[3]):
			current = currentSymbol(cct, current, "esh20")
			adjustment = 0

		default:
			panic("unknown time")
		}

		for _, col := range cols {
			ccv := ts.ValueInFullTimes(dt, col, math.NaN())
			cuv := current.Series().ValueInFullTimes(dt, col, math.NaN())

			if math.IsNaN(ccv) || math.IsNaN(cuv) {
				t.Errorf("value should not be NaN")
			}

			if col != "volume" && col != "openinterest" {
				if ccv != cuv+adjustment {
					t.Errorf("time: %s, symbol: %s, column: %s, expect: %f, get: %f", dt.Format("20060102"), current.Symbol(), col, cuv+adjustment, ccv)
				}
			} else {
				if ccv != cuv {
					t.Errorf("time: %s, symbol: %s, column: %s expect: %f, get: %f", dt.Format("20060102"), current.Symbol(), col, cuv, ccv)
				}
			}
		}

	}

}

func TestContinuousContractFirstOfMonthRatio(t *testing.T) {
	t.Parallel()

	from, err := time.Parse(`20060102`, "20190101")
	if err != nil {
		panic(err)
	}

	to, err := time.Parse(`20060102`, "20191231")
	if err != nil {
		panic(err)
	}

	cct := continuousContract{}

	ts, err := cct.read(
		from,
		to,
		"es",
		futures.FinancialContractMonths,
		futures.BarchartSymbolFormat,
		futures.FirstOfMonth,
		futures.Ratio,
	)
	if err != nil {
		panic(err)
	}

	var (
		current    *futures.Contract
		adjustment float64
	)

	cols := []string{
		"open",
		"high",
		"low",
		"close",
		"volume",
		"openinterest",
	}

	loc := time.Now().Location()
	rollingTimes := []time.Time{
		time.Date(
			2019,
			time.March,
			1,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.June,
			3,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.September,
			3,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.December,
			2,
			0,
			0,
			0,
			0,
			loc,
		),
	}

	tsl := len(ts.FullTimes())

	for i := tsl - 1; i >= 0; i-- {
		dt := ts.FullTimes()[i]

		switch {
		case dt.Before(rollingTimes[0]):
			current = currentSymbol(cct, current, "esh19")
			adjustment = 1.003701071

		case (dt.Equal(rollingTimes[0]) || dt.After(rollingTimes[0])) && dt.Before(rollingTimes[1]):
			current = currentSymbol(cct, current, "esm19")
			adjustment = 1.001825996

		case (dt.Equal(rollingTimes[1]) || dt.After(rollingTimes[1])) && dt.Before(rollingTimes[2]):
			current = currentSymbol(cct, current, "esu19")
			adjustment = 1.000734088

		case (dt.Equal(rollingTimes[2]) || dt.After(rollingTimes[2])) && dt.Before(rollingTimes[3]):
			current = currentSymbol(cct, current, "esz19")
			adjustment = 1.000561933

		case dt.Equal(rollingTimes[3]) || dt.After(rollingTimes[3]):
			current = currentSymbol(cct, current, "esh20")
			adjustment = 1

		default:
			panic("unknown time")
		}

		for _, col := range cols {
			ccv := ts.ValueInFullTimes(dt, col, math.NaN())
			cuv := current.Series().ValueInFullTimes(dt, col, math.NaN())

			if math.IsNaN(ccv) || math.IsNaN(cuv) {
				t.Errorf("value should not be NaN")
			}

			if col != "volume" && col != "openinterest" {
				ndigit := 10.0
				if math.Round(ccv*ndigit) != math.Round(cuv*adjustment*ndigit) {
					t.Errorf(
						"time: %s, symbol: %s, column: %s, expect: %f, get: %f",
						dt.Format("20060102"),
						current.Symbol(),
						col,
						math.Round(cuv*adjustment*ndigit),
						math.Round(ccv*ndigit),
					)
				}
			} else {
				if ccv != cuv {
					t.Errorf("time: %s, symbol: %s, column: %s expect: %f, get: %f", dt.Format("20060102"), current.Symbol(), col, cuv, ccv)
				}
			}
		}

	}

}

func TestContinuousContractOpenInterestPanamaCanal(t *testing.T) {
	t.Parallel()
	t.SkipNow()

	from, err := time.Parse(`20060102`, "20190101")
	if err != nil {
		panic(err)
	}

	to, err := time.Parse(`20060102`, "20191231")
	if err != nil {
		panic(err)
	}

	cct := continuousContract{}

	ts, err := cct.read(
		from,
		to,
		"es",
		futures.FinancialContractMonths,
		futures.BarchartSymbolFormat,
		futures.OpenInterest,
		futures.PanamaCanal,
	)
	if err != nil {
		panic(err)
	}

	var (
		current    *futures.Contract
		adjustment float64
	)

	cols := []string{
		"open",
		"high",
		"low",
		"close",
		"volume",
		"openinterest",
	}

	loc := time.Now().Location()
	rollingTimes := []time.Time{
		time.Date(
			2019,
			time.March,
			11,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.June,
			17,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.September,
			16,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.December,
			16,
			0,
			0,
			0,
			0,
			loc,
		),
	}

	tsl := len(ts.FullTimes())

	for i := tsl - 1; i >= 0; i-- {
		dt := ts.FullTimes()[i]

		switch {
		case dt.Before(rollingTimes[0]):
			current = currentSymbol(cct, current, "esh19")
			adjustment = 16

		case (dt.Equal(rollingTimes[0]) || dt.After(rollingTimes[0])) && dt.Before(rollingTimes[1]):
			current = currentSymbol(cct, current, "esm19")
			adjustment = 11

		case (dt.Equal(rollingTimes[1]) || dt.After(rollingTimes[1])) && dt.Before(rollingTimes[2]):
			current = currentSymbol(cct, current, "esu19")
			adjustment = 6.75

		case (dt.Equal(rollingTimes[2]) || dt.After(rollingTimes[2])) && dt.Before(rollingTimes[3]):
			current = currentSymbol(cct, current, "esz19")
			adjustment = 4.25

		case dt.Equal(rollingTimes[3]) || dt.After(rollingTimes[3]):
			current = currentSymbol(cct, current, "esh20")
			adjustment = 0

		default:
			panic("unknown time")
		}

		for _, col := range cols {
			ccv := ts.ValueInFullTimes(dt, col, math.NaN())
			cuv := current.Series().ValueInFullTimes(dt, col, math.NaN())

			if math.IsNaN(ccv) || math.IsNaN(cuv) {
				t.Errorf("value should not be NaN")
			}

			if col != "volume" && col != "openinterest" {
				if ccv != cuv+adjustment {
					t.Errorf("time: %s, symbol: %s, column: %s, expect: %f, get: %f", dt.Format("20060102"), current.Symbol(), col, cuv+adjustment, ccv)
				}
			} else {
				if ccv != cuv {
					t.Errorf("time: %s, symbol: %s, column: %s expect: %f, get: %f", dt.Format("20060102"), current.Symbol(), col, cuv, ccv)
				}
			}
		}

	}

}

func TestContinuousContractOpenInterestRatio(t *testing.T) {
	t.Parallel()
	t.SkipNow()

	from, err := time.Parse(`20060102`, "20190101")
	if err != nil {
		panic(err)
	}

	to, err := time.Parse(`20060102`, "20191231")
	if err != nil {
		panic(err)
	}

	cct := continuousContract{}

	ts, err := cct.read(
		from,
		to,
		"es",
		futures.FinancialContractMonths,
		futures.BarchartSymbolFormat,
		futures.OpenInterest,
		futures.Ratio,
	)
	if err != nil {
		panic(err)
	}

	var (
		current    *futures.Contract
		adjustment float64
	)

	cols := []string{
		"open",
		"high",
		"low",
		"close",
		"volume",
		"openinterest",
	}

	loc := time.Now().Location()
	rollingTimes := []time.Time{
		time.Date(
			2019,
			time.March,
			11,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.June,
			17,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.September,
			16,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.December,
			16,
			0,
			0,
			0,
			0,
			loc,
		),
	}

	tsl := len(ts.FullTimes())

	for i := tsl - 1; i >= 0; i-- {
		dt := ts.FullTimes()[i]

		switch {
		case dt.Before(rollingTimes[0]):
			current = currentSymbol(cct, current, "esh19")
			adjustment = 1.0054405

		case (dt.Equal(rollingTimes[0]) || dt.After(rollingTimes[0])) && dt.Before(rollingTimes[1]):
			current = currentSymbol(cct, current, "esm19")
			adjustment = 1.003637989

		case (dt.Equal(rollingTimes[1]) || dt.After(rollingTimes[1])) && dt.Before(rollingTimes[2]):
			current = currentSymbol(cct, current, "esu19")
			adjustment = 1.002165236

		case (dt.Equal(rollingTimes[2]) || dt.After(rollingTimes[2])) && dt.Before(rollingTimes[3]):
			current = currentSymbol(cct, current, "esz19")
			adjustment = 1.001330516

		case dt.Equal(rollingTimes[3]) || dt.After(rollingTimes[3]):
			current = currentSymbol(cct, current, "esh20")
			adjustment = 1

		default:
			panic("unknown time")
		}

		for _, col := range cols {
			ccv := ts.ValueInFullTimes(dt, col, math.NaN())
			cuv := current.Series().ValueInFullTimes(dt, col, math.NaN())

			if math.IsNaN(ccv) || math.IsNaN(cuv) {
				t.Errorf("value should not be NaN")
			}

			if col != "volume" && col != "openinterest" {
				ndigit := 1000.0
				if math.Round(ccv*ndigit) != math.Round(cuv*adjustment*ndigit) {
					t.Errorf(
						"time: %s, symbol: %s, column: %s, expect: %f, get: %f",
						dt.Format("20060102"),
						current.Symbol(),
						col,
						math.Round(cuv*adjustment*ndigit),
						math.Round(ccv*ndigit),
					)
				}
			} else {
				if ccv != cuv {
					t.Errorf("time: %s, symbol: %s, column: %s expect: %f, get: %f", dt.Format("20060102"), current.Symbol(), col, cuv, ccv)
				}
			}
		}

	}

}

func TestContinuousContractLastNTradingDayPanamaCanal(t *testing.T) {
	t.Parallel()

	from, err := time.Parse(`20060102`, "20190101")
	if err != nil {
		panic(err)
	}

	to, err := time.Parse(`20060102`, "20191231")
	if err != nil {
		panic(err)
	}

	cct := continuousContract{}

	ts, err := cct.read(
		from,
		to,
		"es",
		futures.FinancialContractMonths,
		futures.BarchartSymbolFormat,
		futures.LastNTradingDay,
		futures.PanamaCanal,
	)
	if err != nil {
		panic(err)
	}

	var (
		current    *futures.Contract
		adjustment float64
	)

	cols := []string{
		"open",
		"high",
		"low",
		"close",
		"volume",
		"openinterest",
	}

	loc := time.Now().Location()
	rollingTimes := []time.Time{
		time.Date(
			2019,
			time.March,
			11,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.June,
			17,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.September,
			16,
			0,
			0,
			0,
			0,
			loc,
		),
		time.Date(
			2019,
			time.December,
			16,
			0,
			0,
			0,
			0,
			loc,
		),
	}

	tsl := len(ts.FullTimes())

	for i := tsl - 1; i >= 0; i-- {
		dt := ts.FullTimes()[i]

		switch {
		case dt.Before(rollingTimes[0]):
			current = currentSymbol(cct, current, "esh19")
			adjustment = 16

		case (dt.Equal(rollingTimes[0]) || dt.After(rollingTimes[0])) && dt.Before(rollingTimes[1]):
			current = currentSymbol(cct, current, "esm19")
			adjustment = 11

		case (dt.Equal(rollingTimes[1]) || dt.After(rollingTimes[1])) && dt.Before(rollingTimes[2]):
			current = currentSymbol(cct, current, "esu19")
			adjustment = 6.75

		case (dt.Equal(rollingTimes[2]) || dt.After(rollingTimes[2])) && dt.Before(rollingTimes[3]):
			current = currentSymbol(cct, current, "esz19")
			adjustment = 4.25

		case dt.Equal(rollingTimes[3]) || dt.After(rollingTimes[3]):
			current = currentSymbol(cct, current, "esh20")
			adjustment = 0

		default:
			panic("unknown time")
		}

		for _, col := range cols {
			ccv := ts.ValueInFullTimes(dt, col, math.NaN())
			cuv := current.Series().ValueInFullTimes(dt, col, math.NaN())

			if math.IsNaN(ccv) || math.IsNaN(cuv) {
				t.Errorf("value should not be NaN")
			}

			if col != "volume" && col != "openinterest" {
				if ccv != cuv+adjustment {
					t.Errorf("time: %s, symbol: %s, column: %s, expect: %f, get: %f", dt.Format("20060102"), current.Symbol(), col, cuv+adjustment, ccv)
				}
			} else {
				if ccv != cuv {
					t.Errorf("time: %s, symbol: %s, column: %s expect: %f, get: %f", dt.Format("20060102"), current.Symbol(), col, cuv, ccv)
				}
			}
		}

	}

}

func TestContinuous(t *testing.T) {
	t.Parallel()

	//from, err := time.Parse(`20060102`, "20090122")
	from, err := time.Parse(`20060102`, "20070920")
	if err != nil {
		panic(err)
	}

	//to, err := time.Parse(`20060102`, "20100122")
	to, err := time.Parse(`20060102`, "20100919")
	if err != nil {
		panic(err)
	}

	cct := continuousContract{}

	_, err = cct.read(
		from,
		to,
		"qr",
		futures.FinancialContractMonths,
		futures.BarchartSymbolFormat,
		futures.LastNTradingDay,
		futures.Ratio,
	)
	if err != nil {
		panic(err)
	}

	//for _, _ = range ts.FullTimes() {
	//fmt.Println(dt)

	//fmt.Println(ts.ValueInFullTimes(dt, "open", math.NaN()))
	//fmt.Println(ts.ValueInFullTimes(dt, "high", math.NaN()))
	//fmt.Println(ts.ValueInFullTimes(dt, "low", math.NaN()))
	//fmt.Println(ts.ValueInFullTimes(dt, "close", math.NaN()))
	//fmt.Println(ts.ValueInFullTimes(dt, "volume", math.NaN()))
	//fmt.Println(ts.ValueInFullTimes(dt, "openinterest", math.NaN()))

	//}
}

func TestReadContract(t *testing.T) {
	t.Parallel()

	tables := []map[string]interface{}{
		map[string]interface{}{
			"contract": "nqh96",
			"err":      os.ErrNotExist,
		},
		map[string]interface{}{
			"contract": "nqh98",
			"err":      errEmptyFile,
		},
		map[string]interface{}{
			"contract": "nqh01",
			"err":      nil,
		},
		map[string]interface{}{
			"contract": "fxh20",
			"err":      nil,
		},
	}

	for _, table := range tables {
		c := continuousContract{}
		_, err := c.readContract(table["contract"].(string))
		if table["err"] != nil {
			if !errors.Is(err, table["err"].(error)) {
				t.Errorf("expect: %T, get %T", table["err"], err)
			}
		} else {
			if err != nil {
				t.Errorf("expect: %T, get %T", table["err"], err)
			}
		}
	}

}

func TestReadContractAll(t *testing.T) {
	t.Parallel()

	root := filepath.Join(
		os.Getenv("HOME"),
		"Documents",
		"data_source",
		"continuous",
	)

	c := continuousContract{}

	symbols, err := os.ReadDir(root)
	if err != nil {
		t.Errorf(err.Error())
	}

	for _, symbol := range symbols {
		path := filepath.Join(root, symbol.Name())

		contracts, err := os.ReadDir(path)
		if err != nil {
			t.Errorf(err.Error())
		}

		for _, contract := range contracts {
			cs := strings.ReplaceAll(contract.Name(), ".csv", "")

			ts, err := c.readContract(cs)
			if errors.Is(err, errEmptyFile) {
				fmt.Printf("empty file: %s\n", cs)
				continue
			}

			if err != nil {
				t.Errorf("contract: %s\nerr: %s", cs, err.Error())
			}

			if ts == nil || len(ts.FullTimes()) <= 0 {
				t.Errorf("contract: %s\ninvalid series", cs)
			}
		}
	}

}
