package operator

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type yahoo struct {
	*operator
}

func NewYahooOperator() *yahoo {
	y := &yahoo{operator: new(operator)}
	y.initDir()

	return y
}

func (y *yahoo) source() map[string]string {
	return map[string]string{
		"^vix": "19900101",
		"^vxn": "20000101",
		/* "^ovx": "20070101",
		"^gvz": "20100101", */
		//
		/* "^gspc": "19270101",
		"^ixic": "19710101",
		"^ndx":  "19850101",
		"^dji":  "19860101",
		"^rut":  "19880101",
		// "^sml": "19890101",
		"^n225": "19650101",
		"ezu":   "20000101",
		"eem":   "20030101",
		"^hsi":  "19860101",
		"fxi":   "20040101", */
		//
		"icsh": "20140101",
		"jpst": "20170701",
		"gsy":  "20080301",
		"near": "20130101",
		"shv":  "20070101",
		"flot": "20110701",
		"ushy": "20180101",
		"shyg": "20140101",
		"hyg":  "20070101",
		"sjnk": "20120501",
		"jnk":  "20080101",
		"faln": "20160701",
		"emb":  "20070101",
		"pcy":  "20071101",
		"emhy": "20120501",
		"bkln": "20110401",
		"srln": "20130501",
		"lqd":  "20020101",
		"mbb":  "20070401",
		// "mub":  "20071001",
		"shy": "20020801",
		"iei": "20070201",
		"ief": "20020801",
		"tlh": "20070201",
		"tlt": "20020801",
		//
		// "iyr":  "20000101",
		// "rem":  "20070101",
		// "reet": "20140801",
		//
		"pff": "20070401",
		"pgx": "20080201",
		// "pgf":  "20070101",
		"hdv": "20110401",
		"dvy": "20031201",
		// "idv": "20070701",
		// "dvye": "20120301",
		"spyd": "20151101",
		"sphd": "20121101",
		// "sdy":  "20051201",
		// "schd": "20111101",
		// "vym":  "20061201",
		// "dgro": "20140701",
		//
		// "usig": "20070201",
		// "slqd": "20140101",
		// "igsb": "20070201",
		// "igib": "20070201",
		// "iglb": "20100101",
		// "qlta": "20120301",
		// "lqdh": "20140701",
		// "govt": "20120301",
		// "igov": "20090201",
		// "stip": "20110101",
		// "tip": "20040101",
		// "sub": "20090101",
		// "icvt": "20150701",
		// "istb": "20130101",
		// "iusb": "20140701",
		// "agg": "20040101",
		// "byld": "20140501",
	}
}

func (y *yahoo) url(symbol, datetime string) string {

	dt, err := time.Parse("20060102", datetime)
	if err != nil {
		panic(err)
	}

	var b strings.Builder

	b.WriteString(
		fmt.Sprintf(
			"https://finance.yahoo.com/quote/%s/history?",
			url.PathEscape(symbol),
		),
	)

	b.WriteString(
		fmt.Sprintf(
			"period1=%d&",
			dt.Unix(),
		),
	)

	b.WriteString(
		fmt.Sprintf(
			"period2=%d&",
			time.Now().Add(5*24*time.Hour).Unix(),
		),
	)

	b.WriteString("interval=1d&filter=history&frequency=1d")

	return b.String()
}

func (y *yahoo) Greeting() {
	y.greetingMessage("Yahoo")
}

func (y *yahoo) Download() {

	for symbol, datetime := range y.source() {
		y.download(y.url(symbol, datetime), symbol)
	}

	y.downloadCompleted()
}

func (y *yahoo) renameSymbol(symbol string) string {
	switch symbol {
	case "n225":
		return "nikk"
	case "gspc":
		return "spx"
	case "ixic":
		return "compq"
	default:
		return symbol
	}
}

func (y *yahoo) Rename() {
	indexRegex := regexp.MustCompile(`(\^(\w+)).csv`)

	fs, err := os.ReadDir(y.srcDir)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		symbol := strings.ToLower(strings.ReplaceAll(f.Name(), filepath.Ext(f.Name()), ""))

		if _, ok := y.source()[symbol]; ok {
			match := indexRegex.FindAllStringSubmatch(f.Name(), -1)
			if len(match) != 0 {
				symbol = strings.ToLower(match[0][2])
			}

			symbol = strings.ReplaceAll(symbol, "-", "")

			symbol = y.renameSymbol(symbol)

			srcPath := filepath.Join(y.srcDir, f.Name())
			dstPath := filepath.Join(y.dstDir, "yahoo", fmt.Sprintf("%s.csv", symbol))

			y.rename(srcPath, dstPath)

		}
	}

	y.renameCompleted()

}

func (y *yahoo) Check() {
	regex := regexp.MustCompile(`^(\^*(\w+))$`)

	for symbol := range y.source() {
		match := regex.FindAllStringSubmatch(symbol, -1)
		if len(match) == 0 {
			panic("invalid yahoo check regex")
		}

		symbol := strings.ToLower(match[0][2])
		symbol = strings.ReplaceAll(symbol, "-", "")
		symbol = y.renameSymbol(symbol)

		path := filepath.Join(y.dstDir, "yahoo", fmt.Sprintf("%s.csv", symbol))
		y.check(path)
	}

	y.checkCompleted()
}
