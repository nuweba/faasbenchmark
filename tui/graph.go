package tui

import (
	"bufio"
	"bytes"
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"io"
	"math"
	"github.com/nuweba/faasbenchmark/cmd"
	"github.com/nuweba/faasbenchmark/config"
	"github.com/nuweba/faasbenchmark/provider"
	"github.com/nuweba/faasbenchmark/report/multi"
	"github.com/nuweba/faasbenchmark/report/output/file"
	"github.com/nuweba/faasbenchmark/report/output/graph"
	"github.com/nuweba/faasbenchmark/testsuite"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type stdch struct {
	*os.File
	ch chan []byte
}

func (s *stdch) Write(p []byte) (n int, err error) {
	s.ch <- p

	return len(p), nil
}

func New(f *os.File) *stdch {
	s := &stdch{f, make(chan []byte, 100)}
	go func() {
		io.Copy(s, f)
	}()
	return s
}

type graphStream struct {
	ch chan float64
}

func (gs *graphStream) parse(p []byte) float64 {
	fields := strings.Fields(string(p))

	f, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		fmt.Println(err)
	}
	return f

}

func (gs *graphStream) Write(p []byte) (n int, err error) {
	gs.ch <- gs.parse(p)
	return len(p), nil
}

func UI() {

	sinFloat64 := (func() []float64 {
		n := 400
		data := make([]float64, n)
		for i := range data {
			data[i] = 1 + math.Sin(float64(i)/5)
		}
		return data
	})()

	sl := widgets.NewSparkline()
	sl.Data = sinFloat64[:100]
	sl.LineColor = ui.ColorCyan
	sl.TitleStyle.Fg = ui.ColorWhite

	slg := widgets.NewSparklineGroup(sl)
	slg.Title = "Sparkline"

	lc := widgets.NewPlot()
	lc.Title = "braille-mode Line Chart"
	lc.Data = append(lc.Data, []float64{0, 0})
	//lc.Marker = widgets.MarkerDot
	lc.HorizontalScale = 3
	lc.AxesColor = ui.ColorWhite
	lc.LineColors[0] = ui.ColorYellow

	gs := make([]*widgets.Gauge, 3)
	for i := range gs {
		gs[i] = widgets.NewGauge()
		gs[i].Percent = i * 10
		gs[i].BarColor = ui.ColorRed
	}

	ls := widgets.NewList()
	ls.Rows = []string{
		"[1] Downloading File 1",
		"",
		"",
		"",
		"[2] Downloading File 2",
		"",
		"",
		"",
		"[3] Uploading File 3",
	}
	ls.Border = false

	testsMenu := testsMenu()

	outputl := widgets.NewList()
	ls.WrapText = true
	//outputl.BorderStyle = ui.NewStyle(ui.ColorWhite)
	outputl.Rows = append(outputl.Rows, "choose a test")
	outputl.SelectedRowStyle = ui.NewStyle(ui.ColorYellow, ui.ColorClear, ui.ModifierBold)

	p := widgets.NewParagraph()
	p.Text = "<> This row has 3 columns\n<- Widgets can be stacked up like left side\n<- Stacked widgets are treated as a single widget"
	p.Title = "Demonstration"
	p.SetRect(0, 0, 25, 8)

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	//grid.Set(
	//	ui.NewRow(1.0/3,
	//		ui.NewCol(1.0/1, testsMenu),
	//	),
	//	ui.NewRow(1.0/3,
	//		ui.NewCol(1.0/2, slg),
	//		ui.NewCol(1.0/2, lc),
	//	),
	//	ui.NewRow(1.0/3,
	//		ui.NewCol(1.0/4, ls),
	//		ui.NewCol(1.0/4,
	//			ui.NewRow(.9/3, gs[0]),
	//			ui.NewRow(.9/3, gs[1]),
	//			ui.NewRow(1.2/3, gs[2]),
	//		),
	//		ui.NewCol(1.0/2,
	//			ui.NewRow(1.0/1, outputl, )),
	//	),
	//)

	grid.Set(
		ui.NewRow(2.0/3,
			ui.NewCol(1.0/6, testsMenu),
			ui.NewCol(5.0/6, lc),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, outputl),
		),
	)

	pr, pw, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
	}
	output := New(pr)
	_ = output
	os.Stdout = pw
	ui.Render(grid)

	//tickerCount := 1
	uiEvents := ui.PollEvents()
	//ticker := time.NewTicker(time.Second).C
	provider, err := provider.NewProvider("aws")
	if err != nil {
		fmt.Println(err)
	}

	const (
		TestsDir = "arsenal"
	)

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileReport, err := file.New(dir)

	if err != nil {
		fmt.Println(err)
		return
	}
	gStream := &graphStream{
		ch: make(chan float64),
	}
	graphReport, err := graph.New(gStream)

	if err != nil {
		fmt.Println(err)
		return
	}

	report := multi.Report(fileReport, graphReport)

	arsenalPath := filepath.Join(dir, TestsDir)
	gConfig, err := config.NewGlobalConfig(provider, arsenalPath, report)
	if err != nil {
		fmt.Println(err)
	}

	currentList := testsMenu
	selectedRowStyle := outputl.SelectedRowStyle
	outputl.SelectedRowStyle = outputl.TextStyle
	prevHadLine := false
	for {
		select {
		case data := <-output.ch:
			d := bytes.NewReader(data)
			scan := bufio.NewScanner(d)
			for scan.Scan() {
				s := scan.Text()
				if prevHadLine == false && len(outputl.Rows) != 0 {
					outputl.Rows[len(outputl.Rows)-1] = outputl.Rows[len(outputl.Rows)-1] + s
				} else {
					outputl.Rows = append(outputl.Rows, s)
				}

				if outputl.SelectedRow == uint(len(outputl.Rows) - 2) {
					outputl.ScrollBottom()
				}

				ui.Render(grid)
				if s == string(data) {
					prevHadLine = false
				} else {
					prevHadLine = true
				}
			}



		case point := <-gStream.ch:
			if lc.Data[0][0] == 0 && len(lc.Data[0]) == 3 {
				lc.Data[0] = lc.Data[0][1:]
			}

		if len(lc.Data[0]) == 72 {
			lc.Data[0] = lc.Data[0][1:]
		}
			lc.Data[0] = append(lc.Data[0], point)
			fmt.Println(point)
			ui.Render(grid)

		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Right>":
				if currentList == outputl {
					continue
				}
				outputl.SelectedRowStyle = selectedRowStyle
				selectedRowStyle = testsMenu.SelectedRowStyle
				testsMenu.SelectedRowStyle = testsMenu.TextStyle
				currentList = outputl

				ui.Render(grid)

			case "<Left>":
				if currentList == testsMenu {
					continue
				}
				testsMenu.SelectedRowStyle = selectedRowStyle
				selectedRowStyle = outputl.SelectedRowStyle
				outputl.SelectedRowStyle = outputl.TextStyle
				currentList = testsMenu
				ui.Render(grid)

			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			case "<Enter>":
				outputl.Rows = []string{}
				ui.Render(grid)
				testId := testsMenu.Rows[testsMenu.SelectedRow]
				go func() {
					err := cmd.RunSpecificTests(gConfig, testId)
					if err != nil {
						fmt.Println(err)
					}
				}()

			default:
				if scroll(currentList, &e) {
					ui.Render(grid)
				}

			}
		}
	}

}

func testsMenu() *widgets.List {
	testsMenu := widgets.NewList()

	for id := range testsuite.Tests.TestFunctions {
		testsMenu.Rows = append(testsMenu.Rows, id)
	}

	testsMenu.Border = true
	testsMenu.SelectedRowStyle = ui.NewStyle(ui.ColorMagenta, ui.ColorClear, ui.ModifierBold)
	//testsMenu.BorderStyle = ui.NewStyle(ui.ColorMagenta)

	return testsMenu
}

func scroll(ls *widgets.List, e *ui.Event) bool {
	switch e.ID {
	case "j", "<Down>", "<MouseWheelDown>":
		ls.ScrollDown()
	case "k", "<Up>", "<MouseWheelUp>":
		ls.ScrollUp()
	case "<C-d>":
		ls.HalfPageDown()
	case "<C-u>":
		ls.HalfPageUp()
	case "<C-f>":
		ls.PageDown()
	case "<C-b>":
		ls.PageUp()
	case "<Home>":
		ls.ScrollTop()
	case "G", "<End>":
		ls.ScrollBottom()
	default:
		return false
	}

	return true
}
