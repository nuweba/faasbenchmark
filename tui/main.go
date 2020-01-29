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
	"sort"
)

func leftTestsMenu() *widgets.List {
	testsMenu := widgets.NewList()

	var ids []string
	for id := range testsuite.Tests.TestFunctions {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	for _, id := range ids {
		testsMenu.Rows = append(testsMenu.Rows, id)
	}

	testsMenu.Border = true
	testsMenu.SelectedRowStyle = ui.NewStyle(ui.ColorMagenta, ui.ColorClear, ui.ModifierBold)

	return testsMenu
}

func grid(plotTabs *widgets.TabPane, linePlot *widgets.Plot, leftTestView *widgets.List, logView *widgets.List, pImage *widgets.Image) *ui.Grid {
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
			ui.NewCol(5.0/6,
				ui.NewRow(1.0/10, plotTabs),
				ui.NewRow(9.0/10, linePlot),
			),
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
	gConfig, err := config.NewGlobalConfig(provider, arsenalPath, report, false)
	if err != nil {
		return nil, err
	}

	return gConfig, nil
}

func Tui(provider string, pImage *widgets.Image) {
	linePlot := lineChart("Invocation OverHead")
	linePlot2 := lineChart("Duration")
	linePlot3 := lineChart("Content Transfer")
	linePlot4 := lineChart("Reused")
	linePlot5 := lineChart("Fresh")
	leftTestsMenu := leftTestsMenu()
	logView := logView()

	tabpane := widgets.NewTabPane(linePlot.Title, linePlot2.Title, linePlot3.Title, linePlot4.Title, linePlot5.Title)
	tabpane.Border = false
	plots := map[string]*widgets.Plot{linePlot.Title: linePlot.Plot, linePlot2.Title: linePlot2.Plot, linePlot3.Title: linePlot3.Plot, linePlot4.Title: linePlot4.Plot, linePlot5.Title: linePlot5.Plot}

	localGrid := grid(tabpane, plots[tabpane.TabNames[tabpane.ActiveTabIndex]], leftTestsMenu, logView.List, pImage)

	tabMoveRight := func() {
		if tabpane.ActiveTabIndex < len(tabpane.TabNames)-1 {
			tabpane.ActiveTabIndex++

		} else {
			tabpane.ActiveTabIndex = 0
		}
		localGrid = grid(tabpane, plots[tabpane.TabNames[tabpane.ActiveTabIndex]], leftTestsMenu, logView.List, pImage)
	}

	result := &graphStream{
		ch: make(chan *graph.Result),
	}

	toggleList := NewToggleList(leftTestsMenu, logView.List)
	logs, err := hookStdout()

	if err != nil {
		fmt.Println(err)
		return
	}

	uiEvents := ui.PollEvents()
	for {
		ui.Render(localGrid)
		select {
		case data := <-logs.ch:
			logView.Modify(data)

		case resultData := <-result.ch:
			linePlot.Modify(resultData.InvocationOverHead)
			linePlot2.Modify(resultData.Duration)
			linePlot3.Modify(resultData.ContentTransfer)
			if resultData.Reused {
				linePlot4.Modify(resultData.InvocationOverHead)
			} else {
				linePlot5.Modify(resultData.InvocationOverHead)
			}

		case e := <-uiEvents:
			switch e.ID {
			case "<Tab>":
				tabMoveRight()
				ui.Clear()
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				localGrid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
			case "<Enter>":
				logView.Reset()
				linePlot.Reset()
				linePlot2.Reset()
				linePlot3.Reset()
				linePlot4.Reset()
				linePlot5.Reset()

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
