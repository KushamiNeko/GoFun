package input

import (
	"fmt"
	"testing"

	"github.com/KushamiNeko/GoFun/Fun/Utility/test"
)

func TestValidateStringWIthRegex(t *testing.T) {
	cases := []struct {
		pattern string
		input   string
		ok      bool
		panic   bool
	}{
		{
			pattern: `^\d+$`,
			input:   "123",
			ok:      true,
			panic:   false,
		},
		{
			pattern: `^\d+$`,
			input:   "abc",
			ok:      false,
			panic:   false,
		},
		{
			pattern: `(`,
			input:   "",
			ok:      false,
			panic:   true,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case@%d", i), func(t *testing.T) {
			if c.panic {
				test.ShouldPanic(t, func() {
					ValidateStringWithRegex(c.pattern, c.input)
				})
			} else {
				ok := ValidateStringWithRegex(c.pattern, c.input)
				if ok != c.ok {
					t.Errorf("expect %v but get %v", c.ok, ok)
				}
			}
		})
	}

}

func TestValidKeyValuePair(t *testing.T) {
	type caseOut struct {
		date   string
		symbol string
	}

	type caseErr struct {
		input  bool
		date   bool
		symbol bool
	}

	cases := []struct {
		input  string
		output caseOut
		err    caseErr
	}{
		{
			input: "d=20200101; s=rty",
			output: caseOut{
				date:   "20200101",
				symbol: "rty",
			},
			err: caseErr{
				input:  false,
				date:   false,
				symbol: false,
			},
		},
		{
			input: " d=20200101;s=rty; ",
			output: caseOut{
				"20200101",
				"rty",
			},
			err: caseErr{
				input:  false,
				date:   false,
				symbol: false,
			},
		},
		{
			input: "d=20190101;s=es;",
			output: caseOut{
				"20190101",
				"es",
			},
			err: caseErr{
				input:  false,
				date:   false,
				symbol: false,
			},
		},
		{
			input: "",
			output: caseOut{
				"",
				"",
			},
			err: caseErr{
				input:  true,
				date:   false,
				symbol: false,
			},
		},
		{
			input: "s=es;",
			output: caseOut{
				"",
				"es",
			},
			err: caseErr{
				input:  false,
				date:   true,
				symbol: false,
			},
		},
		{
			input: "d=20200101",
			output: caseOut{
				"20200101",
				"",
			},
			err: caseErr{
				input:  false,
				date:   false,
				symbol: true,
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case@%d", i), func(t *testing.T) {
			pair, err := keyValuePair(c.input)

			switch c.err.input {
			case true:
				if err == nil {
					t.Errorf("err should be not be nil")
				}

				return
			case false:
				if err != nil {
					t.Errorf("err should be nil but get %s", err.Error())
				}

			}

			d, ok := pair["d"]
			switch c.err.date {
			case true:
				if ok {
					t.Errorf("d should be empty")
				}
			case false:
				if !ok {
					t.Errorf("d should contain date")
				}
			}

			if d != c.output.date {
				t.Errorf("d should be 20200101 but get %s", d)
			}

			s, ok := pair["s"]
			switch c.err.symbol {
			case true:
				if ok {
					t.Errorf("s should be empty")
				}
			case false:
				if !ok {
					t.Errorf("s should contain symbol")
				}
			}

			if s != c.output.symbol {
				t.Errorf("s should be rty but get %s", s)
			}
		})
	}

}

func TestInputAbbreviation(t *testing.T) {
	a := map[string]string{
		"s": "symbol",
		"p": "price",
	}

	i := map[string]string{
		"s": "rty",
		"p": "100.25",
		"d": "20190521",
	}

	v := InputsAbbreviation(i, a)
	s, ok := v["symbol"]
	if !ok {
		t.Errorf("v should contain symbol")
	}

	if s != "rty" {
		t.Errorf("s should be rty")
	}

	p, ok := v["price"]
	if !ok {
		t.Errorf("v should contain price")
	}

	if p != "100.25" {
		t.Errorf("p should be 100.25")
	}

	d, ok := v["d"]
	if !ok {
		t.Errorf("v should contain date")
	}

	if d != "20190521" {
		t.Errorf("d should be 20190521")
	}
}
