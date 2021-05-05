package futures

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFrontContract(t *testing.T) {
	t.Parallel()

	tables := []map[string]interface{}{
		map[string]interface{}{
			"dtime":  "20181205",
			"symbol": "es",
			"months": FinancialContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "esh19",
		},
		map[string]interface{}{
			"dtime":  "20181005",
			"symbol": "es",
			"months": FinancialContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "esz18",
		},
		map[string]interface{}{
			"dtime":  "20180105",
			"symbol": "es",
			"months": FinancialContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "esh18",
		},
		map[string]interface{}{
			"dtime":  "20180605",
			"symbol": "es",
			"months": FinancialContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "esu18",
		},
		map[string]interface{}{
			"dtime":  "20180605",
			"symbol": "cl",
			"months": AllContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "cln18",
		},
		map[string]interface{}{
			"dtime":  "20181206",
			"symbol": "cl",
			"months": AllContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "clf19",
		},
		map[string]interface{}{
			"dtime":  "20180805",
			"symbol": "cl",
			"months": AllContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "clu18",
		},
		map[string]interface{}{
			"dtime":  "20180406",
			"symbol": "cl",
			"months": AllContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "clk18",
		},
		map[string]interface{}{
			"dtime":  "20180605",
			"symbol": "gc",
			"months": EvenContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "gcq18",
		},
		map[string]interface{}{
			"dtime":  "20180305",
			"symbol": "gc",
			"months": EvenContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "gcj18",
		},
		map[string]interface{}{
			"dtime":  "20180205",
			"symbol": "gc",
			"months": EvenContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "gcj18",
		},
		map[string]interface{}{
			"dtime":  "20181206",
			"symbol": "gc",
			"months": EvenContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "gcg19",
		},
		map[string]interface{}{
			"dtime":  "20181106",
			"symbol": "gc",
			"months": EvenContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "gcz18",
		},
	}

	for _, table := range tables {
		dt, _ := time.Parse(`20060102`, table["dtime"].(string))
		contract := FrontContract(dt, table["symbol"].(string), table["months"].(ContractMonths), table["format"].(SymbolFormat))

		if contract != table["expect"] {
			t.Errorf("expect: %s, get %s", table["expect"], contract)
		}
	}
}

func TestContractList(t *testing.T) {
	t.Parallel()

	tables := []map[string]interface{}{
		map[string]interface{}{
			"from":   "20170905",
			"to":     "20181205",
			"symbol": "es",
			"months": FinancialContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "esh19,esz18,esu18,esm18,esh18,esz17,esu17",
		},
		map[string]interface{}{
			"from":   "20170905",
			"to":     time.Now().Add(12 * 7 * 24 * time.Hour).Format("20060102"),
			"symbol": "es",
			"months": FinancialContractMonths,
			"format": BarchartSymbolFormat,
			"expect": fmt.Sprintf(
				"%s,%s",
				FrontContract(time.Now(), "es", FinancialContractMonths, BarchartSymbolFormat),
				"esz19,esu19,esm19,esh19,esz18,esu18,esm18,esh18,esz17,esu17",
			),
		},
		map[string]interface{}{
			"from":   "20171231",
			"to":     "20190101",
			"symbol": "es",
			"months": FinancialContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "esh19,esz18,esu18,esm18,esh18,esz17",
		},
		map[string]interface{}{
			"from":   "20171231",
			"to":     "20190101",
			"symbol": "gc",
			"months": EvenContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "gcg19,gcz18,gcv18,gcq18,gcm18,gcj18,gcg18,gcz17",
		},
		map[string]interface{}{
			"from":   "20171231",
			"to":     "20190101",
			"symbol": "cl",
			"months": AllContractMonths,
			"format": BarchartSymbolFormat,
			"expect": "clg19,clf19,clz18,clx18,clv18,clu18,clq18,cln18,clm18,clk18,clj18,clh18,clg18,clf18,clz17,clx17",
		},
	}

	for _, table := range tables {
		from, _ := time.Parse(`20060102`, table["from"].(string))
		to, _ := time.Parse(`20060102`, table["to"].(string))

		list := ContractList(
			from,
			to,
			table["symbol"].(string),
			table["months"].(ContractMonths),
			table["format"].(SymbolFormat),
		)

		elist := strings.Split(table["expect"].(string), ",")

		for i, c := range list {
			if c != elist[i] {
				t.Errorf("expect: %s, get: %s", elist[i], c)
			}
		}

	}
}

func TestPreviousContract(t *testing.T) {
	t.Parallel()

	tables := map[string][]map[string]interface{}{
		"inputs": []map[string]interface{}{
			map[string]interface{}{
				"contractMonths": FinancialContractMonths,
				"contract":       "esh20",
			},
			map[string]interface{}{
				"contractMonths": FinancialContractMonths,
				"contract":       "esz19",
			},
			map[string]interface{}{
				"contractMonths": FinancialContractMonths,
				"contract":       "esh00",
			},
			map[string]interface{}{
				"contractMonths": FinancialContractMonths,
				"contract":       "esm05",
			},
			map[string]interface{}{
				"contractMonths": FinancialContractMonths,
				"contract":       "esu98",
			},
		},
		"expect": []map[string]interface{}{
			map[string]interface{}{
				"previousContract": "esz19",
			},
			map[string]interface{}{
				"previousContract": "esu19",
			},
			map[string]interface{}{
				"previousContract": "esz99",
			},
			map[string]interface{}{
				"previousContract": "esh05",
			},
			map[string]interface{}{
				"previousContract": "esm98",
			},
		},
	}

	for i := range tables["inputs"] {
		pc := PreviousContract(
			tables["inputs"][i]["contract"].(string),
			tables["inputs"][i]["contractMonths"].(ContractMonths),
			BarchartSymbolFormat,
		)
		if pc != tables["expect"][i]["previousContract"] {
			t.Errorf(
				"expect: %s, get: %s",
				tables["expect"][i]["previousContract"],
				pc,
			)
		}
	}
}
