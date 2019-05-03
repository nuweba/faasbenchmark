package tui

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/nuweba/faasbenchmark/cmd"
	"github.com/nuweba/faasbenchmark/config"
	"github.com/nuweba/faasbenchmark/provider"
	"github.com/nuweba/faasbenchmark/report/multi"
	"github.com/nuweba/faasbenchmark/report/output/graph"
	"github.com/nuweba/faasbenchmark/report/output/json"
	"github.com/nuweba/faasbenchmark/testsuite"
	"os"
	"path/filepath"
)

func leftTestsMenu() *widgets.List {
	testsMenu := widgets.NewList()

	for id := range testsuite.Tests.TestFunctions {
		testsMenu.Rows = append(testsMenu.Rows, id)
	}

	testsMenu.Border = true
	testsMenu.SelectedRowStyle = ui.NewStyle(ui.ColorMagenta, ui.ColorClear, ui.ModifierBold)

	return testsMenu
}

func grid(lineChart *widgets.Plot, leftTestView *widgets.List, logView *widgets.List, pImage *widgets.Image) *ui.Grid {
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	pImage.Border = false
	grid.Set(
		ui.NewRow(2.0/3,
			ui.NewCol(1.0/6,
				ui.NewRow(1.0/2, pImage),
				ui.NewRow(1.0/2, leftTestView),
			),
			ui.NewCol(5.0/6, lineChart),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(1.0/1, logView),
		),
	)

	return grid
}

func faasTestConfig(providerName string, resultCh *graphStream) (*config.Global, error) {
	TestsDir := "arsenal"

	provider, err := provider.NewProvider(providerName)
	if err != nil {
		return nil, err
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fileReport, err := json.New(dir)
	if err != nil {
		return nil, err
	}

	graphReport, err := graph.New(resultCh)
	if err != nil {
		return nil, err
	}

	report := multi.Report(fileReport, graphReport)

	arsenalPath := filepath.Join(dir, TestsDir)
	gConfig, err := config.NewGlobalConfig(provider, arsenalPath, report)
	if err != nil {
		return nil, err
	}

	return gConfig, nil
}

func Tui(provider string, pImage *widgets.Image) {
	lineChart := lineChart()
	leftTestsMenu := leftTestsMenu()
	logView := logView()

	grid := grid(lineChart.Plot, leftTestsMenu, logView.List, pImage)

	result := &graphStream{
		ch: make(chan float64),
	}

	toggleList := NewToggleList(leftTestsMenu, logView.List)
	logs, err := hookStdout()

	if err != nil {
		fmt.Println(err)
		return
	}

	uiEvents := ui.PollEvents()
	for {
		ui.Render(grid)
		select {
		case data := <-logs.ch:
			logView.Modify(data)

		case point := <-result.ch:
			lineChart.Modify(point)

		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
			case "<Enter>":
				logView.Reset()
				lineChart.Reset()

				testId := leftTestsMenu.Rows[leftTestsMenu.SelectedRow]
				go func() {
					faasTestConfig, err := faasTestConfig(provider, result)
					if err != nil {
						fmt.Println(err)
						return
					}

					err = cmd.RunSpecificTests(faasTestConfig, testId)
					if err != nil {
						fmt.Println(err)
					}
				}()
			default:
				toggleList.HandleEvent(&e)
			}
		}
	}

}
