package tui

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type LineChart struct {
	*widgets.Plot
}

func lineChart() *LineChart {
	lc := widgets.NewPlot()
	lc.Title = "braille-mode Line Chart"
	lc.Data = append(lc.Data, []float64{0, 0})
	lc.HorizontalScale = 3
	lc.AxesColor = ui.ColorWhite
	lc.LineColors[0] = ui.ColorYellow

	return &LineChart{Plot: lc}
}

func (l *LineChart) Modify(point float64) {
	if l.Data[0][0] == 0 && len(l.Data[0]) == 3 {
		l.Data[0] = l.Data[0][1:]
	}

	//if the dots came to the right end of the screen we will start moving the graph
	if len(l.Data[0]) == int((l.Max.X - l.Min.X) / 3) - 1{
		l.Data[0] = l.Data[0][1:]
	}
	l.Data[0] = append(l.Data[0], point)
	fmt.Println(point)
}

func (l *LineChart) Reset() {
	l.Data = [][]float64{}
	l.Data = append(l.Data, []float64{0, 0})
}
