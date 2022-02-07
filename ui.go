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

	return bc
}

func buildBottomBars(data[] float64) *BarChart{ 
	bc := NewBarChart(true)
	bc.SetRect(5, 15, 100, 30)
	
	return bc
}

// Everything below is slightly modified copy pasted from the lib so i didnt have to fork
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
		BarColors:    []ui.Color{ui.ColorRed, ui.ColorGreen},
		NumStyles:    ui.Theme.BarChart.Nums,
		LabelStyles:  ui.Theme.BarChart.Labels,
		NumFormatter: func(n float64) string { return fmt.Sprint(n) },
		BarGap:       0,
		BarWidth:     1,
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
					// assign this positional 15 to some variable
					buf.SetCell(c, image.Pt(x, self.Max.Y - y + 15))
				} else {
					buf.SetCell(c, image.Pt(x, y))
				}
			}
		}

		barXCoordinate += (self.BarWidth + self.BarGap)
	}
}

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
