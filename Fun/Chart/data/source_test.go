package data

import (
	"fmt"
	"testing"
	"time"
)

func TestStockCharts(t *testing.T) {
	//t.SkipNow()
	t.Parallel()

	source := NewDataSource(StockCharts)

	s, _ := time.Parse("20060102", "20171230")
	e, _ := time.Parse("20060102", "20190101")

	_, err := source.Read(s, e, "rvx", Daily)
	if err != nil {
		t.Errorf("err should be nil but get %s", err)
	}

	//if len(ts.Times()) != 251 {
	//t.Errorf("expect 251 but get %d", len(ts.Times()))
	//}
}

func TestYahoo(t *testing.T) {
	//t.SkipNow()
	t.Parallel()

	src := NewDataSource(Yahoo)

	s, _ := time.Parse("20060102", "20170101")
	e, _ := time.Parse("20060102", "20190101")

	_, err := src.Read(s, e, "vix", Daily)
	if err != nil {
		t.Errorf("err should be nil but get %s", err)
	}

	_, err = src.Read(s, e, "vxn", Daily)
	if err != nil {
		t.Errorf("err should be nil but get %s", err)
	}

	//for _, q := range ts.SliceValues() {
	//fmt.Println(q.Time())
	//fmt.Println(q.Open())
	//fmt.Println(q.High())
	//fmt.Println(q.Low())
	//fmt.Println(q.Close())
	//fmt.Println(q.Volume())
	//fmt.Println()
	//}
}

func TestBarchart(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	src := NewDataSource(Barchart)

	s, _ := time.Parse("20060102", "20180101")
	e, _ := time.Parse("20060102", "20190101")

	//ts, err := src.Read(s, e, "es", Daily)
	_, err := src.Read(s, e, "no", Daily)
	//ts, err := src.Read(s, e, "vstx", Daily)
	if err != nil {
		t.Errorf("err should be nil but get %s", err)
	}

	//for i := 0; i < len(ts.Times()); i++ {
	//fmt.Println(ts.Times()[i])
	//fmt.Println(ts.ValueAtTimesIndex(i, "open", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "high", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "low", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "close", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "volume", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "openinterest", 0))
	//}
}

func TestAlphaVantage(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	src := NewDataSource(AlphaVantage)

	s, _ := time.Parse("20060102", "20171230")
	e, _ := time.Parse("20060102", "20190101")

	_, err := src.Read(s, e, "spx", Daily)
	if err != nil {
		t.Errorf("err should be nil but get %s", err)
	}

	//qs.TimeSlice(s, e)

	//fmt.Println(s)
	//fmt.Println(e)

	//fmt.Println(qs.SliceStartTime())
	//fmt.Println(qs.SliceEndTime())

	//for _, q := range qs.All() {
	//fmt.Println(q.t)
	//}
}
func TestInvestingCom(t *testing.T) {
	//t.SkipNow()
	t.Parallel()

	src := NewDataSource(InvestingCom)

	s, _ := time.Parse("20060102", "20180101")
	e, _ := time.Parse("20060102", "20190101")

	_, err := src.Read(s, e, "vstx", Daily)
	if err != nil {
		t.Errorf("err should be nil but get %s", err)
	}

	_, err = src.Read(s, e, "jniv", Daily)
	if err != nil {
		t.Errorf("err should be nil but get %s", err)
	}

	//for i := 0; i < len(ts.Times()); i++ {
	//fmt.Println(ts.Times()[i])
	//fmt.Println(ts.ValueAtTimesIndex(i, "open", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "high", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "low", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "close", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "volume", 0))
	//fmt.Println(ts.ValueAtTimesIndex(i, "openinterest", 0))
	//}
}

func TestQuandl(t *testing.T) {
	t.SkipNow()
	t.Parallel()

	src := NewDataSource(Quandl)

	s, _ := time.Parse("20060102", "19900101")
	e, _ := time.Parse("20060102", "20190101")

	ts, err := src.Read(s, e, "nk225mg2020", Daily)
	//ts, err := src.Read(s, e, "fesxu2016", Daily)
	//ts, err := src.Read(s, e, "fvsh2020", Daily)
	if err != nil {
		t.Errorf("err should be nil but get %s", err)
	}

	for i := 0; i < len(ts.FullTimes()); i++ {
		fmt.Println(ts.FullTimes()[i])
		fmt.Println(ts.ValueAtFullTimesIndex(i, "open", 0))
		fmt.Println(ts.ValueAtFullTimesIndex(i, "high", 0))
		fmt.Println(ts.ValueAtFullTimesIndex(i, "low", 0))
		fmt.Println(ts.ValueAtFullTimesIndex(i, "close", 0))
		fmt.Println(ts.ValueAtFullTimesIndex(i, "volume", 0))
		fmt.Println(ts.ValueAtFullTimesIndex(i, "openinterest", 0))
	}
}
