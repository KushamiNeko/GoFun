package continuous

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/KushamiNeko/go_fun/chart/data"
	"github.com/KushamiNeko/go_fun/chart/plotter"
	"github.com/KushamiNeko/go_fun/chart/preset"
	"github.com/KushamiNeko/go_fun/chart/utils"
	"gonum.org/v1/plot"
)

//func TestContinousPlot(t *testing.T) {
func plotContinousFutures(from, to time.Time, symbol, output string) error {
	//from, _ := time.Parse("20060102", "20171230")
	//to, _ := time.Parse("20060102", "20190101")

	src := data.NewDataSource(data.Continuous)

	cc, err := src.Read(
		from,
		to,
		symbol,
		data.Daily,
	)

	//cc, err := data.ContinousContract(
	//from,
	//to,
	//symbol,
	//futures.FinancialContractMonths,
	//futures.BarchartSymbolFormat,
	////data.OpenInterest,
	//futures.FirstOfMonth,
	//futures.Ratio,
	//)
	if err != nil {
		//t.Errorf(err.Error())
		return err
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Y.Scale = plot.LogScale{}
	p.BackgroundColor = preset.ThemeColor("ColorBackground")

	preset.EssentialTickSettings(
		p,
		&plotter.TimeTicker{
			TimeSeries: cc,
			Frequency:  data.Daily,
		},
		&plotter.PriceTicker{
			TimeSeries: cc,
			Step:       20,
		},
	)

	preset.EssentialIndicators(
		cc,
	)

	p.Add(
		plotter.GridPlotter(
			preset.ChartConfig("GridLineWidth"),
			preset.ThemeColor("ColorGrid"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("sma5"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorSMA1"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("sma20"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorSMA2"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("bb+15"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorBB1"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("bb-15"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorBB1"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("bb+20"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorBB2"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("bb-20"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorBB2"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("bb+25"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorBB3"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("bb-25"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorBB3"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("bb+30"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorBB4"),
		),
		plotter.LinePlotter(
			nil,
			cc.Values("bb-30"),
			preset.ChartConfig("LineWidth"),
			preset.ThemeColor("ColorBB4"),
		),
		&plotter.CandleStick{
			TimeSeries:   cc,
			ColorUp:      preset.ThemeColor("ColorUp"),
			ColorDown:    preset.ThemeColor("ColorDown"),
			ColorNeutral: preset.ThemeColor("ColorNeutral"),
			BodyWidth:    preset.ChartConfig("CandleBodyWidth"),
			ShadowWidth:  preset.ChartConfig("CandleShadowWidth"),
		},
	)

	ymin, ymax := utils.RangeExtend(
		utils.Min(cc.Values("low")),
		utils.Max(cc.Values("high")),
		25.0,
	)

	p.Y.Min = ymin
	p.Y.Max = ymax

	p.X.Min = -1
	p.X.Max = float64(len(cc.Times()))

	writer, err := p.WriterTo(
		preset.ChartConfig("ChartWidth"),
		preset.ChartConfig("ChartHeight"),
		"png",
	)
	if err != nil {
		panic(err)
	}

	dst, err := os.Create(output)
	if err != nil {
		panic(err)
	}

	defer dst.Close()

	_, err = writer.WriteTo(dst)
	if err != nil {
		panic(err)
	}

	return nil
}

func TestContinousPlot(t *testing.T) {

	symbol := "es"
	freq := "d"

	for i := 2000; i < 2019; i++ {
		from, _ := time.Parse("20060102", fmt.Sprintf("%d1231", i-1))
		to, _ := time.Parse("20060102", fmt.Sprintf("%d0101", i+1))
		err := plotContinousFutures(from, to, symbol, fmt.Sprintf("%d_%s_%s.png", i, symbol, freq))
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}
