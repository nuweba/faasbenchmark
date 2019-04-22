package tui

import (
	"bytes"
	"github.com/disintegration/gift"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
)


const assestsDir = "_assets/tui"
const imageRatio = 1.8
const imageSize = 60
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
	//fmt.Println(i.images[i.index].Max.X - i.images[i.index].Min.X, i.images[i.index].Max.Y - i.images[i.index].Min.Y)
	i.images[i.index].Border = false
	i.index = (i.index + 1) % len(i.images)
	i.images[i.index].Border = true
}

func (i *ImagesWidget) Previous() {
	i.images[i.index].Border = false
	i.index = (i.index + len(i.images) - 1) % len(i.images)
	i.images[i.index].Border = true
}

func loadImages(providers []string) ([]image.Image, error) {
	var images []image.Image

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
		image1, _, err := image.Decode(bytes.NewReader(img))
		if err != nil {
			return nil, err
		}


		// 1. Create a new filter list and add some filters.
		g := gift.New(
			gift.Contrast(30),
			gift.Pixelate(3),

			gift.ResizeToFit(imageSize * imageRatio, imageSize, gift.LanczosResampling),

		)

		// 2. Create a new image of the corresponding size.
		// dst is a new target image, src is the original image.
		dst := image.NewRGBA(g.Bounds(image1.Bounds()))

		// 3. Use the Draw func to apply the filters to src and store the result in dst.
		g.Draw(dst, image1)
		images = append(images, dst)
	}

	return images, nil
}
