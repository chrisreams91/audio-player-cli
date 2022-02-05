package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)


func buildBars(data [] float64, songName string) *widgets.BarChart{
	bc := widgets.NewBarChart()
	bc.SetRect(5, 5, 100, 25)

	bc.BarWidth = 1
	bc.BarGap = 0
	bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}

	bc.NumFormatter = func(f float64) string { return " "}
	bc.Title = songName
	bc.TitleStyle = ui.NewStyle(ui.ColorWhite)
	bc.Border = false

	// manage here or in wave processing?
	// bc.MaxVal = float64(10)
	
	return bc
}