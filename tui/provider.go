package tui

import (
	"bytes"
	"errors"
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
)



func loadImages(providers []string) ([]image.Image, error) {
	var images []image.Image
	assestsDir := "_assets/tui"

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for _, provider := range providers {
		fileName := provider + ".png"
		img, err := ioutil.ReadFile(filepath.Join(dir, assestsDir, fileName))
		if err != nil {
			return nil, err
		}
		image, _, err := image.Decode(bytes.NewReader(img))
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

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

	upperSpacer := widgets.NewParagraph()
	upperSpacer.Text = "choose a provider"
	upperSpacer.Border = false
	grid.Set(
		ui.NewRow(1.0/6, spacer),
		ui.NewRow(1.0/6,
			ui.NewCol(3.0 / 8, spacer),
			ui.NewCol(2.0 / 8, upperSpacer),
			ui.NewCol(3.0 / 8, spacer),
			),
		ui.NewRow(1.0/3, col...),
		ui.NewRow(1.0/3, spacer),
	)

	return grid
}

type ImagesWidget struct {
	images    []*widgets.Image
	providers []string
	index     int
}

func imagesWidget(images []image.Image, providers []string) *ImagesWidget {
	var widgetImages []*widgets.Image
	for i, imgCol := range images {
		img := widgets.NewImage(imgCol)
		img.BorderStyle = ui.NewStyle(ui.ColorYellow)
		img.Border = false

		img.Title = providers[i]
		widgetImages = append(widgetImages, img)
	}
	iw := &ImagesWidget{images: widgetImages, providers: providers, index: 0}
	iw.images[iw.index].Border = true
	return iw

}

func (i *ImagesWidget) Next() {
	i.images[i.index].Border = false
	i.index = (i.index + 1) % len(i.images)
	i.images[i.index].Border = true
}

func (i *ImagesWidget) Previous() {
	i.images[i.index].Border = false
	i.index = (i.index + len(i.images) - 1) % len(i.images)
	i.images[i.index].Border = true
}

func ChooseProvider() (string, *widgets.Image, error) {
	providers := []string{"aws", "ibm", "google", "azure"}
	images, err := loadImages(providers)
	if err != nil {
		fmt.Println(err)
		return "", nil, err
	}

	widgetImages := imagesWidget(images, providers)
	grid := providerGrid(widgetImages.images)

	uiEvents := ui.PollEvents()
	for {
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

		ui.Render(grid)

	}
}
