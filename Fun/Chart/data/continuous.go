package data

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/KushamiNeko/GoFun/Chart/futures"
	"github.com/KushamiNeko/GoFun/Chart/utils"
	"github.com/KushamiNeko/GoFun/Model/finance"
	"github.com/KushamiNeko/GoFun/Utility/pretty"
)

var errDataFormat error = fmt.Errorf("invalid data format")
var errEmptyFile error = fmt.Errorf("empty file")

type continuousContract struct{}

func (c continuousContract) readContract(contract string) (*finance.TimeSeries, error) {
	root := filepath.Join(
		os.Getenv("HOME"),
		"Documents/data_source/continuous",
	)

	body, err := os.ReadFile(
		filepath.Join(
			root,
			contract[:2],
			fmt.Sprintf("%s.csv", contract),
		),
	)
	if err != nil {
		return nil, err
	}

	if len(bytes.TrimSpace(body)) == 0 {
		return nil, errEmptyFile
	}

	var series *finance.TimeSeries

	series, err = c.barchartHistorical(bytes.NewBuffer(body))
	if err == nil {
		goto re
	}

	series, err = c.barchartInteractive(bytes.NewBuffer(body))
	if err == nil {
		goto re
	}

	series, err = c.barchartOnDemand(bytes.NewBuffer(body))
	if err == nil {
		goto re
	}

	return nil, errDataFormat

re:
	return series, nil
}

func (c continuousContract) barchartInteractive(buffer *bytes.Buffer) (*finance.TimeSeries, error) {
	//"Date Time","Symbol","Open","High","Low","Close","Change"
	//2000-01-03,ESH00,1549.5,1556.5,1512.75,1527.25,0
	//const regex string = `\s*(\d{4}-\d{2}-\d{2})\s*,\s*([\w\d]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([-0-9.]+)\s*`

	//re := regexp.MustCompile(regex)
	//match := re.FindAllStringSubmatch(buffer.String(), -1)

	bl := len(strings.Split(buffer.String(), "\n"))
	if bl <= 3 {
		return nil, errEmptyFile
	}

	reader := csv.NewReader(buffer)
	match := make([][]string, 0, bl)
	for {
		record, err := reader.Read()
		if errors.Is(err, csv.ErrFieldCount) || errors.Is(err, io.EOF) {
			break
		}

		match = append(match, record)
	}

	ml := len(match)
	if ml <= 3 {
		return nil, errDataFormat
	}

	match = match[2:]
	ml -= 2

	ts, vs := initTimeValues(ml)

	const timeFormat = `2006-01-02`

	for i, m := range match {
		//setTimeValuesRow(timeFormat, i, ts, vs, m[1], m[3], m[4], m[5], m[6], "", "")
		err := setTimeValuesRow(timeFormat, i, ts, vs, m[0], m[2], m[3], m[4], m[5], "", "")
		if err != nil {
			return nil, err
		}
	}

	series := finance.NewTimeSeries(ts, vs)

	return series, nil
}

func (c continuousContract) barchartHistorical(buffer *bytes.Buffer) (*finance.TimeSeries, error) {
	//Symbol,Time,Open,High,Low,Last,Change,Volume,"Open Int"
	//const regex string = `(?:\s*([a-zA-Z0-9]+)\s*,)*\s*(\d{2}/\d{2}/\d{4})\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([-0-9.]+)\s*,\s*([0-9]+)\s*(?:,\s*([0-9]+))*\s*`

	//re := regexp.MustCompile(regex)
	//match := re.FindAllStringSubmatch(buffer.String(), -1)

	bl := len(strings.Split(buffer.String(), "\n"))
	if bl <= 2 {
		return nil, errEmptyFile
	}

	reader := csv.NewReader(buffer)
	match := make([][]string, 0, bl)
	for {
		record, err := reader.Read()
		if errors.Is(err, csv.ErrFieldCount) || errors.Is(err, io.EOF) {
			break
		}

		match = append(match, record)
	}

	ml := len(match)
	if ml <= 2 {
		return nil, errDataFormat
	}

	match = match[1:]
	ml -= 1

	ts, vs := initTimeValues(ml)

	for i, m := range match {
		//err := setTimeValuesRow(timeFormat, ml-1-i, ts, vs, m[1], m[2], m[3], m[4], m[5], m[7], m[8])
		err := setTimeValuesRow(`01/02/2006`, ml-1-i, ts, vs, m[0], m[1], m[2], m[3], m[4], m[6], m[7])
		if err != nil {
			err := setTimeValuesRow(`01/02/06`, ml-1-i, ts, vs, m[0], m[1], m[2], m[3], m[4], m[6], m[7])
			if err != nil {
				return nil, err
			}
		}
	}

	series := finance.NewTimeSeries(ts, vs)

	return series, nil
}

func (c continuousContract) barchartOnDemand(buffer *bytes.Buffer) (*finance.TimeSeries, error) {
	//"esu19","2018-01-02T00:00:00-06:00","2018-01-02","2692.5","2713.25","2691.75","2710.25","1003194","3051095"
	//const pattern = `"([a-zA-Z0-9]+)","(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{2}:\d{2})","(\d{4}-\d{2}-\d{2})","([0-9.]+)","([0-9.]+)","([0-9.]+)","([0-9.]+)","(\d+)"(?:,"(\d+)")*`

	//regex := regexp.MustCompile(pattern)
	//match := regex.FindAllStringSubmatch(buffer.String(), -1)

	bl := len(strings.Split(buffer.String(), "\n"))
	if bl <= 2 {
		return nil, errEmptyFile
	}

	reader := csv.NewReader(buffer)
	match := make([][]string, 0, bl)
	for {
		record, err := reader.Read()
		if errors.Is(err, csv.ErrFieldCount) || errors.Is(err, io.EOF) {
			break
		}

		match = append(match, record)
	}

	ml := len(match)
	if ml <= 2 {
		return nil, errDataFormat
	}

	contract := match[0]
	match = match[1:]
	ml -= 1

	ts, vs := initTimeValues(ml)

	for i, m := range match {

		d := strings.ReplaceAll(m[2], `"`, "")
		o := strings.ReplaceAll(m[3], `"`, "")
		h := strings.ReplaceAll(m[4], `"`, "")
		l := strings.ReplaceAll(m[5], `"`, "")
		c := strings.ReplaceAll(m[6], `"`, "")
		v := strings.ReplaceAll(m[7], `"`, "")
		oi := strings.ReplaceAll(m[8], `"`, "")

		regex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})-\d{2}:\d{2}`)
		ld := regex.FindAllStringSubmatch(strings.ReplaceAll(m[1], `"`, ""), -1)[0][1]
		ldt, err := time.Parse(`2006-01-02T15:04:05`, ld)
		if err != nil {
			return nil, err
		}

		if ldt.Hour() != 0 {
			pretty.ColorPrintln(
				pretty.PaperRed400,
				fmt.Sprintf("skip unusual intra-day data in contract %s at: %s", contract, ldt),
			)
			continue
		}

		if oi == "" {
			oi = "0"
		}

		err = setTimeValuesRow(`2006-01-02`, i, ts, vs, d, o, h, l, c, v, oi)
		if err != nil {
			return nil, err
		}
	}

	series := finance.NewTimeSeries(ts, vs)

	return series, nil
}

func (c continuousContract) Read(
	from,
	to time.Time,
	symbol string,
	freq Frequency,
) (*finance.TimeSeries, error) {

	var months futures.ContractMonths
	switch symbol {
	case "gc":
		months = futures.EvenContractMonths
	case "cl":
		months = futures.AllContractMonths
	default:
		months = futures.FinancialContractMonths
	}

	var roll futures.RollingMethod
	switch symbol {
	case "gc":
		roll = futures.OpenInterest
	case "cl":
		roll = futures.OpenInterest
	default:
		roll = futures.LastNTradingDay
	}

	ts, err := c.read(
		from,
		to,
		symbol,
		months,
		futures.BarchartSymbolFormat,
		roll,
		futures.Ratio,
	)
	if err != nil {
		return nil, err
	}

	switch freq {
	case Hourly:
		panic("unimplemented")
	case Daily:
	case Weekly:
		ts = utils.DailyToWeekly(ts)
	case Monthly:
		panic("unimplemented")
	default:
		panic("unknown frequency")
	}

	return ts, nil
}

func (cct continuousContract) read(
	from,
	to time.Time,
	symbol string,
	contractMonths futures.ContractMonths,
	format futures.SymbolFormat,
	roll futures.RollingMethod,
	adjust futures.AdjustingMethod,
) (*finance.TimeSeries, error) {

	contracts := futures.ContractList(from, to, symbol, contractMonths, format)

	if len(contracts) <= 0 {
		return nil, fmt.Errorf("empyt contracts")
	}

	contractSeries := make([]*finance.TimeSeries, len(contracts))

	for i, c := range contracts {
		ts, err := cct.readContract(c)
		if err != nil {

			switch {
			case errors.Is(err, os.ErrNotExist):
				continue
			case errors.Is(err, errDataFormat):
				continue
			case errors.Is(err, errEmptyFile):
				continue
			default:
				return nil, err
			}

		}

		contractSeries[i] = ts
	}

	contractTimes := make([]time.Time, 0, int(to.Sub(from).Hours()/24.0)+1)

	for _, cs := range contractSeries {

		if cs == nil {
			//panic("empty series")
			continue
		}

		for _, t := range cs.FullTimes() {

			if t.After(to) || t.Before(from) {
				continue
			}

			inTimes := false
			for _, ct := range contractTimes {
				if ct.Equal(t) {
					inTimes = true
					break
				}
			}

			if !inTimes {
				contractTimes = append(contractTimes, t)
			}
		}

	}

	sort.Slice(contractTimes, func(i, j int) bool {
		return contractTimes[i].Before(contractTimes[j])
	})

	ctl := len(contractTimes)

	if ctl <= 0 {
		return nil, fmt.Errorf("empty contracts")
	}

	t := make([]time.Time, 0, ctl)
	o := make([]float64, 0, ctl)
	h := make([]float64, 0, ctl)
	l := make([]float64, 0, ctl)
	c := make([]float64, 0, ctl)
	v := make([]float64, 0, ctl)
	oi := make([]float64, 0, ctl)

	cl := len(contracts)

	currentIndex := 0

	current := futures.NewContract(
		contracts[currentIndex],
		contractSeries[currentIndex],
		format,
	)

	previous := futures.NewContract(
		contracts[currentIndex+1],
		contractSeries[currentIndex+1],
		format,
	)

	var adjustment float64
	switch adjust {
	case futures.PanamaCanal:
		adjustment = 0
	case futures.Ratio:
		adjustment = 1
	default:
		panic("unknown adjusting method")
	}

	for i := ctl - 1; i >= 0; i-- {

		dt := contractTimes[i]

		var rolling bool

		if previous.Series() != nil {

			switch roll {
			case futures.LastNTradingDay:
				rolling = futures.ContractRollLastNTradingDay(dt, current, previous, 4)
			case futures.FirstOfMonth:
				rolling = futures.ContractRollFirstOfMonth(dt, current, previous)
			case futures.OpenInterest:
				rolling = futures.ContractRollOpenInterest(dt, current, previous)
			default:
				panic("unknown rolling method")
			}

		}

		// if you put rolling here
		// the data on the rolling date will be the front contract
		// ex:
		// front contract: esu19
		// back contract: esz19
		// rolling date: 0901
		// the data at 0901 will be esu19

		cvo := current.Series().ValueInFullTimes(dt, "open", math.NaN())
		cvh := current.Series().ValueInFullTimes(dt, "high", math.NaN())
		cvl := current.Series().ValueInFullTimes(dt, "low", math.NaN())
		cvc := current.Series().ValueInFullTimes(dt, "close", math.NaN())
		cvv := current.Series().ValueInFullTimes(dt, "volume", math.NaN())
		cvoi := current.Series().ValueInFullTimes(dt, "openinterest", math.NaN())

		if math.IsNaN(cvo) || math.IsNaN(cvh) || math.IsNaN(cvl) || math.IsNaN(cvc) {
			pretty.ColorPrintln(pretty.PaperRed400, fmt.Sprintf("nan value at time: %s", dt))
			continue
		}

		t = append(t, dt)
		v = append(v, cvv)
		oi = append(oi, cvoi)

		switch adjust {
		case futures.PanamaCanal:

			o = append(o, cvo+adjustment)
			h = append(h, cvh+adjustment)
			l = append(l, cvl+adjustment)
			c = append(c, cvc+adjustment)

		case futures.Ratio:

			o = append(o, cvo*adjustment)
			h = append(h, cvh*adjustment)
			l = append(l, cvl*adjustment)
			c = append(c, cvc*adjustment)

		default:
			panic("unknown adjusting method")
		}

		// if you put rolling here
		// the data on the rolling date will be the back contract
		// ex:
		// front contract: esu19
		// back contract: esz19
		// rolling date: 0901
		// the data at 0901 will be esz19

		if rolling {
			currentIndex += 1

			if currentIndex == cl {
				break
			}

			cc := current.Series().ValueInFullTimes(dt, "close", math.NaN())
			pc := previous.Series().ValueInFullTimes(dt, "close", math.NaN())

			switch adjust {
			case futures.PanamaCanal:
				adjustment += cc - pc
			case futures.Ratio:
				adjustment *= cc / pc
			default:
				panic("unknown adjusting method")
			}

			current = futures.NewContract(
				contracts[currentIndex],
				contractSeries[currentIndex],
				format,
			)

			if currentIndex < cl-1 {

				previous = futures.NewContract(
					contracts[currentIndex+1],
					contractSeries[currentIndex+1],
					format,
				)
			}
		}

	}

	tl := len(t)
	nt := make([]time.Time, tl)
	no := make([]float64, tl)
	nh := make([]float64, tl)
	nl := make([]float64, tl)
	nc := make([]float64, tl)
	nv := make([]float64, tl)
	noi := make([]float64, tl)

	for i := 0; i < tl; i++ {
		nt[tl-1-i] = t[i]
		no[tl-1-i] = o[i]
		nh[tl-1-i] = h[i]
		nl[tl-1-i] = l[i]
		nc[tl-1-i] = c[i]
		nv[tl-1-i] = v[i]
		noi[tl-1-i] = oi[i]
	}

	cc := finance.NewTimeSeries(
		nt,
		map[string][]float64{
			"open":         no,
			"high":         nh,
			"low":          nl,
			"close":        nc,
			"volume":       nv,
			"openinterest": noi,
		},
	)

	return cc, nil
}
