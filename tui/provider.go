package tui

import (
	"errors"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/nuweba/faasbenchmark/provider"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func providerGrid(widgetImages []*widgets.Image) *ui.Grid {
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	spacer := widgets.NewParagraph()
	spacer.Border = false

	slots := float64(len(widgetImages) + 2)

	var col []interface{}

	col = append(col, ui.NewCol(1.0/slots, spacer))
	for i := range widgetImages {
		col = append(col, ui.NewCol(1.0/slots, widgetImages[i]))
	}

	col = append(col, ui.NewCol(1.0/slots, spacer))

	upperSpacer := widgets.NewTable()
	upperSpacer.Rows = [][]string{
		{"choose a provider"},
	}

	upperSpacer.TextAlignment = ui.AlignCenter
	upperSpacer.Border = false
	grid.Set(
		ui.NewRow(1.0/6, spacer),
		ui.NewRow(1.0/6,
			ui.NewCol(1.0 / 3, spacer),
			ui.NewCol(1.0 / 3, upperSpacer),
			ui.NewCol(1.0 / 3, spacer),
			),
		ui.NewRow(1.0/3, col...),
		ui.NewRow(1.0/3, spacer),
	)

	return grid
}



func ChooseProvider() (string, *widgets.Image, error) {
	//todo: for testing
	//providers := []string{"aws", "ibm", "google", "azure"}
	providers := provider.List()
	images, err := loadImages(providers)
	if err != nil {
		fmt.Println(err)
		return "", nil, err
	}

	widgetImages := imagesWidget(images, providers)
	grid := providerGrid(widgetImages.images)

	uiEvents := ui.PollEvents()
	for {
		ui.Render(grid)
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return "", nil, errors.New("provider was not chosen")
		case "<Left>", "h":
			widgetImages.Previous()
		case "<Right>", "l":
			widgetImages.Next()
		case "<Enter>":
			ui.Clear()
			return widgetImages.providers[widgetImages.index], widgetImages.images[widgetImages.index], nil
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			grid.SetRect(0, 0, payload.Width, payload.Height)
			ui.Clear()
		}
	}
}
