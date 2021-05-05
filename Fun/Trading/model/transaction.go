package model

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/KushamiNeko/GoFun/Trading/config"
	"github.com/KushamiNeko/GoFun/Utility/hashutils"
)

type FuturesTransaction struct {
	index     string
	timeStamp string

	time      string
	symbol    string
	operation string
	quantity  string
	price     string
	note      string
}

func NewFuturesTransaction(
	t string,
	s string,
	o string,
	q string,
	p string,
	n string,
	index string,
	timeStamp string,
) (*FuturesTransaction, error) {

	f := new(FuturesTransaction)
	f.time = t
	f.symbol = s
	f.operation = o
	f.quantity = q
	f.price = p
	f.note = n

	if index == "" {
		f.index = hashutils.RandString(config.IDLen)
	} else {
		f.index = index
	}

	if timeStamp == "" {
		f.timeStamp = strconv.FormatInt(time.Now().UnixNano(), 10)
	} else {
		f.timeStamp = timeStamp
	}

	err := f.validateInput()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func NewFuturesTransactionFromInputs(entity map[string]string) (*FuturesTransaction, error) {

	time, ok := entity["time"]
	if !ok {
		return nil, fmt.Errorf("missing time")
	}

	symbol, ok := entity["symbol"]
	if !ok {
		return nil, fmt.Errorf("missing symbol")
	}

	operation, ok := entity["operation"]
	if !ok {
		return nil, fmt.Errorf("missing operation")
	}

	quantity, ok := entity["quantity"]
	if !ok {
		return nil, fmt.Errorf("missing quantity")
	}

	price, ok := entity["price"]
	if !ok {
		return nil, fmt.Errorf("missing price")
	}

	note := entity["note"]

	f, err := NewFuturesTransaction(time, symbol, operation, quantity, price, note, "", "")
	if err != nil {
		return nil, err
	}

	return f, nil
}

func NewFuturesTransactionFromEntity(entity map[string]string) (*FuturesTransaction, error) {

	time, ok := entity["time"]
	if !ok {
		return nil, fmt.Errorf("missing time")
	}

	symbol, ok := entity["symbol"]
	if !ok {
		return nil, fmt.Errorf("missing symbol")
	}

	operation, ok := entity["operation"]
	if !ok {
		return nil, fmt.Errorf("missing operation")
	}

	quantity, ok := entity["quantity"]
	if !ok {
		return nil, fmt.Errorf("missing quantity")
	}

	price, ok := entity["price"]
	if !ok {
		return nil, fmt.Errorf("missing price")
	}

	note := entity["note"]

	index, ok := entity["index"]
	if !ok {
		return nil, fmt.Errorf("missing index")
	}

	timeStamp, ok := entity["time_stamp"]
	if !ok {
		return nil, fmt.Errorf("missing timeStamp")
	}

	f, err := NewFuturesTransaction(time, symbol, operation, quantity, price, note, index, timeStamp)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *FuturesTransaction) validateInput() error {

	const (
		reTime      = `^\d{8}(?:\s*\d{2}:\d{2}:\d{2})*$`
		reSymbol    = `^[a-z]+$`
		reOperation = `^(?:\+|\-)$`
		reQuantity  = `^[0-9]+$`
		rePrice     = `^[0-9.]+$`
		reNote      = `^[^;]*$`
		reIndex     = `^[a-zA-Z0-9]+$`
		reTimeStamp = `^[0-9.]+$`
	)

	var re *regexp.Regexp

	re = regexp.MustCompile(reTime)
	if !re.MatchString(f.time) {
		return fmt.Errorf("invalid time: %s", f.time)
	}

	re = regexp.MustCompile(reSymbol)
	if !re.MatchString(f.symbol) {
		return fmt.Errorf("invalid symbol: %s", f.symbol)
	}

	c := NewContractSpecs()
	if !c.ValidateSymbol(f.symbol) {
		return fmt.Errorf("invalid symbol")
	}

	re = regexp.MustCompile(reOperation)
	if !re.MatchString(f.operation) {
		return fmt.Errorf("invalid operation: %s", f.operation)
	}

	re = regexp.MustCompile(reQuantity)
	if !re.MatchString(f.quantity) {
		return fmt.Errorf("invalid quantity: %s", f.quantity)
	}

	re = regexp.MustCompile(rePrice)
	if !re.MatchString(f.price) {
		return fmt.Errorf("invalid price: %s", f.price)
	}

	re = regexp.MustCompile(reNote)
	if !re.MatchString(f.note) {
		return fmt.Errorf("invalid note: %s", f.note)
	}

	re = regexp.MustCompile(reIndex)
	if !re.MatchString(f.index) {
		return fmt.Errorf("invalid index: %s", f.index)
	}

	re = regexp.MustCompile(reTimeStamp)
	if !re.MatchString(f.timeStamp) {
		return fmt.Errorf("invalid timeStamp: %s", f.timeStamp)
	}

	return nil
}

func (f *FuturesTransaction) Index() string {
	return f.index
}

func (f *FuturesTransaction) TimeStamp() int64 {
	t, err := strconv.ParseInt(f.timeStamp, 10, 64)
	if err != nil {
		panic(err)
	}

	return t
}

func (f *FuturesTransaction) Time() time.Time {
	var d time.Time
	m, err := regexp.MatchString(`^\d{8}$`, f.time)
	if err != nil {
		panic(err)
	}

	if m {
		d, err = time.Parse(config.TimeFormatS, f.time)
		if err != nil {
			panic(err)
		}
	} else {
		d, err = time.Parse(config.TimeFormatL, f.time)
		if err != nil {
			panic(err)
		}
	}

	return d
}

func (f *FuturesTransaction) Symbol() string {
	return f.symbol
}

func (f *FuturesTransaction) Operation() string {
	return f.operation
}

func (f *FuturesTransaction) Quantity() int {
	q, err := strconv.ParseInt(f.quantity, 10, 64)
	if err != nil {
		panic(err)
	}

	return int(q)
}

func (f *FuturesTransaction) Price() float64 {
	p, err := strconv.ParseFloat(f.price, 64)
	if err != nil {
		panic(err)
	}

	return p
}

func (f *FuturesTransaction) Note() string {
	return f.note
}

func (f *FuturesTransaction) TotalPrice() float64 {
	return f.Price() * float64(f.Quantity())
}

func (f *FuturesTransaction) Action() int {
	switch f.operation {
	case "+":
		return f.Quantity()
	case "-":
		return f.Quantity() * -1
	default:
		panic(fmt.Sprintf("invalid operation: %s", f.operation))
	}
}

func (f *FuturesTransaction) Entity() map[string]string {
	return map[string]string{
		"index":      f.index,
		"time_stamp": f.timeStamp,
		"time":       f.time,
		"symbol":     f.symbol,
		"operation":  f.operation,
		"quantity":   f.quantity,
		"price":      f.price,
		"note":       f.note,
	}
}

const (
	futuresTransactionFmtString = "%-[2]*[4]s%-[1]*[5]s%-[2]*[6]s%-[2]*[7]s%-[2]*[8]s%[9]s"
)

func (f *FuturesTransaction) Fmt() string {
	return fmt.Sprintf(
		futuresTransactionFmtString,
		config.FmtWidth,
		config.FmtWidthL,
		config.FmtWidthXL,
		f.time,
		f.symbol,
		f.operation,
		f.quantity,
		fmt.Sprintf("%.[1]*f", config.DollarDecimals, f.Price()),
		f.note,
	)
}

func FuturesTransactionFmtLabels() string {
	return fmt.Sprintf(
		futuresTransactionFmtString,
		config.FmtWidth,
		config.FmtWidthL,
		config.FmtWidthXL,
		"Date",
		"Symbol",
		"Operation",
		"Quantity",
		"Price",
		"Note",
	)
}
