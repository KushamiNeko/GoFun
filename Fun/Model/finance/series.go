package finance

import (
	"fmt"
	"strings"
	"time"
)

type TimeSeries struct {
	si int
	ei int

	times  []time.Time
	values map[string][]float64
}

func NewTimeSeries(times []time.Time, values map[string][]float64) *TimeSeries {
	tl := len(times)
	if tl <= 0 {
		panic("empty times")
	}

	for _, v := range values {
		if len(v) != tl {
			panic("the length of values must be equal to the length of indexes")
		}
	}

	return &TimeSeries{
		times:  times,
		values: values,
		si:     0,
		ei:     len(times) - 1,
	}
}

func (t *TimeSeries) FullTimes() []time.Time {
	return t.times
}

func (t *TimeSeries) FullValues(col string) []float64 {
	return t.values[col]
}

func (t *TimeSeries) SetColumn(col string, val []float64) {
	if len(val) != len(t.times) {
		panic(
			fmt.Sprintf(
				"the length of the val should be the same as the time series\nval %d, t: %d",
				len(val),
				len(t.times),
			),
		)
	}

	t.values[col] = val
}

func (t *TimeSeries) Times() []time.Time {
	return t.times[t.si : t.ei+1]
}

func (t *TimeSeries) Values(col string) []float64 {
	return t.values[col][t.si : t.ei+1]
}

func (t *TimeSeries) IndexInFullTimes(dt time.Time) int {
	for i, v := range t.times {
		if v.Equal(dt) {
			return i
		}
	}

	return -1
}

func (t *TimeSeries) IndexInTimes(dt time.Time) int {
	for i, v := range t.times[t.si : t.ei+1] {
		if v.Equal(dt) {
			return i
		}
	}

	return -1
}

func (t *TimeSeries) StartIndex() int {
	return t.si
}

func (t *TimeSeries) EndIndex() int {
	return t.ei
}

func (t *TimeSeries) StartTime() time.Time {
	return t.times[t.si]
}

func (t *TimeSeries) EndTime() time.Time {
	return t.times[t.ei]
}

func (t *TimeSeries) StartValue(col string) float64 {
	return t.values[col][t.si]
}

func (t *TimeSeries) EndValue(col string) float64 {
	return t.values[col][t.ei]
}

func (t *TimeSeries) ValueInFullTimes(index time.Time, col string, def float64) float64 {
	vs, ok := t.values[col]
	if ok {
		for i, v := range t.times {
			if v.Equal(index) {
				return vs[i]
			}
		}
	}

	return def
}

func (t *TimeSeries) ValueInTimes(index time.Time, col string, def float64) float64 {
	vs, ok := t.values[col]
	if ok {
		for i, v := range t.times[t.si : t.ei+1] {
			if v.Equal(index) {
				return vs[t.si : t.ei+1][i]
			}
		}
	}

	return def
}

func (t *TimeSeries) ValueAtFullTimesIndex(index int, col string, def float64) float64 {
	if vs, ok := t.values[col]; ok {
		return vs[index]
	}

	return def
}

func (t *TimeSeries) ValueAtTimesIndex(index int, col string, def float64) float64 {
	if vs, ok := t.values[col]; ok {
		return vs[t.si : t.ei+1][index]
	}

	return def
}

func (t *TimeSeries) All() {
	t.si = 0
	t.ei = len(t.times) - 1
}

func (t *TimeSeries) LastNSessions(n int) *TimeSeries {
	t.All()

	t.si = len(t.values) - n
	t.ei = len(t.values) - 1

	return t
}

func (t TimeSeries) IndexSlice(start, end int) {
	t.si = start
	t.ei = end
}

func (t *TimeSeries) TimeSlice(start, end time.Time) *TimeSeries {
	t.All()

	si := -1
	ei := -1

	if start.Before(t.times[0]) {
		si = 0
	}

	if end.Before(t.times[0]) {
		ei = 0
	}

	if start.After(t.times[len(t.times)-1]) {
		si = len(t.times) - 1
	}

	if end.After(t.times[len(t.times)-1]) {
		ei = len(t.times) - 1
	}

	for i, v := range t.times {
		if si != -1 && ei != -1 {
			break
		}

		if (v.After(start) || v.Equal(start)) && si == -1 {
			si = i
		}

		if ei == -1 {
			if v.After(end) {
				ei = i - 1
			}

			if v.Equal(end) {
				ei = i
			}
		}
	}

	t.si = si
	t.ei = ei

	return t
}

func (t *TimeSeries) Forward() *TimeSeries {
	if (t.ei + 1) < len(t.times) {
		t.si += 1
		t.ei += 1
	}

	return t
}

func (t *TimeSeries) Backward() *TimeSeries {
	if (t.si - 1) > 0 {
		t.si -= 1
		t.ei -= 1
	}

	return t
}

func (t *TimeSeries) CSV() string {
	out := make([]string, 0, len(t.times)+1)
	out = append(out, "date,open,high,low,close,volume,openinterest")

	for i := 0; i < len(t.times); i++ {
		out = append(
			out,
			fmt.Sprintf(
				"%s,%f,%f,%f,%f,%f,%f",
				t.times[i].Format(`2006-01-02`),
				t.values["open"][i],
				t.values["high"][i],
				t.values["low"][i],
				t.values["close"][i],
				t.values["volume"][i],
				t.values["openinterest"][i],
			),
		)
	}

	return strings.Join(out, "\n")
}
