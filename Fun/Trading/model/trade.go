package model

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/KushamiNeko/GoFun/Fun/Trading/config"
)

const (
	commissionFees = 1.5
)

type FuturesTrade struct {
	transactions []*FuturesTransaction
	size         int

	o []*FuturesTransaction
	c []*FuturesTransaction
}

func NewFuturesTrade(transactions []*FuturesTransaction) (*FuturesTrade, error) {
	if len(transactions) == 0 {
		return nil, fmt.Errorf("invalid transactions")
	}

	f := new(FuturesTrade)

	sort.Slice(transactions, func(i, j int) bool {
		if transactions[i].Time().Equal(transactions[j].Time()) {
			return transactions[i].TimeStamp() < transactions[j].TimeStamp()
		} else {
			return transactions[i].Time().Before(transactions[j].Time())
		}
	})

	f.transactions = transactions
	err := f.processing()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *FuturesTrade) processing() error {
	o := make([]*FuturesTransaction, 0)
	c := make([]*FuturesTransaction, 0)

	oq := 0
	cq := 0

	ta := 0

	for _, t := range f.transactions {
		if t.Symbol() != f.Symbol() {
			return fmt.Errorf("inconsistence symbols: %s, %s", t.Symbol(), f.Symbol())
		}

		op := t.Operation()

		if op == f.Operation() {
			oq = oq + t.Quantity()
			o = append(o, t)
		} else {
			cq = cq + t.Quantity()
			c = append(c, t)
		}

		ta += t.Action()
		switch f.Operation() {
		case "+":
			f.size = int(math.Max(float64(f.size), float64(ta)))
		case "-":
			f.size = int(math.Min(float64(f.size), float64(ta)))
		default:
			panic(fmt.Sprintf("unknonw operation: %s", f.Operation()))
		}
	}

	if oq != cq {
		return fmt.Errorf("inconsistence quantity: %d, %d", oq, cq)
	}

	f.o = o
	f.c = c

	f.size = int(math.Abs(float64(f.size)))

	return nil
}

func (f *FuturesTrade) averagePrice(transactions []*FuturesTransaction) float64 {
	q := 0
	var tp float64 = 0.0

	for _, t := range transactions {
		q += t.Quantity()
		tp += t.Price() * float64(t.Quantity())
	}

	return tp / float64(q)
}

func (f *FuturesTrade) Transactions() []*FuturesTransaction {
	return f.transactions
}

func (f *FuturesTrade) OpenOrders() []*FuturesTransaction {
	return f.o
}

func (f *FuturesTrade) CloseOrders() []*FuturesTransaction {
	return f.c
}

func (f *FuturesTrade) Operation() string {
	return f.transactions[0].Operation()
}

func (f *FuturesTrade) Symbol() string {
	return f.transactions[0].Symbol()
}

func (f *FuturesTrade) OpenTime() time.Time {
	return f.transactions[0].Time()
}

func (f *FuturesTrade) CloseTime() time.Time {
	return f.transactions[len(f.transactions)-1].Time()
}

func (f *FuturesTrade) OpenTimeStamp() int64 {
	return f.transactions[0].TimeStamp()
}

func (f *FuturesTrade) CloseTimeStamp() int64 {
	return f.transactions[len(f.transactions)-1].TimeStamp()
}

func (f *FuturesTrade) Size() int {
	return f.size
}

func (f *FuturesTrade) CommissionFees() float64 {
	return commissionFees * float64(f.Size()) * 2
}

func (f *FuturesTrade) AvgOpenPrice() float64 {
	c := NewContractSpecs()
	unit, _ := c.LookupContractUnit(f.Symbol())
	return f.averagePrice(f.o) * float64(f.Size()) * unit
}

func (f *FuturesTrade) AvgClosePrice() float64 {
	c := NewContractSpecs()
	unit, _ := c.LookupContractUnit(f.Symbol())
	return f.averagePrice(f.c) * float64(f.Size()) * unit
}

func (f *FuturesTrade) GL() float64 {

	var o float64
	var c float64

	if f.Operation() == "+" {
		o = -1 * f.AvgOpenPrice()
		c = f.AvgClosePrice()
	} else {
		o = f.AvgOpenPrice()
		c = -1 * f.AvgClosePrice()
	}

	return o + c - f.CommissionFees()
}

func (f *FuturesTrade) GLP() float64 {
	glp := (f.GL() / f.AvgOpenPrice()) * 100.0
	return glp
}

func (f *FuturesTrade) Entity() map[string]string {
	return map[string]string{
		"operation":       f.Operation(),
		"symbol":          f.Symbol(),
		"open_time":       f.transactions[0].time,
		"close_time":      f.transactions[len(f.transactions)-1].time,
		"average_open":    fmt.Sprintf("%.2f", f.AvgOpenPrice()),
		"average_close":   fmt.Sprintf("%.2f", f.AvgClosePrice()),
		"gl":              fmt.Sprintf("%.2f", f.GL()),
		"glp":             fmt.Sprintf("%.2f", f.GLP()),
		"commission_fees": fmt.Sprintf("%.2f", f.CommissionFees()),
		"size":            fmt.Sprintf("%d", f.Size()),
	}
}

const (
	futuresTradeFmtString = "%-[1]*[4]s%-[2]*[5]s%-[1]*[6]s%-[2]*[7]s%-[2]*[8]s%-[3]*[9]s%-[3]*[10]s%-[3]*[11]s%-[3]*[12]s%-[3]*[13]s"
)

func (f *FuturesTrade) Fmt() string {
	return fmt.Sprintf(
		futuresTradeFmtString,
		config.FmtWidth,
		config.FmtWidthL,
		config.FmtWidthXL,
		f.Symbol(),
		f.Operation(),
		fmt.Sprintf("%d", f.Size()),
		f.transactions[0].time,
		f.transactions[len(f.transactions)-1].time,
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.AvgOpenPrice()),
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.AvgClosePrice()),
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.CommissionFees()),
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.GL()),
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.GLP()),
	)
}

func FuturesTradeFmtLabels() string {
	return fmt.Sprintf(
		futuresTradeFmtString,
		config.FmtWidth,
		config.FmtWidthL,
		config.FmtWidthXL,
		"Symbol",
		"Operation",
		"Size",
		"Open Date",
		"Close Date",
		"Avg Open Price",
		"Avg Close Price",
		"Commission Fees",
		"GL($)",
		"GL(%)",
	)
}
