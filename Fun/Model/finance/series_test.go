package finance

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestTimeSlice(t *testing.T) {
	var err error

	loc, err := time.LoadLocation("")
	if err != nil {
		t.Errorf("err should be nil, but get: %s", err)
	}

	ts := &TimeSeries{
		times: []time.Time{
			time.Date(2018, time.January, 1, 0, 0, 0, 0, loc),
			time.Date(2018, time.January, 3, 0, 0, 0, 0, loc),
			time.Date(2018, time.January, 6, 0, 0, 0, 0, loc),
			time.Date(2018, time.January, 10, 0, 0, 0, 0, loc),
			time.Date(2018, time.January, 15, 0, 0, 0, 0, loc),
		},
		values: map[string][]float64{
			"open": []float64{
				100, 300, 500, 500, 500,
			},
			"high": []float64{
				105, 305, 505, 505, 505,
			},
			"low": []float64{
				95, 295, 495, 495, 495,
			},
			"close": []float64{
				97, 297, 497, 497, 497,
			},
			"volume": []float64{
				1000, 3000, 5000, 5000, 5000,
			},
			"openinterest": []float64{
				1000, 3000, 5000, 5000, 5000,
			},
		},
	}

	ts.All()

	tbs := []map[string]string{
		map[string]string{
			"start":             "20180101",
			"end":               "20180106",
			"expect si":         "0",
			"expect ei":         "2",
			"expect all len":    "5",
			"expect slice len":  "3",
			"expect start time": "20180101",
			"expect end time":   "20180106",
		},
		map[string]string{
			"start":             "20180102",
			"end":               "20180106",
			"expect si":         "1",
			"expect ei":         "2",
			"expect all len":    "5",
			"expect slice len":  "2",
			"expect start time": "20180103",
			"expect end time":   "20180106",
		},
		map[string]string{
			"start":             "20180103",
			"end":               "20180115",
			"expect si":         "1",
			"expect ei":         "4",
			"expect all len":    "5",
			"expect slice len":  "4",
			"expect start time": "20180103",
			"expect end time":   "20180115",
		},
		map[string]string{
			"start":             "20180108",
			"end":               "20180116",
			"expect si":         "3",
			"expect ei":         "4",
			"expect all len":    "5",
			"expect slice len":  "2",
			"expect start time": "20180110",
			"expect end time":   "20180115",
		},
	}

	for _, tb := range tbs {
		s, _ := time.Parse("20060102", tb["start"])
		e, _ := time.Parse("20060102", tb["end"])

		ts.TimeSlice(s, e)

		if strings.Compare(fmt.Sprintf("%d", ts.si), tb["expect si"]) != 0 {
			t.Errorf("expect %s but get %d", tb["expect si"], ts.si)
		}

		if strings.Compare(fmt.Sprintf("%d", ts.ei), tb["expect ei"]) != 0 {
			t.Errorf("expect %s but get %d", tb["expect ei"], ts.ei)
		}

		if strings.Compare(fmt.Sprintf("%d", len(ts.FullTimes())), tb["expect all len"]) != 0 {
			t.Errorf("expect %s but get %d", tb["expect all len"], len(ts.FullTimes()))
		}

		if strings.Compare(fmt.Sprintf("%d", len(ts.Times())), tb["expect slice len"]) != 0 {
			t.Errorf("expect %s but get %d, %v", tb["expect slice len"], len(ts.Times()), ts.Times())
		}

		if strings.Compare(ts.StartTime().Format("20060102"), tb["expect start time"]) != 0 {
			t.Errorf("expect %s but get %s", tb["expect start time"], ts.StartTime().Format("20060102"))
		}

		if strings.Compare(ts.EndTime().Format("20060102"), tb["expect end time"]) != 0 {
			t.Errorf("expect %s but get %s", tb["expect end time"], ts.EndTime().Format("20060102"))
		}
	}
}
