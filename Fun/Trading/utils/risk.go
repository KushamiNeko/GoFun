package utils

import (
	"fmt"
	"math"
	"time"

	"github.com/KushamiNeko/GoFun/Fun/Model/finance"
	"github.com/KushamiNeko/GoFun/Fun/Trading/model"
)

func transactionSlice(
	records []*model.FuturesTransaction,
	from,
	to time.Time,
	op string,
) []*model.FuturesTransaction {

	ts := make([]*model.FuturesTransaction, 0)

	sliced := false
	for _, r := range records {
		if (r.Time().Equal(from) || r.Time().After(from)) && !sliced {
			sliced = true
		}

		if (r.Time().After(to)) && sliced {
			sliced = false
		}

		if sliced {
			if r.Operation() == op {
				ts = append(ts, r)
			}
		}
	}

	if len(ts) > 1 && ts[0].Time().Equal(ts[1].Time()) {
		ts = ts[1:]
	}

	if len(ts) > 1 && ts[len(ts)-1].Time().Equal(ts[len(ts)-2].Time()) {
		ts = ts[:len(ts)-1]
	}

	return ts
}

type Risk struct {
	time time.Time
	risk float64

	adjucted bool
	combined bool
}

func (r *Risk) String() string {
	var l string
	if r.combined {
		l = "Combined Risk"
	} else {
		l = "Risk"
	}

	if r.adjucted {
		l = fmt.Sprintf("%s (Adj)", l)
	} else {
		l = fmt.Sprintf("%s (Sim)", l)
	}

	return fmt.Sprintf("%-15s @ %s: %.4f%%", l, r.time.Format("20060102"), r.risk)
}

func (r *Risk) Label() string {
	var l string
	if r.combined {
		l = "Combined Risk"
	} else {
		l = "Risk"
	}

	if r.adjucted {
		l = fmt.Sprintf("%s(Adj)", l)
	} else {
		l = fmt.Sprintf("%s(Sim)", l)
	}

	return l
}

func (r *Risk) Time() time.Time {
	return r.time
}

func (r *Risk) Risk() float64 {
	return r.risk
}

func (r *Risk) Adjusted() bool {
	return r.adjucted
}

func (r *Risk) Combined() bool {
	return r.combined
}

func CalculateRisk(
	series *finance.TimeSeries,
	records []*model.FuturesTransaction,
	from,
	to,
	flip time.Time,
	//ct *model.FuturesTransaction,
	ct time.Time,
	op string,
	adjusted bool,
) []*Risk {

	fi := series.IndexInTimes(flip)

	//ci := series.IndexInTimes(ct.Time())
	ci := series.IndexInTimes(ct)

	ts := transactionSlice(records, from, to, op)

	risks := make([]*Risk, 0)

	positions := make([]float64, 0)
	sizes := make([]float64, 0)

	for _, record := range ts {

		var r float64

		index := series.IndexInTimes(record.Time())
		if index == -1 {
			panic(fmt.Errorf("unknown index"))
		}

		c := series.ValueAtTimesIndex(index, "close", 0)

		var (
			nl float64 = math.Inf(1)
			nh float64 = math.Inf(-1)
		)

		gain := 0.0

		for i := 1; i <= 14; i++ {

			cic := series.ValueAtTimesIndex(index+i, "close", 0)

			switch op {
			case "+":
				gain = math.Max(((cic-c)/c)*100.0, gain)
			case "-":
				gain = math.Max(((-cic+c)/c)*100.0, gain)
			default:
				panic("unknown operation")
			}

			if index+i != ci && index+1 != fi {

				l := series.ValueAtTimesIndex(index+i, "low", 0)
				nl = math.Min(
					nl,
					math.Max(
						l,
						c*0.975,
					),
				)

				h := series.ValueAtTimesIndex(index+i, "high", 0)
				nh = math.Max(
					nh,
					math.Min(
						h,
						c*1.025,
					),
				)

			} else {

				cio := series.ValueAtTimesIndex(index+i, "open", 0)
				cih := series.ValueAtTimesIndex(index+i, "high", 0)
				cil := series.ValueAtTimesIndex(index+i, "low", 0)

				if gain >= 1.5 {

					nl = math.Min(
						nl,
						math.Max(
							math.Min(
								c,
								cio,
							),
							math.Max(
								c*0.975,
								cil,
							),
						),
					)

					nh = math.Max(
						nh,
						math.Min(
							math.Max(
								c,
								cio,
							),
							math.Min(
								c*1.025,
								cih,
							),
						),
					)

				} else {

					l := series.ValueAtTimesIndex(index+i, "low", 0)
					nl = math.Min(
						nl,
						math.Max(
							l,
							c*0.975,
						),
					)

					h := series.ValueAtTimesIndex(index+i, "high", 0)
					nh = math.Max(
						nh,
						math.Min(
							h,
							c*1.025,
						),
					)

				}

			}

			if index+i >= fi || index+i >= ci {
				break
			}
		}

		var q float64
		if adjusted {
			q = float64(record.Quantity())
		} else {
			q = 1
		}

		switch op {
		case "+":
			r += (((nl - c) / c) * 100.0) * q

			risks = append(risks, &Risk{
				time:     record.Time(),
				risk:     r,
				adjucted: adjusted,
				combined: false,
			})

			if len(positions) > 0 {
				for i, p := range positions {
					r += (((nl - p) / p) * 100.0) * sizes[i]
				}

				risks = append(risks, &Risk{
					time:     record.Time(),
					risk:     r,
					adjucted: adjusted,
					combined: true,
				})
			} else {
				risks = append(risks, &Risk{
					time:     record.Time(),
					risk:     r,
					adjucted: adjusted,
					combined: true,
				})
			}

		case "-":
			r += (((c - nh) / c) * 100.0) * q

			risks = append(risks, &Risk{
				time:     record.Time(),
				risk:     r,
				adjucted: adjusted,
				combined: false,
			})

			if len(positions) > 0 {
				for i, p := range positions {
					r += (((p - nh) / p) * 100.0) * sizes[i]
				}

				risks = append(risks, &Risk{
					time:     record.Time(),
					risk:     r,
					adjucted: adjusted,
					combined: true,
				})
			} else {
				risks = append(risks, &Risk{
					time:     record.Time(),
					risk:     r,
					adjucted: adjusted,
					combined: true,
				})
			}

		default:
			panic(fmt.Errorf("unknown op"))
		}

		positions = append(positions, c)
		sizes = append(sizes, q)
	}

	return risks
}
