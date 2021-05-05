package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/KushamiNeko/go_fun/chart/data"
	cp "github.com/KushamiNeko/go_fun/chart/plotter"
	"github.com/KushamiNeko/go_fun/trading/agent"
	"github.com/KushamiNeko/go_fun/trading/utils"
	"github.com/KushamiNeko/go_fun/utils/pretty"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func VixToRisk(year int, op string) ([]float64, []float64) {
	a, err := agent.NewTradingAgentCompact(
		filepath.Join(
			os.Getenv("HOME"),
			"Documents/database/filedb/futures_wizards",
		),
		"aa",
		fmt.Sprintf("es_%d", year),
	)
	if err != nil {
		panic(err)
	}

	trades, err := a.Trades()
	if err != nil {
		panic(err)
	}

	records, err := a.Transactions()
	if err != nil {
		panic(err)
	}

	f := time.Date(
		year,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.Now().Location(),
	)

	t := time.Date(
		year+1,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.Now().Location(),
	)

	exf := f.Add(-500 * 24 * time.Hour)
	ext := t.Add(500 * 24 * time.Hour)

	ssrc := data.NewDataSource(data.StockCharts)
	ysrc := data.NewDataSource(data.Yahoo)

	es, err := ssrc.Read(exf, ext, "es", data.Daily)
	if err != nil {
		panic(err)
	}

	vix, err := ysrc.Read(exf, ext, "vix", data.Daily)
	if err != nil {
		panic(err)
	}

	vs := make([]float64, 0)
	rs := make([]float64, 0)

	for _, trade := range trades {
		if trade.Operation() != op {
			continue
		}

		period := utils.TradeEntry(trade.Transactions(), trade.Operation())

		for _, p := range period.Times() {
			vs = append(vs, vix.ValueInTimes(p, "close", 0))
		}

		risks := utils.CalculateRisk(
			es,
			records,
			period.From(),
			period.To(),
			period.Flip(),
			trade.Operation(),
			false,
		)

		for _, r := range risks {
			if !r.Combined() {
				if r.Risk() < -1.0 {
					fmt.Printf("%s: %.4f%%  VIX: %.2f\n", r.Time().Format("20060102"), r.Risk(), vix.ValueInTimes(r.Time(), "close", 0))
				}

				rs = append(rs, r.Risk())
			}
		}

		//fmt.Printf("%s ~ %s\n", period[0].Format("20060102"), period[len(period)-1].Format("20060102"))
	}

	pos, err := a.Positions()
	if err != nil {
		panic(err)
	}

	if len(pos) > 0 {
		if pos[0].Operation() == op {
			period := utils.TradeEntry(pos, pos[0].Operation())

			for _, p := range period.Times() {
				vs = append(vs, vix.ValueInTimes(p, "close", 0))
			}

			risks := utils.CalculateRisk(
				es,
				records,
				period.From(),
				period.To(),
				period.Flip(),
				pos[0].Operation(),
				false,
			)

			for _, r := range risks {
				if !r.Combined() {
					if r.Risk() < -1.0 {
						//fmt.Printf("%s: %.4f%%\n", r.Time().Format("20060102"), r.Risk())
						fmt.Printf("%s: %.4f%%  VIX: %.2f\n", r.Time().Format("20060102"), r.Risk(), vix.ValueInTimes(r.Time(), "close", 0))
					}

					rs = append(rs, r.Risk())
				}
			}
		}
	}

	return vs, rs
}

func XYs(v, f []float64) plotter.XYs {
	if len(v) != len(f) {
		panic("length should be the same")
	}

	pts := make(plotter.XYs, 0, len(f))
	for i := range f {
		if math.IsNaN(v[i]) || math.IsNaN(f[i]) {
			continue
		}

		pts = append(pts, plotter.XY{X: v[i], Y: f[i]})
	}

	return pts
}

func main() {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Vix to Risk"
	p.X.Label.Text = "Vix"
	p.Y.Label.Text = "Risk"

	p.X.Tick.Marker = &cp.BasicTicker{Step: 1, Digit: 2}
	p.Y.Tick.Marker = &cp.BasicTicker{Step: 0.05, Digit: 2}

	p.Add(plotter.NewGrid())

	vixs := make([]float64, 0)
	risks := make([]float64, 0)

	ys := 2016
	ye := 2017

	for i := ys; i < ye; i++ {
		vs, rs := VixToRisk(i, "+")

		vixs = append(vixs, vs...)
		risks = append(risks, rs...)
	}

	xys := XYs(vixs, risks)

	s, err := plotter.NewScatter(xys)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = pretty.HexToColor(pretty.PaperBlue300, 0.8)
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	s.GlyphStyle.Radius = 5

	p.Add(s)
	p.Legend.Add("+", s)

	vixs = make([]float64, 0)
	risks = make([]float64, 0)

	for i := ys; i < ye; i++ {
		vs, rs := VixToRisk(i, "-")

		vixs = append(vixs, vs...)
		risks = append(risks, rs...)
	}

	xys = XYs(vixs, risks)

	s, err = plotter.NewScatter(xys)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = pretty.HexToColor(pretty.PaperRed300, 0.8)
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	s.GlyphStyle.Radius = 5

	p.Add(s)
	p.Legend.Add("-", s)

	err = p.Save(10*vg.Inch, 10*vg.Inch, "vix_risk.png")
	if err != nil {
		panic(err)
	}
}
