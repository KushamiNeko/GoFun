package statistic

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func plRange(records []*TradeRecord) (min float64, max float64) {
	max = float64(math.MinInt32)
	min = float64(math.MaxInt32)

	for _, r := range records {
		if r.PL() > max {
			max = r.PL()
		}

		if r.PL() < min {
			min = r.PL()
		}
	}

	return min, max
}

func holdingRangeHour(records []*TradeRecord) (min float64, max float64) {
	max = float64(math.MinInt32)
	min = float64(math.MaxInt32)

	for _, r := range records {
		diff := holdingLengthHour(r)

		if diff > max {
			max = diff
		}

		if diff < min {
			min = diff
		}
	}

	return min, max
}

func plAvg(records []*TradeRecord) float64 {
	total := 0.0

	for _, r := range records {
		total += r.PL()
	}

	return total / float64(len(records))
}

func holdingAvgHour(records []*TradeRecord) float64 {
	total := 0.0

	for _, r := range records {
		total += holdingLengthHour(r)
	}

	return total / float64(len(records))
}

func holdingLengthHour(r *TradeRecord) float64 {
	to := r.TradeOpen()
	tc := r.TradeClose()

	diffM := tc.Sub(to).Minutes()

	return diffM / 60.0
}

type TradingStatistic struct {
	start time.Time
	end   time.Time

	statements []*MonthlyStatement

	all     []*TradeRecord
	winners []*TradeRecord
	losers  []*TradeRecord

	long  []*TradeRecord
	short []*TradeRecord

	longWinners []*TradeRecord
	longLosers  []*TradeRecord

	shortWinners []*TradeRecord
	shortLosers  []*TradeRecord

	battingAvg  float64
	battingAvgL float64
	battingAvgS float64

	winAvg   float64
	winAvgL  float64
	winAvgS  float64
	winMax   float64
	winMaxL  float64
	winMaxS  float64
	winMin   float64
	winMinL  float64
	winMinS  float64
	lossAvg  float64
	lossAvgL float64
	lossAvgS float64
	lossMax  float64
	lossMaxL float64
	lossMaxS float64
	lossMin  float64
	lossMinL float64
	lossMinS float64

	winLossRatio  float64
	winLossRatioL float64
	winLossRatioS float64

	winHoldAvg   float64
	winHoldAvgL  float64
	winHoldAvgS  float64
	winHoldMax   float64
	winHoldMaxL  float64
	winHoldMaxS  float64
	winHoldMin   float64
	winHoldMinL  float64
	winHoldMinS  float64
	lossHoldAvg  float64
	lossHoldAvgL float64
	lossHoldAvgS float64
	lossHoldMax  float64
	lossHoldMaxL float64
	lossHoldMaxS float64
	lossHoldMin  float64
	lossHoldMinL float64
	lossHoldMinS float64

	expectedValue  float64
	expectedValueL float64
	expectedValueS float64
}

func NewTradingStatistic(start, end time.Time, statements []*MonthlyStatement) *TradingStatistic {
	t := new(TradingStatistic)

	t.start = start
	t.end = end

	t.statements = statements

	t.all = make([]*TradeRecord, 0)
	t.winners = make([]*TradeRecord, 0)
	t.losers = make([]*TradeRecord, 0)
	t.long = make([]*TradeRecord, 0)
	t.short = make([]*TradeRecord, 0)
	t.longWinners = make([]*TradeRecord, 0)
	t.longLosers = make([]*TradeRecord, 0)
	t.shortWinners = make([]*TradeRecord, 0)
	t.shortLosers = make([]*TradeRecord, 0)

	for _, s := range statements {
		for _, r := range s.TradingRecords() {
			if !t.start.IsZero() && !t.end.IsZero() {
				if r.TradeClose().Before(t.start) || r.TradeClose().After(t.end) {
					continue
				}
			}

			t.all = append(t.all, r)

			if r.PL() > 0 {
				t.winners = append(t.winners, r)
			}

			if r.PL() < 0 {
				t.losers = append(t.losers, r)
			}

			if r.Direction() == LongTrade {
				t.long = append(t.long, r)
			}

			if r.Direction() == ShortTrade {
				t.short = append(t.short, r)
			}

			if r.PL() > 0 && r.Direction() == LongTrade {
				t.longWinners = append(t.longWinners, r)
			}

			if r.PL() < 0 && r.Direction() == LongTrade {
				t.longLosers = append(t.longLosers, r)
			}

			if r.PL() > 0 && r.Direction() == ShortTrade {
				t.shortWinners = append(t.shortWinners, r)
			}

			if r.PL() < 0 && r.Direction() == ShortTrade {
				t.shortLosers = append(t.shortLosers, r)
			}
		}
	}

	return t
}

func (t *TradingStatistic) calculate() {

	t.battingAvg = float64(len(t.winners)) / float64(len(t.all))
	t.battingAvgL = float64(len(t.longWinners)) / float64(len(t.long))
	t.battingAvgS = float64(len(t.shortWinners)) / float64(len(t.short))

	t.winAvg = plAvg(t.winners)
	t.winAvgL = plAvg(t.longWinners)
	t.winAvgS = plAvg(t.shortWinners)

	t.lossAvg = plAvg(t.losers)
	t.lossAvgL = plAvg(t.longLosers)
	t.lossAvgS = plAvg(t.shortLosers)

	t.winMin, t.winMax = plRange(t.winners)
	t.winMinL, t.winMaxL = plRange(t.longWinners)
	t.winMinS, t.winMaxS = plRange(t.shortWinners)

	t.lossMin, t.lossMax = plRange(t.losers)
	t.lossMinL, t.lossMaxL = plRange(t.longLosers)
	t.lossMinS, t.lossMaxS = plRange(t.shortLosers)

	t.winLossRatio = t.winAvg / math.Abs(t.lossAvg)
	t.winLossRatioL = t.winAvgL / math.Abs(t.lossAvgL)
	t.winLossRatioS = t.winAvgS / math.Abs(t.lossAvgS)

	t.winHoldAvg = holdingAvgHour(t.winners)
	t.winHoldAvgL = holdingAvgHour(t.longWinners)
	t.winHoldAvgS = holdingAvgHour(t.shortWinners)

	t.lossHoldAvg = holdingAvgHour(t.losers)
	t.lossHoldAvgL = holdingAvgHour(t.longLosers)
	t.lossHoldAvgS = holdingAvgHour(t.shortLosers)

	t.winHoldMin, t.winHoldMax = holdingRangeHour(t.winners)
	t.winHoldMinL, t.winHoldMaxL = holdingRangeHour(t.longWinners)
	t.winHoldMinS, t.winHoldMaxS = holdingRangeHour(t.shortWinners)

	t.lossHoldMin, t.lossHoldMax = holdingRangeHour(t.losers)
	t.lossHoldMinL, t.lossHoldMaxL = holdingRangeHour(t.longLosers)
	t.lossHoldMinS, t.lossHoldMaxS = holdingRangeHour(t.shortLosers)

	t.expectedValue = (t.winAvg * t.battingAvg) + (t.lossAvg * (1.0 - t.battingAvg))
	t.expectedValueL = (t.winAvgL * t.battingAvgL) + (t.lossAvgL * (1.0 - t.battingAvgL))
	t.expectedValueS = (t.winAvgS * t.battingAvgS) + (t.lossAvgS * (1.0 - t.battingAvgS))
}

func (t *TradingStatistic) String() string {
	t.calculate()

	const fieldLayout = `%- 18s`

	var builder strings.Builder

	builder.WriteString(fieldLayout)
	builder.WriteString(" | ")
	builder.WriteString(fieldLayout)
	builder.WriteString(" | ")
	builder.WriteString(fieldLayout)
	builder.WriteString(" | ")
	builder.WriteString(fieldLayout)
	builder.WriteString(" | ")
	builder.WriteString(fieldLayout)
	builder.WriteString(" | ")
	builder.WriteString(fieldLayout)
	builder.WriteString(" | ")
	builder.WriteString(fieldLayout)
	builder.WriteString(" | ")
	builder.WriteString(fieldLayout)
	builder.WriteString("\n")

	row := builder.String()

	separator := fmt.Sprintf(
		row,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	)

	separator = strings.ReplaceAll(separator, " ", "-")

	// regex := regexp.MustCompile(`-\|-`)
	// separator = regex.ReplaceAllString(separator, " | ")

	builder = strings.Builder{}

	fmt.Fprintf(
		&builder,
		row,
		"",
		"Batting Avg. (%)",
		"Win PL Avg. ($|¥)",
		"Loss PL Avg. ($|¥)",
		"Win Loss Ratio",
		"Win Hold Avg. (H)",
		"Loss Hold Avg. (H)",
		"EV. ($|¥)",
	)

	builder.WriteString(separator)

	fmt.Fprintf(
		&builder,
		row,
		"ALL",
		fmt.Sprintf("%.3f", t.battingAvg*100.0),
		fmt.Sprintf("%.3f", t.winAvg),
		fmt.Sprintf("%.3f", t.lossAvg),
		fmt.Sprintf("%.3f", t.winLossRatio),
		fmt.Sprintf("%.3f", t.winHoldAvg),
		fmt.Sprintf("%.3f", t.lossHoldAvg),
		fmt.Sprintf("%.3f", t.expectedValue),
	)

	builder.WriteString(separator)

	fmt.Fprintf(
		&builder,
		row,
		"LONG",
		fmt.Sprintf("%.3f", t.battingAvgL*100.0),
		fmt.Sprintf("%.3f", t.winAvgL),
		fmt.Sprintf("%.3f", t.lossAvgL),
		fmt.Sprintf("%.3f", t.winLossRatioL),
		fmt.Sprintf("%.3f", t.winHoldAvgL),
		fmt.Sprintf("%.3f", t.lossHoldAvgL),
		fmt.Sprintf("%.3f", t.expectedValueL),
	)

	builder.WriteString(separator)

	fmt.Fprintf(
		&builder,
		row,
		"SHORT",
		fmt.Sprintf("%.3f", t.battingAvgS*100.0),
		fmt.Sprintf("%.3f", t.winAvgS),
		fmt.Sprintf("%.3f", t.lossAvgS),
		fmt.Sprintf("%.3f", t.winLossRatioS),
		fmt.Sprintf("%.3f", t.winHoldAvgS),
		fmt.Sprintf("%.3f", t.lossHoldAvgS),
		fmt.Sprintf("%.3f", t.expectedValueS),
	)

	return builder.String()
}
