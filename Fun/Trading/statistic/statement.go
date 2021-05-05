package statistic

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ErrMissingAccountBalance = errors.New("missing account balance")
var ErrMissingTradingRecords = errors.New("missing trading records")

type MonthlyStatement struct {

	// the objective of the statement
	// eg.
	// trading or investing
	objective string

	/*
		list all accounts used for specified objective

		eg.
		objective: trading
		balance:
				1 futures account for main trading
				1 equities account for holding money market ETF
				1 crypto account for earning crypto interests
	*/
	balance []*AccountBalance

	// account percentage risk per trade
	percentageRisk float64

	// trading pl records
	records []*TradeRecord
}

func (t *MonthlyStatement) Objective() string {
	return t.objective
}

func (t *MonthlyStatement) EndingBalance() float64 {
	balance := 0.0
	for _, b := range t.balance {
		balance += b.endingBalance
	}

	return balance
}

func (t *MonthlyStatement) AcountPercentageRisk() float64 {
	return t.percentageRisk
}

func (t *MonthlyStatement) TradingRecords() []*TradeRecord {
	return t.records
}

func NewMonthlyStatement(objective, statement string) *MonthlyStatement {
	t := &MonthlyStatement{
		objective: objective,
	}

	t.parseBalance(statement)
	t.parseRecords(statement)

	return t
}

func (t *MonthlyStatement) parseBalance(statement string) {

	const pattern = `(?m)^\s*(\w+)\s*balance:\s*([0-9,.-]+)\s*$`

	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllStringSubmatch(statement, -1)
	if len(matches) <= 0 {
		panic("missing balance records")
	}

	t.balance = make([]*AccountBalance, len(matches))

	for i, m := range matches {
		a := strings.TrimSpace(m[1])

		s := strings.ReplaceAll(m[2], ",", "")
		b, err := strconv.ParseFloat(s, 10)
		if err != nil {
			panic(err)
		}

		t.balance[i] = &AccountBalance{
			accountType:   a,
			endingBalance: b,
		}
	}
}

func (t *MonthlyStatement) parseRecords(statement string) error {
	const dateLayout = `2006-01-02`
	const timeLayout = `15:04`

	var err error

	var pattern string
	var regex *regexp.Regexp
	var matches [][]string

	pattern = `(?m)^\s*account percentage risk per trade\s*:\s*([0-9.]+)\s*%\s*$`
	regex = regexp.MustCompile(pattern)

	matches = regex.FindAllStringSubmatch(statement, -1)
	if len(matches) <= 0 {
		return ErrMissingTradingRecords
	}

	t.percentageRisk, err = strconv.ParseFloat(matches[0][1], 10)
	if err != nil {
		panic(err)
	}

	const datetimeInputPattern = `(\d{4}\s*[-]\s*\d{2}\s*[-]\s*\d{2}(?:[\s]\d{2}:\d{2})*)`

	pattern = fmt.Sprintf(`(?m)^\s*%s(?:\s*[-~]\s*%s)*\s*:\s*([0-9,.-]+)\s*(?:@\s*([LS])(\d*))$`, datetimeInputPattern, datetimeInputPattern)
	regex = regexp.MustCompile(pattern)

	matches = regex.FindAllStringSubmatch(statement, -1)
	if len(matches) <= 0 {
		return ErrMissingTradingRecords
	}

	trades := make([]*TradeRecord, len(matches))
	for i, m := range matches {
		var tradeOpen, tradeClose time.Time

		tradeOpen, err = time.Parse(dateLayout, m[1])
		if err != nil {
			tradeOpen, err = time.Parse(fmt.Sprintf("%s %s", dateLayout, timeLayout), m[1])
			if err != nil {
				panic(err)
			}
		}

		if m[2] != "" {
			tradeClose, err = time.Parse(dateLayout, m[2])
			if err != nil {
				tradeClose, err = time.Parse(fmt.Sprintf("%s %s", dateLayout, timeLayout), m[2])
				if err != nil {
					panic(err)
				}
			}
		} else {
			tradeClose = tradeOpen
		}

		var direction TradeDirection
		switch m[4] {
		case "L":
			direction = LongTrade
		case "S":
			direction = ShortTrade

		default:
			panic("unkonw trade direction")
		}

		s := strings.ReplaceAll(m[3], ",", "")
		pl, err := strconv.ParseFloat(s, 10)
		if err != nil {
			panic(err)
		}

		trades[i] = NewTradeRecord(tradeOpen, tradeClose, direction, pl)

	}

	t.records = trades

	return nil
}
