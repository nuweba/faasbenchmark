package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/nuweba/faasbenchmark/tui"
	"os"
)

func main() {
	if err := ui.Init(); err != nil {
		fmt.Printf("failed to initialize termui: %v", err)
		os.Exit(1)
	}
	defer ui.Close()

	provider,pImage, err := tui.ChooseProvider()
	if err != nil {
		fmt.Println(err)
		return
	}
	tui.Tui(provider, pImage)
}
