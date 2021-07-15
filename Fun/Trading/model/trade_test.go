package model

import (
	"testing"

	"github.com/KushamiNeko/GoFun/Fun/Trading/config"
)

func TestNewFuturesTradeInconsistenceSymbol(t *testing.T) {
	to1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190704",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	to2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190705",
			"symbol":    "qr",
			"operation": "+",
			"quantity":  "1",
			"price":     "10200",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190706",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10600",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190707",
			"symbol":    "qr",
			"operation": "-",
			"quantity":  "1",
			"price":     "10800",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	ts := []*FuturesTransaction{
		to1,
		to2,
		tc1,
		tc2,
	}

	_, err = NewFuturesTrade(ts)
	if err == nil {
		t.Errorf("err should not be nil")
	}
}

func TestNewFuturesTradeInconsistenceQuantity(t *testing.T) {
	to1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190704",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "2",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190706",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10600",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	ts := []*FuturesTransaction{
		to1,
		tc1,
	}

	_, err = NewFuturesTrade(ts)
	if err == nil {
		t.Errorf("err should not be nil")
	}
}

func TestNewFuturesTradeSimpleLong(t *testing.T) {
	to1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190704",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190705",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10200",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	ts := []*FuturesTransaction{
		to1,
		tc1,
	}

	trade, err := NewFuturesTrade(ts)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if trade.Symbol() != "ym" {
		t.Errorf("symbol should be ym but get %s", trade.Symbol())
	}

	if trade.Operation() != "+" {
		t.Errorf("operation should be + but get %s", trade.Operation())
	}

	if trade.OpenTime().Format(config.TimeFormatS) != "20190704" {
		t.Errorf("openTime should be 20190704 but get %s", trade.OpenTime().Format(config.TimeFormatS))
	}

	if trade.CloseTime().Format(config.TimeFormatS) != "20190705" {
		t.Errorf("closeTime should be 20190705 but get %s", trade.CloseTime().Format(config.TimeFormatS))
	}

	if trade.GL() != 200*5-1.5*2 {
		t.Errorf("GL should be 997 but get %f", trade.GL())
	}

	if trade.CommissionFees() != 1.5*2 {
		t.Errorf("CommissionFees should be 3 but get %f", trade.CommissionFees())
	}

	if trade.AvgOpenPrice() != 10000*5 {
		t.Errorf("AvgOpenPrice should be 50000 but get %f", trade.AvgOpenPrice())
	}

	if trade.AvgClosePrice() != 10200*5 {
		t.Errorf("AvgClosePrice should be 51000 but get %f", trade.AvgClosePrice())
	}

	if trade.Size() != 1 {
		t.Errorf("size should be 1 but get %d", trade.Size())
	}
}

func TestNewFuturesTradeMiddleLong(t *testing.T) {
	to1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190704",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	to2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190705",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10200",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190706",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10600",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190707",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10800",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	ts := []*FuturesTransaction{
		to1,
		to2,
		tc1,
		tc2,
	}

	trade, err := NewFuturesTrade(ts)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if trade.Symbol() != "ym" {
		t.Errorf("symbol should be ym but get %s", trade.Symbol())
	}

	if trade.Operation() != "+" {
		t.Errorf("operation should be + but get %s", trade.Operation())
	}

	if trade.OpenTime().Format(config.TimeFormatS) != "20190704" {
		t.Errorf("openTime should be 20190704 but get %s", trade.OpenTime().Format(config.TimeFormatS))
	}

	if trade.CloseTime().Format(config.TimeFormatS) != "20190707" {
		t.Errorf("closeTime should be 20190707 but get %s", trade.CloseTime().Format(config.TimeFormatS))
	}

	if trade.GL() != 600*5*2-1.5*4 {
		t.Errorf("GL should be 5994 but get %f", trade.GL())
	}

	if trade.CommissionFees() != 1.5*4 {
		t.Errorf("CommissionFees should be 6 but get %f", trade.CommissionFees())
	}

	if trade.AvgOpenPrice() != 10100*5*2 {
		t.Errorf("AvgOpenPrice should be 101000 but get %f", trade.AvgOpenPrice())
	}

	if trade.AvgClosePrice() != 10700*5*2 {
		t.Errorf("AvgClosePrice should be 107000 but get %f", trade.AvgClosePrice())
	}

	if trade.Size() != 2 {
		t.Errorf("size should be 2 but get %d", trade.Size())
	}
}

func TestNewFuturesTradeComplexLong(t *testing.T) {
	to1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190704",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	to2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190705",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10100",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190706",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "2",
			"price":     "10500",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	ts := []*FuturesTransaction{
		to1,
		to2,
		tc1,
	}

	trade, err := NewFuturesTrade(ts)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if trade.Symbol() != "ym" {
		t.Errorf("symbol should be ym but get %s", trade.Symbol())
	}

	if trade.Operation() != "+" {
		t.Errorf("operation should be + but get %s", trade.Operation())
	}

	if trade.OpenTime().Format(config.TimeFormatS) != "20190704" {
		t.Errorf("openTime should be 20190704 but get %s", trade.OpenTime().Format(config.TimeFormatS))
	}

	if trade.CloseTime().Format(config.TimeFormatS) != "20190706" {
		t.Errorf("closeTime should be 20190706 but get %s", trade.CloseTime().Format(config.TimeFormatS))
	}

	if trade.GL() != 450*5*2-1.5*4 {
		t.Errorf("GL should be 4494 but get %f", trade.GL())
	}

	if trade.CommissionFees() != 1.5*4 {
		t.Errorf("CommissionFees should be 6 but get %f", trade.CommissionFees())
	}

	if trade.AvgOpenPrice() != 10050*5*2 {
		t.Errorf("AvgOpenPrice should be 100500 but get %f", trade.AvgOpenPrice())
	}

	if trade.AvgClosePrice() != 10500*5*2 {
		t.Errorf("AvgClosePrice should be 105000 but get %f", trade.AvgClosePrice())
	}

	if trade.Size() != 2 {
		t.Errorf("size should be 2 but get %d", trade.Size())
	}
}

func TestNewFuturesTradeSimpleShort(t *testing.T) {
	to1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190704",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10200",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190705",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	ts := []*FuturesTransaction{
		to1,
		tc1,
	}

	trade, err := NewFuturesTrade(ts)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if trade.Symbol() != "ym" {
		t.Errorf("symbol should be ym but get %s", trade.Symbol())
	}

	if trade.Operation() != "-" {
		t.Errorf("operation should be - but get %s", trade.Operation())
	}

	if trade.OpenTime().Format(config.TimeFormatS) != "20190704" {
		t.Errorf("openTime should be 20190704 but get %s", trade.OpenTime().Format(config.TimeFormatS))
	}

	if trade.CloseTime().Format(config.TimeFormatS) != "20190705" {
		t.Errorf("closeTime should be 20190705 but get %s", trade.CloseTime().Format(config.TimeFormatS))
	}

	if trade.GL() != 200*5-1.5*2 {
		t.Errorf("GL should be 200 but get %f", trade.GL())
	}

	if trade.CommissionFees() != 1.5*2 {
		t.Errorf("CommissionFees should be 3 but get %f", trade.CommissionFees())
	}

	if trade.AvgOpenPrice() != 10200*5 {
		t.Errorf("AvgOpenPrice should be 51000 but get %f", trade.AvgOpenPrice())
	}

	if trade.AvgClosePrice() != 10000*5 {
		t.Errorf("AvgClosePrice should be 50000 but get %f", trade.AvgClosePrice())
	}

	if trade.Size() != 1 {
		t.Errorf("size should be 1 but get %d", trade.Size())
	}
}

func TestNewFuturesTradeMiddleShort(t *testing.T) {
	to1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190704",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10800",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	to2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190705",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "1",
			"price":     "10600",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190706",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10200",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190707",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	ts := []*FuturesTransaction{
		to1,
		to2,
		tc1,
		tc2,
	}

	trade, err := NewFuturesTrade(ts)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if trade.Symbol() != "ym" {
		t.Errorf("symbol should be ym but get %s", trade.Symbol())
	}

	if trade.Operation() != "-" {
		t.Errorf("operation should be + but get %s", trade.Operation())
	}

	if trade.OpenTime().Format(config.TimeFormatS) != "20190704" {
		t.Errorf("openTime should be 20190704 but get %s", trade.OpenTime().Format(config.TimeFormatS))
	}

	if trade.CloseTime().Format(config.TimeFormatS) != "20190707" {
		t.Errorf("closeTime should be 20190707 but get %s", trade.CloseTime().Format(config.TimeFormatS))
	}

	if trade.GL() != 600*5*2-1.5*4 {
		t.Errorf("GL should be 5994 but get %f", trade.GL())
	}

	if trade.CommissionFees() != 1.5*4 {
		t.Errorf("CommissionFees should be 6 but get %f", trade.CommissionFees())
	}

	if trade.AvgOpenPrice() != 10700*5*2 {
		t.Errorf("AvgOpenPrice should be 107000 but get %f", trade.AvgOpenPrice())
	}

	if trade.AvgClosePrice() != 10100*5*2 {
		t.Errorf("AvgClosePrice should be 101000 but get %f", trade.AvgClosePrice())
	}

	if trade.Size() != 2 {
		t.Errorf("size should be 2 but get %d", trade.Size())
	}
}

func TestNewFuturesTradeComplexShort(t *testing.T) {
	to1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190704",
			"symbol":    "ym",
			"operation": "-",
			"quantity":  "2",
			"price":     "10500",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc1, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190705",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10100",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	tc2, err := NewFuturesTransactionFromInputs(
		map[string]string{
			"time":      "20190706",
			"symbol":    "ym",
			"operation": "+",
			"quantity":  "1",
			"price":     "10000",
		},
	)

	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	ts := []*FuturesTransaction{
		to1,
		tc1,
		tc2,
	}

	trade, err := NewFuturesTrade(ts)
	if err != nil {
		t.Errorf("err should be nil but get %s", err.Error())
	}

	if trade.Symbol() != "ym" {
		t.Errorf("symbol should be ym but get %s", trade.Symbol())
	}

	if trade.Operation() != "-" {
		t.Errorf("operation should be - but get %s", trade.Operation())
	}

	if trade.OpenTime().Format(config.TimeFormatS) != "20190704" {
		t.Errorf("openTime should be 20190704 but get %s", trade.OpenTime().Format(config.TimeFormatS))
	}

	if trade.CloseTime().Format(config.TimeFormatS) != "20190706" {
		t.Errorf("closeTime should be 20190706 but get %s", trade.CloseTime().Format(config.TimeFormatS))
	}

	if trade.GL() != 450*5*2-1.5*4 {
		t.Errorf("GL should be 4494 but get %f", trade.GL())
	}

	if trade.CommissionFees() != 1.5*4 {
		t.Errorf("CommissionFees should be 6 but get %f", trade.CommissionFees())
	}

	if trade.AvgOpenPrice() != 10500*5*2 {
		t.Errorf("AvgOpenPrice should be 105000 but get %f", trade.AvgOpenPrice())
	}

	if trade.AvgClosePrice() != 10050*5*2 {
		t.Errorf("AvgClosePrice should be 100500 but get %f", trade.AvgClosePrice())
	}

	if trade.Size() != 2 {
		t.Errorf("size should be 2 but get %d", trade.Size())
	}
}
