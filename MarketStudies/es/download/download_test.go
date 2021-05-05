package download

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/KushamiNeko/go_fun/chart/data"
	"github.com/KushamiNeko/go_fun/chart/futures"
)

//func TestOpenHistoricalPage(t *testing.T) {
//historicalPage()
//}

func download(ys, ye int, dst string, symbols, monthCodes []string, format futures.SymbolFormat, src data.DataSource) {
	var wg sync.WaitGroup

	for _, symbol := range symbols {
		for y := ys; y < ye; y++ {
			for _, m := range monthCodes {

				wg.Add(1)

				var code string
				switch format {
				case futures.BarchartSymbolFormat:
					code = fmt.Sprintf("%s%s%02d", symbol, m, y%100)
				case futures.QuandlSymbolFormat:
					code = fmt.Sprintf("%s%s%d", symbol, m, y)
				default:
					panic("unknown symbol format")
				}

				//go func(symbol string, y int, m string) {
				go func(symbol, code string, src data.DataSource) {

					//code := fmt.Sprintf("%s%s%02d", symbol, m, y%100)
					fmt.Println(code)

					//from, _ := time.Parse(`20060102`, "19900101")
					//to := time.Now()
					//series, err := src.Read(
					//from,
					//to,
					//code,
					//data.Daily,
					//)
					//if err != nil {
					//panic(err)
					//}

					//contract := series.CSV()

					contract, err := readOnDemand(code)
					if err != nil {
						//t.Errorf(err.Error())
						panic(err)
					}

					f, err := os.Create(
						filepath.Join(
							//output,
							dst,
							symbol,
							fmt.Sprintf("%s.csv", code),
						),
					)
					if err != nil {
						//t.Errorf(err.Error())
						panic(err)
					}

					_, err = f.WriteString(contract)
					if err != nil {
						//t.Errorf(err.Error())
						panic(err)
					}

					f.Sync()
					f.Close()

					wg.Done()

					//}(symbol, y, m)
				}(symbol, code, src)

			}
		}

	}

	wg.Wait()
}

func TestDownloadBarchartOnDemand(t *testing.T) {
	//t.SkipNow()

	dst := filepath.Join(
		os.Getenv("HOME"),
		"Documents/data_source/continuous",
	)

	var (
		symbols    []string
		monthCodes []string
	)

	//symbols = []string{
	////"ge",
	////"ny",
	////"nl",
	////"es",
	////"nq",
	////"qr",
	////"zn",
	//"gc",
	//"cl",
	//}

	//monthCodes = []string{
	////"h", "m", "u", "z",
	//"f", "g", "h", "j", "k", "m", "n", "q", "u", "v", "x", "z",
	//}

	//download(1997, 2020, dst, symbols, monthCodes, futures.BarchartSymbolFormat, data.NewDataSource(data.Barchart))

	symbols = []string{
		"ge",
		"ny",
		"nl",
		"es",
		"nq",
		"qr",
		"zn",
		//"gc",
		//"cl",
	}

	monthCodes = []string{
		"h", "m", "u", "z",
		//"f", "g", "h", "j", "k", "m", "n", "q", "u", "v", "x", "z",
	}

	download(1997, 1999, dst, symbols, monthCodes, futures.BarchartSymbolFormat, data.NewDataSource(data.Barchart))
}

func TestDownloadQuandl(t *testing.T) {
	t.SkipNow()

	dst := filepath.Join(
		os.Getenv("HOME"),
		"Documents/data_source/continuous",
	)

	symbols := []string{
		//"fesx",
		//"fgbl",
		//"nk225",
		"nk225m",
	}

	monthCodes := []string{
		//"h", "m", "u", "z",
		"f", "g", "h", "j", "k", "m", "n", "q", "u", "v", "x", "z",
	}

	download(2016, 2020, dst, symbols, monthCodes, futures.QuandlSymbolFormat, data.NewDataSource(data.Quandl))
}
