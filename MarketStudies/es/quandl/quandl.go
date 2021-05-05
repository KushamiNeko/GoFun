package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {
	const url = "https://www.quandl.com/api/v3/datatables/SCF/PRICES.csv?api_key=DowgUznF329gkXrXpf_-"

	resp, err := http.Get(
		fmt.Sprintf("%s&quandl_code=%s", url, "CME_MD1_OR"),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var buffer bytes.Buffer

	io.Copy(&buffer, resp.Body)

	fmt.Println(buffer.String())
}
