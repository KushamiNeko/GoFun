package download

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func historicalPage() {

	const interactive = `https://www.barchart.com/futures/quotes/%s/interactive-chart`
	const historical = `https://www.barchart.com/futures/quotes/%s/historical-download`

	symbols := []string{
		//"fx",
		//"np",
		//"gg",
		//"tj",
		"es",
		"nq",
		"qr",
		"zn",
	}

	months := []string{
		"h", "m", "u", "z",
		//"f", "g", "h", "j", "k", "m", "n", "q", "u", "v", "x", "z",
	}

	ys := 1998
	ye := 1999

	for _, symbol := range symbols {
		for y := ys; y < ye; y++ {
			for _, month := range months {

				var outb bytes.Buffer
				var errb bytes.Buffer

				cmd := exec.Command(
					"google-chrome",
					fmt.Sprintf(
						`https://www.barchart.com/futures/quotes/%s/historical-download`,
						fmt.Sprintf("%s%s%02d", symbol, month, y),
					),
				)

				cmd.Stdout = &outb
				cmd.Stderr = &errb

				err := cmd.Start()
				if err != nil {
					panic(err)
				}

			}

			time.Sleep(20 * time.Second)
		}
	}
}

func readOnDemand(symbol string) (string, error) {

	const root = `https://ondemand.websol.barchart.com/getHistory.csv?`

	const timeFormat = "20060102"

	url := fmt.Sprintf(
		"%sapikey=%s&symbol=%s&type=%s&startDate=%s&endDate=%s&interval=%s&volume=%s&backAdjust=%s&contractRoll=%s",
		root,
		os.Getenv("BARCHART"),
		symbol,
		"daily",
		"19900101",
		time.Now().Format(timeFormat),
		"1",
		"contract",
		"false",
		"combined",
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil

}
