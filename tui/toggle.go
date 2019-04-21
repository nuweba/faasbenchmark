package tui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type ToggleList struct {
	listA   *widgets.List
	listB   *widgets.List
	current *widgets.List
	currentRowStyle ui.Style
}

func NewToggleList(listA *widgets.List, listB *widgets.List) *ToggleList {
	t := &ToggleList{
		listA:   listA,
		listB:   listB,
		current: listA,
		currentRowStyle: listB.SelectedRowStyle,
	}
	t.listB.SelectedRowStyle = t.listB.TextStyle
	return t
}

func (t *ToggleList) ToggleTo(list *widgets.List) bool {
	if t.current == list {
		return false
	}
	list.SelectedRowStyle = t.currentRowStyle
	t.currentRowStyle = t.current.SelectedRowStyle
	t.current.SelectedRowStyle = t.current.TextStyle
	t.current = list

	return true
}

func (t *ToggleList) HandleEvent(e *ui.Event) bool {
	switch e.ID {
	case "<Right>":
		t.ToggleTo(t.listB)
	case "<Left>":
		t.ToggleTo(t.listA)
	case "j", "<Down>", "<MouseWheelDown>":
		t.current.ScrollDown()
	case "k", "<Up>", "<MouseWheelUp>":
		t.current.ScrollUp()
	case "<C-d>":
		t.current.ScrollHalfPageDown()
	case "<C-u>":
		t.current.ScrollHalfPageUp()
	case "<C-f>":
		t.current.ScrollPageDown()
	case "<C-b>":
		t.current.ScrollPageUp()
	case "<Home>":
		t.current.ScrollTop()
	case "G", "<End>":
		t.current.ScrollBottom()
	default:
		return false
	}

	return true

}
