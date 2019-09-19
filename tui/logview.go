package tui

import (
	"bufio"
	"bytes"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type LogView struct {
	*widgets.List
	prevHadLine bool
}

func logView() *LogView {
	textList := widgets.NewList()

	textList.Rows = append(textList.Rows, "choose a test\nkeys:\nLeft and Right arrow keys to toggle tests and log view\nTab to toggle results\nQ to quit\nEnter to start test")
	textList.SelectedRowStyle = ui.NewStyle(ui.ColorYellow, ui.ColorClear, ui.ModifierBold)
	return &LogView{List: textList, prevHadLine: false}
}


func (l *LogView) Modify(data []byte) {
	d := bytes.NewReader(data)
	scan := bufio.NewScanner(d)
	for scan.Scan() {
		s := scan.Text()
		if l.prevHadLine == false && len(l.Rows) != 0 {
			l.Rows[len(l.Rows)-1] = l.Rows[len(l.Rows)-1] + s
		} else {
			l.Rows = append(l.Rows, s)
		}

		if l.SelectedRow == int(len(l.Rows) - 2) {
			l.ScrollBottom()
		}

		if s == string(data) {
			l.prevHadLine = false
		} else {
			l.prevHadLine = true
		}
	}
}

func (l *LogView) Reset() {
	l.Rows = []string{}
	l.ScrollTop()
}