package futures

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SymbolFormat int

//type ContractMonths string

type RollingMethod int
type AdjustingMethod int

const (
	BarchartSymbolFormat SymbolFormat = iota
	QuandlSymbolFormat

	//AllContractMonths  ContractMonths = "fghjkmnquvxz"
	//EvenContractMonths ContractMonths = "gjmqvz"

	//FinancialContractMonths ContractMonths = "hmuz"

	LastNTradingDay RollingMethod = iota
	FirstOfMonth
	OpenInterest

	PanamaCanal AdjustingMethod = iota
	Ratio

	barchartContractPattern = `(\w{2})([fghjkmnquvxz])(\d{2})`
	quandlContractPattern   = `([\d\w]+)([fghjkmnquvxz])(\d{4})`
)

func FrontContract(dtime time.Time, symbol string, contractMonths ContractMonths, format SymbolFormat) string {
	if dtime.IsZero() {
		dtime = time.Now()
	}

	yearCode := dtime.Year()
	month := dtime.Month()

	cl := len(contractMonths)

	var monthCode string

	for i := 0; i < cl; i++ {
		j := (i + 1) % cl

		if month == MonthCode(string(contractMonths[i])).Month() {
			monthCode = string(contractMonths[j])

			if i == cl-1 {
				yearCode += 1
			}

			break
		} else {
			if month < MonthCode(string(contractMonths[i])).Month() {
				monthCode = string(contractMonths[i])
				break
			}
		}
	}

	if monthCode == "" {
		panic("this should never happen")
	}

	switch format {
	case QuandlSymbolFormat:
	case BarchartSymbolFormat:
		yearCode %= 100
	default:
		panic("unknown year code")
	}

	return fmt.Sprintf("%s%s%02d", symbol, monthCode, yearCode)
}

func PreviousContract(contract string, contractMonths ContractMonths, format SymbolFormat) string {

	var pattern string
	switch format {
	case BarchartSymbolFormat:
		pattern = barchartContractPattern
	case QuandlSymbolFormat:
		pattern = quandlContractPattern
	default:
		panic("unknown symbol format")
	}

	regex := regexp.MustCompile(pattern)
	match := regex.FindAllStringSubmatch(contract, -1)

	if len(match) == 0 || len(match[0]) != 4 {
		panic("invalid contract for the specified contractFormat")
	}

	symbol := match[0][1]
	month := match[0][2]
	year := match[0][3]

	pYear, err := strconv.ParseInt(fmt.Sprintf("20%s", year), 10, 32)
	if err != nil {
		panic(err)
	}

	if int(pYear) > time.Now().Year() {
		pYear -= 100
	}

	mi := strings.Index(string(contractMonths), month) - 1
	if mi < 0 {
		pYear -= 1
		mi = len(contractMonths) - 1
	}

	pMonth := string(contractMonths[mi])

	switch format {
	case BarchartSymbolFormat:
		pYear = pYear % 100
	case QuandlSymbolFormat:
	default:
		panic("unknown symbol format")
	}

	return fmt.Sprintf("%s%s%02d", symbol, pMonth, pYear)
}

func ContractList(from, to time.Time, symbol string, contractMonths ContractMonths, format SymbolFormat) []string {

	if to.After(time.Now()) {
		to = time.Now()
	}

	last := FrontContract(to, symbol, contractMonths, format)

	first := FrontContract(from, symbol, contractMonths, format)
	first = PreviousContract(first, contractMonths, format)

	contracts := []string{
		last,
	}

	for {
		previous := PreviousContract(last, contractMonths, format)
		contracts = append(contracts, previous)
		last = previous

		if last == first {
			break
		}
	}

	return contracts
}
