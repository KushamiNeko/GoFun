package utils

import (
	"math"
	"time"

	"github.com/KushamiNeko/GoFun/Model/finance"
)

func DailyToWeekly(ts *finance.TimeSeries) *finance.TimeSeries {

	ts.All()

	wdCur := math.MinInt64

	var si, ei int

	times := make([]time.Time, 0)
	os := make([]float64, 0)
	hs := make([]float64, 0)
	ls := make([]float64, 0)
	cs := make([]float64, 0)
	vs := make([]float64, 0)
	ois := make([]float64, 0)

	var end time.Time

	for i, t := range ts.Times() {
		dint := weekdayToInt(t.Weekday())

		if dint < wdCur {

			end = t

			t, o, h, l, c, v, oi := weeklyRange(ts, si, ei)

			times = append(times, t)
			os = append(os, o)
			hs = append(hs, h)
			ls = append(ls, l)
			cs = append(cs, c)
			vs = append(vs, v)
			ois = append(ois, oi)

			si = i
		}

		wdCur = dint
		ei = i
	}

	t, o, h, l, c, v, oi := weeklyRange(
		ts,
		ts.IndexInFullTimes(end),
		len(ts.FullTimes())-1,
	)

	times = append(times, t)
	os = append(os, o)
	hs = append(hs, h)
	ls = append(ls, l)
	cs = append(cs, c)
	vs = append(vs, v)
	ois = append(ois, oi)

	wts := finance.NewTimeSeries(
		times,
		map[string][]float64{
			"open":         os,
			"high":         hs,
			"low":          ls,
			"close":        cs,
			"volume":       vs,
			"openinterest": ois,
		},
	)

	return wts
}

func weeklyRange(ts *finance.TimeSeries, si, ei int) (
	t time.Time, o float64, h float64, l float64, c float64, v float64, oi float64) {

	if ei < si {
		panic("invalid index")
	}

	h = math.Inf(-1)
	l = math.Inf(1)

	t = ts.Times()[si]
	o = ts.Values("open")[si]
	c = ts.Values("close")[ei]

	for i := si; i <= ei; i++ {
		//if !math.IsNaN(ts.Values("high")[i]) {
		h = math.Max(h, ts.Values("high")[i])
		//}

		//if !math.IsNaN(ts.Values("low")[i]) {
		l = math.Min(l, ts.Values("low")[i])
		//}

		//if !math.IsNaN(ts.ValueAtTimesIndex(i, "volume", 0)) {
		v += ts.ValueAtTimesIndex(i, "volume", 0)
		//}

		//if !math.IsNaN(ts.ValueAtTimesIndex(i, "openinterest", 0)) {
		oi += ts.ValueAtTimesIndex(i, "openinterest", 0)
		//}
	}

	return
}

func weekdayToInt(t time.Weekday) int {
	switch t {
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	case time.Saturday:
		return 6
	case time.Sunday:
		return 7
	default:
		panic("unknown weekday")
	}
}
