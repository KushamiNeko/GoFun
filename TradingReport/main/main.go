package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/KushamiNeko/GoFun/Trading/statistic"
)

// var ErrMissingAccountBalance = errors.New("missing account balance")
// var ErrMissingTradingRecords = errors.New("missing trading records")

// type MonthlyTradingStatement struct {
// 	// start time.Time
// 	// end   time.Time

// 	// account percentage risk per trade
// 	percentageRisk float64

// 	records []*TradeRecord
// }

// func (t *MonthlyTradingStatement) Parse(record string) error {
// 	const dateLayout = `2006-01-02`
// 	const timeLayout = `15:04`

// 	var err error

// 	var pattern string
// 	var regex *regexp.Regexp
// 	var matches [][]string

// 	pattern = `(?m)^\s*account percentage risk per trade\s*:\s*([0-9.]+)\s*%\s*$`
// 	regex = regexp.MustCompile(pattern)

// 	matches = regex.FindAllStringSubmatch(record, -1)
// 	if len(matches) <= 0 {
// 		return ErrMissingTradingRecords
// 	}

// 	t.percentageRisk, err = strconv.ParseFloat(matches[0][1], 10)
// 	if err != nil {
// 		panic(err)
// 	}

// 	const datetimeInputPattern = `(\d{4}\s*[-]\s*\d{2}\s*[-]\s*\d{2}(?:[\s]\d{2}:\d{2})*)`

// 	pattern = fmt.Sprintf(`(?m)^\s*%s(?:\s*[-~]\s*%s)*\s*:\s*([0-9,.-]+)\s*(?:@\s*([LS])(\d*))$`, datetimeInputPattern, datetimeInputPattern)
// 	regex = regexp.MustCompile(pattern)

// 	matches = regex.FindAllStringSubmatch(record, -1)
// 	if len(matches) <= 0 {
// 		return ErrMissingTradingRecords
// 	}

// 	trades := make([]*TradeRecord, len(matches))
// 	for i, m := range matches {
// 		var tradeOpen, tradeClose time.Time

// 		tradeOpen, err = time.Parse(dateLayout, m[1])
// 		if err != nil {
// 			tradeOpen, err = time.Parse(fmt.Sprintf("%s %s", dateLayout, timeLayout), m[1])
// 			if err != nil {
// 				panic(err)
// 			}
// 		}

// 		if m[2] != "" {
// 			tradeClose, err = time.Parse(dateLayout, m[2])
// 			if err != nil {
// 				tradeClose, err = time.Parse(fmt.Sprintf("%s %s", dateLayout, timeLayout), m[2])
// 				if err != nil {
// 					panic(err)
// 				}
// 			}
// 		} else {
// 			tradeClose = tradeOpen
// 		}

// 		var direction TradeDirection
// 		switch m[4] {
// 		case "L":
// 			direction = LongTrade
// 		case "S":
// 			direction = ShortTrade

// 		default:
// 			panic("unkonw trade direction")
// 		}

// 		// var s string
// 		// s = strings.ReplaceAll(m[5], ",", "")
// 		// size, err := strconv.ParseFloat(s, 10)
// 		// if err != nil {
// 		// 	panic(err)
// 		// }

// 		s := strings.ReplaceAll(m[3], ",", "")
// 		pl, err := strconv.ParseFloat(s, 10)
// 		if err != nil {
// 			panic(err)
// 		}

// 		trades[i] = NewTradeRecord(tradeOpen, tradeClose, direction, pl)

// 		fmt.Println(trades[i])
// 	}

// 	t.records = trades

// 	return nil
// }

// /* func (t *TradingStatistic) calculate() {
// 	if len(t.records) <= 0 {
// 		panic("invalid trading records")
// 	}

// 	var count, countL, countS float64
// 	var winCount, winCountL, winCountS float64
// 	var lossCount, lossCountL, lossCountS float64

// 	var winAvgPL, lossAvgPL float64
// 	// var winAvgHold, lossAvgHold float64

// 	for _, r := range t.records {
// 		if !t.start.IsZero() && !t.end.IsZero() {
// 			if r.TradeClose().Before(t.start) || r.TradeClose().After(t.end) {
// 				continue
// 			}
// 		}

// 		count++

// 		switch {
// 		case r.Direction() == LongTrade:
// 		case r.Direction() == ShortTrade:
// 		case r.PL() > 0:
// 			winCount++
// 			winAvgPL += r.PL()

// 		case r.PL() < 0:
// 			lossCount++
// 			lossAvgPL += r.PL()

// 		default:
// 			continue
// 		}
// 	}

// 	t.battingAvg = winCount / count
// } */

// type MonthlyAccountBalance struct {
// 	accountType string

// 	// startingBalance float64
// 	endingBalance float64
// }

// type TradeDirection int

// func (t TradeDirection) String() string {
// 	switch t {
// 	case LongTrade:
// 		return "Long"

// 	case ShortTrade:
// 		return "Short"

// 	default:
// 		panic("unknown trade direction")
// 	}
// }

// const (
// 	LongTrade TradeDirection = iota
// 	ShortTrade
// )

// type TradeRecord struct {
// 	tradeOpen  time.Time
// 	tradeClose time.Time

// 	direction TradeDirection

// 	// size float64

// 	pl float64
// }

// func NewTradeRecord(tradeOpen, tradeClose time.Time, direction TradeDirection, pl float64) *TradeRecord {
// 	if tradeOpen.IsZero() || tradeClose.IsZero() {
// 		panic("invalid time for opening or closing trade")
// 	}

// 	return &TradeRecord{
// 		tradeOpen:  tradeOpen,
// 		tradeClose: tradeClose,
// 		direction:  direction,
// 		// size:       size,
// 		pl: pl,
// 	}
// }

// func (t *TradeRecord) TradeOpen() time.Time {
// 	return t.tradeOpen
// }

// func (t *TradeRecord) TradeClose() time.Time {
// 	return t.tradeClose
// }

// func (t *TradeRecord) Direction() TradeDirection {
// 	return t.direction
// }

// // func (t *TradeRecord) Size() float64 {
// // 	return t.size
// // }

// func (t *TradeRecord) PL() float64 {
// 	return t.pl
// }

// func parseBalance(record string) {
// 	const pattern = `(?m)^\s*(\w+)\s*balance:\s*([0-9,.-]+)\s*$`

// 	regex := regexp.MustCompile(pattern)

// 	matches := regex.FindAllStringSubmatch(record, -1)
// 	if len(matches) <= 0 {
// 		panic("missing balance records")
// 	}

// 	var builder strings.Builder
// 	// builder.WriteString("Account Balance:\n")
// 	fmt.Fprintf(&builder, "Account Balance: %d\n", len(matches))

// 	balance := 0.0
// 	for _, m := range matches {
// 		a := strings.TrimSpace(m[1])
// 		fmt.Fprintf(&builder, "%s\n", a)

// 		s := strings.ReplaceAll(m[2], ",", "")
// 		b, err := strconv.ParseFloat(s, 10)
// 		if err != nil {
// 			panic(err)
// 		}

// 		balance += b
// 	}

// 	fmt.Println(builder.String())
// 	fmt.Println(balance)
// }

// func parseTrading(record string) {

// }

func main() {
	root := filepath.Join(os.Getenv("HOME"), "Documents", "TRADING_NOTES", "records", "KushamiNeko", "trading")

	fs, err := os.ReadDir(root)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		bs, err := os.ReadFile(filepath.Join(root, f.Name()))
		if err != nil {
			panic(err)
		}

		content := bytes.NewBuffer(bs)
		statements := statistic.NewMonthlyStatement("trading", content.String())
		stat := statistic.NewTradingStatistic(time.Time{}, time.Time{}, []*statistic.MonthlyStatement{statements})
		fmt.Println(stat)

		// parseBalance(content.String())

		// t := &MonthlyTradingStatement{}
		// t.Parse(content.String())
		// parseTrading(content.String())
	}
}
