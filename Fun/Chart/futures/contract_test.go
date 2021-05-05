package futures

import (
	"testing"
	"time"
)

func TestContract(t *testing.T) {
	c := &Contract{
		symbol: "esh20",
		format: BarchartSymbolFormat,
	}

	if c.ContractYear() != 2020 {
		t.Errorf("espect: 2020, get: %d", c.ContractYear())
	}

	if c.ContractMonth() != time.March {
		t.Errorf("espect: march, get: %d", c.ContractMonth())
	}

	c = &Contract{
		symbol: "esu99",
		format: BarchartSymbolFormat,
	}

	if c.ContractYear() != 1999 {
		t.Errorf("espect: 1999, get: %d", c.ContractYear())
	}

	if c.ContractMonth() != time.September {
		t.Errorf("espect: september, get: %d", c.ContractMonth())
	}

	c = &Contract{
		symbol: "nk225mm1999",
		format: QuandlSymbolFormat,
	}

	if c.ContractYear() != 1999 {
		t.Errorf("espect: 1999, get: %d", c.ContractYear())
	}

	if c.ContractMonth() != time.June {
		t.Errorf("espect: june, get: %d", c.ContractMonth())
	}

	c = &Contract{
		symbol: "nk225mu2019",
		format: QuandlSymbolFormat,
	}

	if c.ContractYear() != 2019 {
		t.Errorf("espect: 2019, get: %d", c.ContractYear())
	}

	if c.ContractMonth() != time.September {
		t.Errorf("espect: september, get: %d", c.ContractMonth())
	}
}
