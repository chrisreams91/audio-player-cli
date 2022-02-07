package main

import (
	"fmt"
	"image"
	"sync"

	ui "github.com/gizak/termui/v3"
)

func buildTopBars(data[] float64) *BarChart{
	bc := NewBarChart(false)
	bc.SetRect(5, 0, 100, 15)

	bc.BarWidth = 1
	bc.BarGap = 0
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.NumFormatter = func(f float64) string { return " "}
	
	return bc
}

func buildBottomBars(data[] float64) *BarChart{ 
	bc := NewBarChart(true)
	bc.SetRect(5, 15, 100, 30)

	bc.BarWidth = 1
	bc.BarGap = 0
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
	bc.NumFormatter = func(f float64) string { return " "}
	return bc
}

type BarChart struct {
	Block
	BarColors    []ui.Color
	LabelStyles  []ui.Style
	NumStyles    []ui.Style // only Fg and Modifier are used
	NumFormatter func(float64) string
	Data         []float64
	Labels       []string
	BarWidth     int
	BarGap       int
	MaxVal       float64
	inverted     bool
}

func NewBarChart(inverted bool) *BarChart {
	return &BarChart{
		Block:        *NewBlock(),
		BarColors:    ui.Theme.BarChart.Bars,
		NumStyles:    ui.Theme.BarChart.Nums,
		LabelStyles:  ui.Theme.BarChart.Labels,
		NumFormatter: func(n float64) string { return fmt.Sprint(n) },
		BarGap:       1,
		BarWidth:     3,
		inverted: inverted,
	}
}

func (self *BarChart) Draw(buf *ui.Buffer) {
	self.Block.Draw(buf)

	maxVal := self.MaxVal
	if maxVal == 0 {
		maxVal, _ = ui.GetMaxFloat64FromSlice(self.Data)
	}

	barXCoordinate := self.Inner.Min.X
	for i, data := range self.Data {
		// draw bar
		height := int((data / maxVal) * float64(self.Inner.Dy()-2))
		for x := barXCoordinate; x < ui.MinInt(barXCoordinate+self.BarWidth, self.Inner.Max.X); x++ {
			for y := self.Inner.Max.Y ; y > (self.Inner.Max.Y-1)-height; y-- {
				c := ui.NewCell(' ', ui.NewStyle(ui.ColorClear, ui.SelectColor(self.BarColors, i)))
				if (self.inverted) {
					buf.SetCell(c, image.Pt(x, self.Max.Y - y + 15))
				} else {
					buf.SetCell(c, image.Pt(x, y))
				}
			}
		}

		numberXCoordinate := barXCoordinate + int((float64(self.BarWidth) / 2))
		if numberXCoordinate <= self.Inner.Max.X {
			// barLocation := self.Inner.Max.Y
		// 	if (self.inverted) {
		// 	barLocation = self.Inner.Min.Y
		// 	buf.SetString(
		// 		self.NumFormatter(data),
		// 		ui.NewStyle(
		// 			ui.SelectStyle(self.NumStyles, i+1).Fg,
		// 			ui.SelectColor(self.BarColors, i),
		// 			ui.SelectStyle(self.NumStyles, i+1).Modifier,
		// 		),
		// 		image.Pt(numberXCoordinate, barLocation - 10),
		// 	)
		// }

		}

		barXCoordinate += (self.BarWidth + self.BarGap)
	}
}


// Block is the base struct inherited by most widgets.
// Block manages size, position, border, and title.
// It implements all 3 of the methods needed for the `Drawable` interface.
// Custom widgets will override the Draw method.
type Block struct {
	Border      bool
	BorderStyle ui.Style

	BorderLeft, BorderRight, BorderTop, BorderBottom bool

	PaddingLeft, PaddingRight, PaddingTop, PaddingBottom int

	image.Rectangle
	Inner image.Rectangle

	Title      string
	TitleStyle ui.Style

	sync.Mutex
}

func NewBlock() *Block {
	return &Block{
		Border:       true,
		BorderStyle:  ui.Theme.Block.Border,
		BorderLeft:   true,
		BorderRight:  true,
		BorderTop:    true,
		BorderBottom: true,

		TitleStyle: ui.Theme.Block.Title,
	}
}

// Draw implements the Drawable interface.
func (self *Block) Draw(buf *ui.Buffer) {
	buf.SetString(
		self.Title,
		self.TitleStyle,
		image.Pt(self.Min.X+2, self.Min.Y),
	)
}

// SetRect implements the Drawable interface.
func (self *Block) SetRect(x1, y1, x2, y2 int) {
	self.Rectangle = image.Rect(x1, y1, x2, y2)
	self.Inner = image.Rect(x1, y1, x2, y2)
}

// GetRect implements the Drawable interface.
func (self *Block) GetRect() image.Rectangle {
	return self.Rectangle
}
