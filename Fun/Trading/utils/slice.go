package utils

import (
	"time"

	"github.com/KushamiNeko/GoFun/Trading/model"
)

type timing struct {
	times []time.Time
	flip  time.Time
	ct    time.Time
	//ct *model.FuturesTransaction
}

func (t timing) Times() []time.Time {
	return t.times
}

func (t timing) From() time.Time {
	return t.times[0]
}

func (t timing) To() time.Time {
	return t.times[len(t.times)-1]
}

func (t timing) Flip() time.Time {
	return t.flip
}

func (t timing) CloseTime() time.Time {
	return t.ct
}

//func (t timing) CloseRecord() *model.FuturesTransaction {
//return t.ct
//}

func TradeEntry(records []*model.FuturesTransaction, op string) *timing {
	rt := make([]time.Time, 0, len(records))
	var f time.Time

	for _, t := range records {
		if t.Operation() == op {
			rt = append(rt, t.Time())
		} else {
			f = t.Time()
			break
		}
	}

	return &timing{
		times: rt,
		flip:  f,
		//ct:    records[len(records)-1],
		ct: records[len(records)-1].Time(),
	}
}

func TradeExit(records []*model.FuturesTransaction, op string) *timing {
	rt := make([]time.Time, 0, len(records))
	var f time.Time

	for i := len(records) - 1; i >= 0; i-- {
		t := records[i]
		if t.Operation() != op {
			rt = append(rt, t.Time())
		} else {
			f = t.Time()
			break
		}
	}

	return &timing{
		times: rt,
		flip:  f,
		//ct:    records[len(records)-1],
		ct: records[len(records)-1].Time(),
	}
}

func TradeSameOperation(records []*model.FuturesTransaction, op string) *timing {
	rt := make([]time.Time, 0, len(records))

	action := 0

	for _, t := range records {
		if t.Operation() == op {
			rt = append(rt, t.Time())
		}

		action += t.Action()
	}

	flip := records[len(records)-1].Time()
	if action != 0 {
		flip = flip.Add(30 * 24 * time.Hour)
	} else {
		for i := len(records) - 1; i > 0; i-- {
			if records[i].Operation() != op {
				flip = records[i].Time()
			} else {
				break
			}
		}
	}

	return &timing{
		times: rt,
		flip:  flip,
		//ct:    records[len(records)-1],
		ct: records[len(records)-1].Time(),
	}
}
