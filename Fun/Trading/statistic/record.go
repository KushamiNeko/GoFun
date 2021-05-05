package statistic

import "time"

const (
	LongTrade TradeDirection = iota
	ShortTrade
)

type TradeDirection int

func (t TradeDirection) String() string {
	switch t {
	case LongTrade:
		return "Long"

	case ShortTrade:
		return "Short"

	default:
		panic("unknown trade direction")
	}
}

type TradeRecord struct {
	tradeOpen  time.Time
	tradeClose time.Time

	direction TradeDirection

	pl float64
}

func NewTradeRecord(tradeOpen, tradeClose time.Time, direction TradeDirection, pl float64) *TradeRecord {
	if tradeOpen.IsZero() || tradeClose.IsZero() {
		panic("invalid time for opening or closing trade")
	}

	return &TradeRecord{
		tradeOpen:  tradeOpen,
		tradeClose: tradeClose,
		direction:  direction,
		pl:         pl,
	}
}

func (t *TradeRecord) TradeOpen() time.Time {
	return t.tradeOpen
}

func (t *TradeRecord) TradeClose() time.Time {
	return t.tradeClose
}

func (t *TradeRecord) Direction() TradeDirection {
	return t.direction
}

func (t *TradeRecord) PL() float64 {
	return t.pl
}
