package utils

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/KushamiNeko/GoFun/Model/finance"
)

func randomTimeSeries(start, end, ex string) *finance.TimeSeries {
	now := time.Now()

	rand.Seed(now.Unix())

	times := make([]time.Time, 0)

	from, _ := time.Parse("20060102", start)
	to, _ := time.Parse("20060102", end)

	exclude := time.Time{}
	if ex != "" {
		exclude, _ = time.Parse("20060102", ex)
	}

	dt := from
	for {
		if dt.After(to) {
			break
		}

		if dt.Weekday() == time.Saturday || dt.Weekday() == time.Sunday {
			dt = dt.Add(24 * time.Hour)
			continue
		}

		if dt.Equal(exclude) {
			dt = dt.Add(24 * time.Hour)
			continue
		}

		times = append(times, dt)
		dt = dt.Add(24 * time.Hour)
	}

	tl := len(times)

	values := make(map[string][]float64)
	values["open"] = make([]float64, tl)
	values["high"] = make([]float64, tl)
	values["low"] = make([]float64, tl)
	values["close"] = make([]float64, tl)
	values["volume"] = make([]float64, tl)
	values["openinterest"] = make([]float64, tl)

	for i := range times {
		values["open"][i] = rand.Float64() * 10000
		values["high"][i] = rand.Float64() * 10000
		values["low"][i] = rand.Float64() * 10000
		values["close"][i] = rand.Float64() * 10000
		values["volume"][i] = rand.Float64() * 10000
		values["openinterest"][i] = rand.Float64() * 10000
	}

	series := finance.NewTimeSeries(times, values)

	return series
}

func TestDailyToWeekly(t *testing.T) {
	tables := []map[string]interface{}{
		map[string]interface{}{
			"start":   "20200101",
			"end":     "20200116",
			"exclude": "",
			"expect": []map[string]interface{}{
				map[string]interface{}{
					"si": 0,
					"ei": 2,
				},
				map[string]interface{}{
					"si": 3,
					"ei": 7,
				},
				map[string]interface{}{
					"si": 8,
					"ei": 11,
				},
			},
		},
		map[string]interface{}{
			"start":   "20200103",
			"end":     "20200117",
			"exclude": "",
			"expect": []map[string]interface{}{
				map[string]interface{}{
					"si": 0,
					"ei": 0,
				},
				map[string]interface{}{
					"si": 1,
					"ei": 5,
				},
				map[string]interface{}{
					"si": 6,
					"ei": 10,
				},
			},
		},
		map[string]interface{}{
			"start":   "20191230",
			"end":     "20200113",
			"exclude": "",
			"expect": []map[string]interface{}{
				map[string]interface{}{
					"si": 0,
					"ei": 4,
				},
				map[string]interface{}{
					"si": 5,
					"ei": 9,
				},
				map[string]interface{}{
					"si": 10,
					"ei": 10,
				},
			},
		},
		map[string]interface{}{
			"start":   "20191230",
			"end":     "20200113",
			"exclude": "20200108",
			"expect": []map[string]interface{}{
				map[string]interface{}{
					"si": 0,
					"ei": 4,
				},
				map[string]interface{}{
					"si": 5,
					"ei": 8,
				},
				map[string]interface{}{
					"si": 9,
					"ei": 9,
				},
			},
		},
		map[string]interface{}{
			"start":   "20191230",
			"end":     "20200114",
			"exclude": "20200110",
			"expect": []map[string]interface{}{
				map[string]interface{}{
					"si": 0,
					"ei": 4,
				},
				map[string]interface{}{
					"si": 5,
					"ei": 8,
				},
				map[string]interface{}{
					"si": 9,
					"ei": 10,
				},
			},
		},
		map[string]interface{}{
			"start":   "20191231",
			"end":     "20200114",
			"exclude": "20200113",
			"expect": []map[string]interface{}{
				map[string]interface{}{
					"si": 0,
					"ei": 3,
				},
				map[string]interface{}{
					"si": 4,
					"ei": 8,
				},
				map[string]interface{}{
					"si": 9,
					"ei": 9,
				},
			},
		},
		map[string]interface{}{
			"start":   "20191231",
			"end":     "20200115",
			"exclude": "20200113",
			"expect": []map[string]interface{}{
				map[string]interface{}{
					"si": 0,
					"ei": 3,
				},
				map[string]interface{}{
					"si": 4,
					"ei": 8,
				},
				map[string]interface{}{
					"si": 9,
					"ei": 10,
				},
			},
		},
	}

	for _, table := range tables {
		series := randomTimeSeries(table["start"].(string), table["end"].(string), table["exclude"].(string))
		weekly := DailyToWeekly(series)

		//for i := range series.FullTimes() {
		//fmt.Println(series.FullTimes()[i])
		//}

		//fmt.Println()

		//for i := range weekly.FullTimes() {
		//fmt.Println(weekly.FullTimes()[i])
		//}

		//fmt.Println()

		for i := range weekly.FullTimes() {
			dt := weekly.FullTimes()[i]

			esi := table["expect"].([]map[string]interface{})[i]["si"].(int)
			eei := table["expect"].([]map[string]interface{})[i]["ei"].(int)

			wt := weekly.FullTimes()[i].Format("20060102")
			et := series.FullTimes()[esi].Format("20060102")

			if wt != et {
				t.Errorf("time: %s, expect: %s, get: %s", dt, et, wt)
			}

			wo := weekly.Values("open")[i]
			eo := series.Values("open")[esi]

			if wo != eo {
				t.Errorf("open time: %s, expect: %f, get: %f", dt, eo, wo)
			}

			wc := weekly.Values("close")[i]
			ec := series.Values("close")[eei]

			if wc != ec {
				t.Errorf("close time: %s, expect: %f, get: %f", dt, ec, wc)
			}

			wh := weekly.Values("high")[i]
			wl := weekly.Values("low")[i]
			eh, el := weeklyHL(series, esi, eei)

			if wh != eh {
				t.Errorf("high time: %s, expect: %f, get: %f", dt, eh, wh)
			}

			if wl != el {
				t.Errorf("low time: %s, expect: %f, get: %f", dt, el, wl)
			}

			wv := weekly.Values("volume")[i]
			woi := weekly.Values("openinterest")[i]
			ev, eoi := weeklyVO(series, esi, eei)

			if wv != ev {
				t.Errorf("volume time: %s, expect: %f, get: %f", dt, ev, wv)
			}

			if woi != eoi {
				t.Errorf("open interest time: %s, expect: %f, get: %f", dt, eoi, woi)
			}

			//fmt.Println(dt)
			//fmt.Println(wo)
			//fmt.Println(wh)
			//fmt.Println(wl)
			//fmt.Println(wc)
			//fmt.Println(wv)
			//fmt.Println(woi)
		}
	}
}

func weeklyHL(ts *finance.TimeSeries, si, ei int) (float64, float64) {
	h := math.Inf(-1)
	l := math.Inf(1)

	for i := si; i <= ei; i++ {
		h = math.Max(ts.Values("high")[i], h)
		l = math.Min(ts.Values("low")[i], l)
	}

	return h, l
}

func weeklyVO(ts *finance.TimeSeries, si, ei int) (float64, float64) {
	var v, o float64

	for i := si; i <= ei; i++ {
		v += ts.Values("volume")[i]
		o += ts.Values("openinterest")[i]
	}

	return v, o
}
