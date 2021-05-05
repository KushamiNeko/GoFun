package futures

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ContractMonths string
type MonthCode string

const (
	AllContractMonths  ContractMonths = "fghjkmnquvxz"
	EvenContractMonths ContractMonths = "gjmqvz"

	OddContractMonths       ContractMonths = "fhknux"
	FinancialContractMonths ContractMonths = "hmuz"

	CryptoContractMonths ContractMonths = "fjnv"
)

const (
	January   MonthCode = "f"
	February  MonthCode = "g"
	March     MonthCode = "h"
	April     MonthCode = "j"
	May       MonthCode = "k"
	June      MonthCode = "m"
	July      MonthCode = "n"
	August    MonthCode = "q"
	September MonthCode = "u"
	October   MonthCode = "v"
	November  MonthCode = "x"
	December  MonthCode = "z"
)

func DefaultContractMonths(symbol string) ContractMonths {
	switch symbol {
	case "cl":
		return AllContractMonths
	case "ng":
		return AllContractMonths
	case "gc":
		return EvenContractMonths
	case "si":
		return "hknuz"
	case "hg":
		return "hknuz"
	case "zs":
		return "fhknqux"
	case "zc":
		return "hknuz"
	case "zw":
		return "hknuz"
	default:
		return FinancialContractMonths
	}
}

func NewContractMonthFromCodes(months ...MonthCode) ContractMonths {
	var b strings.Builder

	for _, m := range months {
		b.WriteString(string(m))
	}

	return NewContractMonthFromString(b.String())
}

func NewContractMonthFromString(months string) ContractMonths {
	const pattern = `^[fghjkmnquvxz]+$`

	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(months) {
		panic(fmt.Sprintf("invalid contract months: %s", months))
	}

	return ContractMonths(months)
}

func (c ContractMonths) ForEach(fun func(month MonthCode)) {
	for _, m := range c {
		fun(MonthCode(m))
	}
}

func NewMonthCode(month string) MonthCode {
	const pattern = `^[fghjkmnquvxz]$`

	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(month) {
		panic(fmt.Sprintf("invalid month code: %s", month))
	}

	return MonthCode(month)
}

func NewMonthCodeFromMonth(month time.Month) MonthCode {
	switch month {
	case time.January:
		return January
	case time.February:
		return February
	case time.March:
		return March
	case time.April:
		return April
	case time.May:
		return May
	case time.June:
		return June
	case time.July:
		return July
	case time.August:
		return August
	case time.September:
		return September
	case time.October:
		return October
	case time.November:
		return November
	case time.December:
		return December
	default:
		panic("unknown month")
	}
}

func NewMonthCodeFromMonthValue(month int) MonthCode {
	switch month {
	case 1:
		return January
	case 2:
		return February
	case 3:
		return March
	case 4:
		return April
	case 5:
		return May
	case 6:
		return June
	case 7:
		return July
	case 8:
		return August
	case 9:
		return September
	case 10:
		return October
	case 11:
		return November
	case 12:
		return December
	default:
		panic("unknown month")
	}
}

func (m MonthCode) Month() time.Month {
	switch m {
	case January:
		return time.January
	case February:
		return time.February
	case March:
		return time.March
	case April:
		return time.April
	case May:
		return time.May
	case June:
		return time.June
	case July:
		return time.July
	case August:
		return time.August
	case September:
		return time.September
	case October:
		return time.October
	case November:
		return time.November
	case December:
		return time.December
	default:
		panic("unknown month code")
	}
}

func (m MonthCode) MonthValue() int {
	return int(m.Month())
}
