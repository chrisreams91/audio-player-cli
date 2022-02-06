package main

import (
	"fmt"
	"image"
	"sync"

	ui "github.com/gizak/termui/v3"

	rw "github.com/mattn/go-runewidth"
)


func buildTopBars(data[] float64, songName string) *BarChart{
	bc := NewBarChart(false)
	bc.SetRect(5, 0, 100, 18)

	bc.BarWidth = 1
	bc.BarGap = 0
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}

	bc.NumFormatter = func(f float64) string { return " "}
	bc.Title = songName
	bc.TitleStyle = ui.NewStyle(ui.ColorWhite)
	bc.Border = false

	bc.PaddingTop = 0
	bc.PaddingBottom = 0
	bc.PaddingLeft = 0
	bc.PaddingRight = 0
	
	return bc
}

func buildBottomBars(data[] float64) *BarChart{ 
	bc := NewBarChart(true)
	bc.SetRect(5, 0, 100, 30)

	bc.BarWidth = 1
	bc.BarGap = 0
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}

	bc.NumFormatter = func(f float64) string { return " "}
	bc.TitleStyle = ui.NewStyle(ui.ColorWhite)
	bc.Border = false


	bc.PaddingTop = 0
	bc.PaddingBottom = 0
	bc.PaddingLeft = 0
	bc.PaddingRight = 0

	return bc
}

// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

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
			for y := self.Inner.Max.Y ; y > (self.Inner.Max.Y-2)-height; y-- {
				c := ui.NewCell(' ', ui.NewStyle(ui.ColorClear, ui.SelectColor(self.BarColors, i)))
				buf.SetCell(c, image.Pt(x, y))
			}
		}

		// draw label
		if i < len(self.Labels) {
			labelXCoordinate := barXCoordinate +
				int((float64(self.BarWidth) / 2)) -
				int((float64(rw.StringWidth(self.Labels[i])) / 2))
			buf.SetString(
				self.Labels[i],
				ui.SelectStyle(self.LabelStyles, i),
				image.Pt(labelXCoordinate, self.Inner.Max.Y-1),
			)
		}

		// draw number
		numberXCoordinate := barXCoordinate + int((float64(self.BarWidth) / 2))
		if numberXCoordinate <= self.Inner.Max.X {
			buf.SetString(
				self.NumFormatter(data),
				ui.NewStyle(
					ui.SelectStyle(self.NumStyles, i+1).Fg,
					ui.SelectColor(self.BarColors, i),
					ui.SelectStyle(self.NumStyles, i+1).Modifier,
				),
				image.Pt(numberXCoordinate, self.Inner.Max.Y-2),
			)
		}

		barXCoordinate += (self.BarWidth + self.BarGap)
	}
}

// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.




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

func (self *Block) drawBorder(buf *ui.Buffer) {
	verticalCell := ui.Cell{ui.VERTICAL_LINE, self.BorderStyle}
	horizontalCell := ui.Cell{ui.HORIZONTAL_LINE, self.BorderStyle}

	// draw lines
	if self.BorderTop {
		buf.Fill(horizontalCell, image.Rect(self.Min.X, self.Min.Y, self.Max.X, self.Min.Y+1))
	}
	if self.BorderBottom {
		buf.Fill(horizontalCell, image.Rect(self.Min.X, self.Max.Y-1, self.Max.X, self.Max.Y))
	}
	if self.BorderLeft {
		buf.Fill(verticalCell, image.Rect(self.Min.X, self.Min.Y, self.Min.X+1, self.Max.Y))
	}
	if self.BorderRight {
		buf.Fill(verticalCell, image.Rect(self.Max.X-1, self.Min.Y, self.Max.X, self.Max.Y))
	}

	// draw corners
	if self.BorderTop && self.BorderLeft {
		buf.SetCell(ui.Cell{ui.TOP_LEFT, self.BorderStyle}, self.Min)
	}
	if self.BorderTop && self.BorderRight {
		buf.SetCell(ui.Cell{ui.TOP_RIGHT, self.BorderStyle}, image.Pt(self.Max.X-1, self.Min.Y))
	}
	if self.BorderBottom && self.BorderLeft {
		buf.SetCell(ui.Cell{ui.BOTTOM_LEFT, self.BorderStyle}, image.Pt(self.Min.X, self.Max.Y-1))
	}
	if self.BorderBottom && self.BorderRight {
		buf.SetCell(ui.Cell{ui.BOTTOM_RIGHT, self.BorderStyle}, self.Max.Sub(image.Pt(1, 1)))
	}
}

// Draw implements the Drawable interface.
func (self *Block) Draw(buf *ui.Buffer) {
	if self.Border {
		self.drawBorder(buf)
	}
	buf.SetString(
		self.Title,
		self.TitleStyle,
		image.Pt(self.Min.X+2, self.Min.Y),
	)
}

// SetRect implements the Drawable interface.
func (self *Block) SetRect(x1, y1, x2, y2 int) {
	self.Rectangle = image.Rect(x1, y1, x2, y2)
	// self.Inner = image.Rect(
	// 	self.Min.X+1+self.PaddingLeft,
	// 	self.Min.Y+1+self.PaddingTop,
	// 	self.Max.X-1-self.PaddingRight,
	// 	self.Max.Y-1-self.PaddingBottom,
	// )
	self.Inner = image.Rect(x1, y1, x2, y2)
}

// GetRect implements the Drawable interface.
func (self *Block) GetRect() image.Rectangle {
	return self.Rectangle
}
