package main

//func readSymbol(from, to, symbol string, freq data.Frequency) *data.TimeSeries {
//var src data.DataSource

//switch {
//case symbol == "vix" || symbol == "vxn":
//src = data.NewDataSource(data.Yahoo)
//default:
//src = data.NewDataSource(data.StockCharts)
//}

//f, err := time.Parse("20060102", from)
//if err != nil {
//panic(err)
//}

//t, err := time.Parse("20060102", to)
//if err != nil {
//panic(err)
//}

//ts, err := src.Read(f, t, symbol, freq)
//if err != nil {
//panic(err)
//}

//return ts
//}

//var (
//vix     *data.TimeSeries
//vxn     *data.TimeSeries
//spxhilo *data.TimeSeries
//ndxhilo *data.TimeSeries
//nyhilo  *data.TimeSeries
//nahilo  *data.TimeSeries
//ryratmm *data.TimeSeries
//)

//func printData(date string) string {
//s := make([]string, 0)

//from := "19900101"
//to := "20200101"

//freq := data.Daily

//d, err := time.Parse("20060102", date)
//if err != nil {
//panic(err)
//}

//if vix == nil {
//vix = readSymbol(from, to, "vix", freq)
//}

//if vxn == nil {
//vxn = readSymbol(from, to, "vxn", freq)
//}

//if spxhilo == nil {
//spxhilo = readSymbol(from, to, "spxhilo", freq)
//}

//if ndxhilo == nil {
//ndxhilo = readSymbol(from, to, "ndxhilo", freq)
//}

//if nyhilo == nil {
//nyhilo = readSymbol(from, to, "nyhilo", freq)
//}

//if nahilo == nil {
//nahilo = readSymbol(from, to, "nahilo", freq)
//}

//if ryratmm == nil {
//ryratmm = readSymbol(from, to, "ryratmm", freq)
//}

//s = append(s, fmt.Sprintf("VIX: %.2f", vix.ValueInTimes(d, "close", 0)))
//s = append(s, fmt.Sprintf("VXN: %.2f", vxn.ValueInTimes(d, "close", 0)))
//s = append(s, fmt.Sprintf("SPXHILO: %.2f", spxhilo.ValueInTimes(d, "close", 0)))
//s = append(s, fmt.Sprintf("NDXHILO: %.2f", ndxhilo.ValueInTimes(d, "close", 0)))
//s = append(s, fmt.Sprintf("NYHILO: %.2f", nyhilo.ValueInTimes(d, "close", 0)))
//s = append(s, fmt.Sprintf("NAHILO: %.2f", nahilo.ValueInTimes(d, "close", 0)))
//s = append(s, fmt.Sprintf("RYRATMM: %.2f", ryratmm.ValueInTimes(d, "close", 0)))

//return strings.Join(s, "\n")
//}

//func createNote(book string) {

//tradingAgent, err := newTradingAgent(book)
//if err != nil {
//panic(err)
//}

//ts, err := tradingAgent.Transactions()
//if err != nil {
//panic(err)
//}

//ns := make([]map[string]string, 0, len(ts))

//var last time.Time

//for _, t := range ts {
//if t.Time().Equal(last) {
//continue
//}

//d := t.Time().Format("20060102")
//ns = append(ns, map[string]string{
//"time": d,
//"note": printData(d),
//})

//last = t.Time()
//}

//jd, err := json.MarshalIndent(ns, "", "  ")
//if err != nil {
//panic(err)
//}

//b, _ := tradingAgent.Reading()

//f, err := os.Create(fmt.Sprintf("trading_note_%s.json", b.NoteIndex()))
//if err != nil {
//panic(err)
//}

//defer f.Close()

//f.Write(jd)
//}
