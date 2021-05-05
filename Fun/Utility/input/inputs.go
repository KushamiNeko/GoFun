package input

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//func ValidateWithRegex(input, pattern string) bool {
//regex := regexp.MustCompile(pattern)
//return regex.MatchString(input)
//}

func ValidateStringWithRegex(pattern, input string) bool {
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(input)
}

//func MustValidateWithRegex(input, pattern string) {
//if !ValidateWithRegex(input, pattern) {
//panic(fmt.Sprintf(`invalid input "%s" with regex "%s"`, input, pattern))
//}
//}

func keyValuePair(inputs string) (map[string]string, error) {
	const pattern = `([^;]*)=([^;]*)`

	pair := make(map[string]string)

	re := regexp.MustCompile(pattern)

	match := re.FindAllStringSubmatch(inputs, -1)

	for _, m := range match {
		k := strings.TrimSpace(m[1])
		v := strings.TrimSpace(m[2])

		if k == "" {
			return nil, fmt.Errorf("empty key: %s=%s", k, v)
		}

		if v == "" {
			return nil, fmt.Errorf("empty value: %s=%s", k, v)
		}

		pair[k] = v
	}

	if len(pair) == 0 {
		return nil, fmt.Errorf("empty inputs")
	}

	return pair, nil
}

func InputsAbbreviation(inputs, abbreviation map[string]string) map[string]string {
	n := make(map[string]string)

	for k, v := range inputs {
		fk, ok := abbreviation[k]
		if ok {
			n[fk] = v
		} else {
			n[k] = v
		}
	}

	return n
}

func YearRangeInput(input string) (int, int) {
	const pattern = `^(\d{4})(?:(?:\-|\~)(\d{4}))*$`

	var err error

	regex := regexp.MustCompile(pattern)
	m := regex.FindAllStringSubmatch(input, -1)

	s := m[0][1]
	e := m[0][2]

	var sy, ey int64

	sy, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid year: %s", s))
	}

	if e == "" {
		ey = sy + 1
	} else {
		ey, err = strconv.ParseInt(e, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid year: %s", e))
		}

		//ey += 1
	}

	if ey <= sy {
		panic(fmt.Sprintf("invalid years: %d - %d", sy, ey))
	}

	return int(sy), int(ey)
}
