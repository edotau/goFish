package api

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/vdobler/chart"
	"github.com/vdobler/chart/imgg"
	"github.com/vdobler/chart/svgg"
	"github.com/vdobler/chart/txtg"
)

type Dumper struct {
	N, M, W, H, Cnt           int
	S                         *svg.SVG
	I                         *image.RGBA
	svgFile, imgFile, txtFile *os.File
}

func NewDumper(name string, n, m, w, h int) *Dumper {
	var err error
	dumper := Dumper{N: n, M: m, W: w, H: h}

	dumper.svgFile, err = os.Create(name + ".svg")
	if err != nil {
		panic(err)
	}
	dumper.S = svg.New(dumper.svgFile)
	dumper.S.Start(n*w, m*h)
	dumper.S.Title(name)
	dumper.S.Rect(0, 0, n*w, m*h, "fill: #ffffff")

	dumper.imgFile, err = os.Create(name + ".png")
	if err != nil {
		panic(err)
	}
	dumper.I = image.NewRGBA(image.Rect(0, 0, n*w, m*h))
	bg := image.NewUniform(color.RGBA{0xff, 0xff, 0xff, 0xff})
	draw.Draw(dumper.I, dumper.I.Bounds(), bg, image.ZP, draw.Src)

	dumper.txtFile, err = os.Create(name + ".txt")
	if err != nil {
		panic(err)
	}

	return &dumper
}

func (d *Dumper) Plot(c chart.Chart) {
	row, col := d.Cnt/d.N, d.Cnt%d.N

	igr := imgg.AddTo(d.I, col*d.W, row*d.H, d.W, d.H, color.RGBA{0xff, 0xff, 0xff, 0xff}, nil, nil)
	c.Plot(igr)

	sgr := svgg.AddTo(d.S, col*d.W, row*d.H, d.W, d.H, "", 12, color.RGBA{0xff, 0xff, 0xff, 0xff})
	c.Plot(sgr)

	tgr := txtg.New(100, 30)
	c.Plot(tgr)
	d.txtFile.Write([]byte(tgr.String() + "\n\n\n"))

	d.Cnt++

}
func (d *Dumper) Close() {
	png.Encode(d.imgFile, d.I)
	d.imgFile.Close()

	d.S.End()
	d.svgFile.Close()

	d.txtFile.Close()
}

func Kernels() {
	dumper := NewDumper("xkernels", 1, 1, 600, 400)
	defer dumper.Close()

	p := chart.ScatterChart{Title: "Kernels"}
	p.XRange.Label, p.YRange.Label = "u", "K(u)"
	p.XRange.MinMode.Fixed, p.XRange.MaxMode.Fixed = true, true
	p.XRange.MinMode.Value, p.XRange.MaxMode.Value = -2, 2
	p.YRange.MinMode.Fixed, p.YRange.MaxMode.Fixed = true, true
	p.YRange.MinMode.Value, p.YRange.MaxMode.Value = -0.1, 1.1

	p.XRange.TicSetting.Delta = 1
	p.YRange.TicSetting.Delta = 0.2
	p.XRange.TicSetting.Mirror = 1
	p.YRange.TicSetting.Mirror = 1

	p.AddFunc("Bisquare", chart.BisquareKernel,
		chart.PlotStyleLines, chart.Style{Symbol: 'o', LineWidth: 1, LineColor: color.NRGBA{0xa0, 0x00, 0x00, 0xff}, LineStyle: 1})
	p.AddFunc("Epanechnikov", chart.EpanechnikovKernel,
		chart.PlotStyleLines, chart.Style{Symbol: 'X', LineWidth: 1, LineColor: color.NRGBA{0x00, 0xa0, 0x00, 0xff}, LineStyle: 1})
	p.AddFunc("Rectangular", chart.RectangularKernel,
		chart.PlotStyleLines, chart.Style{Symbol: '=', LineWidth: 1, LineColor: color.NRGBA{0x00, 0x00, 0xa0, 0xff}, LineStyle: 1})
	p.AddFunc("Gauss", chart.GaussKernel,
		chart.PlotStyleLines, chart.Style{Symbol: '*', LineWidth: 1, LineColor: color.NRGBA{0xa0, 0x00, 0xa0, 0xff}, LineStyle: 1})

	dumper.Plot(&p)
}

func ScatterChart(filename string, x, y []float64) {
	dumper := NewDumper(filename, 1, 1, 800, 600)
	defer dumper.Close()

	pl := chart.ScatterChart{Title: "Scatter + Lines"}
	pl.XRange.Label, pl.YRange.Label = "X - Value", "Y - Value"
	pl.Key.Pos = "itl"
	// pl.XRange.TicSetting.Delta = 5
	pl.XRange.TicSetting.Grid = 1

	pl.AddDataPair("Data", x, y, chart.PlotStyleLinesPoints,
		chart.Style{Symbol: '#', SymbolColor: color.NRGBA{0x00, 0x00, 0xff, 0xff}, LineStyle: chart.SolidLine})
	last := len(pl.Data) - 1
	pl.Data[last].Samples[6].DeltaX = 2.5
	pl.Data[last].Samples[6].OffX = 0.5
	pl.Data[last].Samples[6].DeltaY = 16
	pl.Data[last].Samples[6].OffY = 2

	pl.AddData("Points", []chart.EPoint{{-4, 40, 0, 0, 0, 0}, {-3, 45, 0, 0, 0, 0},
		{-2, 35, 0, 0, 0, 0}}, chart.PlotStylePoints,
		chart.Style{Symbol: '0', SymbolColor: color.NRGBA{0xff, 0x00, 0xff, 0xff}})
	pl.AddFunc("Theory", func(x float64) float64 {
		if x > 5.25 && x < 5.75 {
			return 75
		}
		if x > 7.25 && x < 7.75 {
			return 500
		}
		return x * x
	}, chart.PlotStyleLines,
		chart.Style{Symbol: '%', LineWidth: 2, LineColor: color.NRGBA{0xa0, 0x00, 0x00, 0xff}, LineStyle: chart.DashDotDotLine})
	pl.AddFunc("30", func(x float64) float64 { return 30 }, chart.PlotStyleLines,
		chart.Style{Symbol: '+', LineWidth: 1, LineColor: color.NRGBA{0x00, 0xa0, 0x00, 0xff}, LineStyle: 1})
	pl.AddFunc("", func(x float64) float64 { return 7 }, chart.PlotStyleLines,
		chart.Style{Symbol: '@', LineWidth: 1, LineColor: color.NRGBA{0x00, 0x00, 0xa0, 0xff}, LineStyle: 1})

	pl.XRange.ShowZero = true
	pl.XRange.TicSetting.Mirror = 1
	pl.YRange.TicSetting.Mirror = 1
	pl.XRange.TicSetting.Grid = 1
	pl.XRange.Label = "X-Range"
	pl.YRange.Label = "Y-Range"
	pl.Key.Cols = 2
	pl.Key.Pos = "orb"

	dumper.Plot(&pl)
}
