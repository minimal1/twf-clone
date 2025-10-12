package views

import (
	"github.com/minimal1/twf-clone/internal/state"
	"github.com/minimal1/twf-clone/internal/terminal"
)

type Layout struct {
	treeView   *TreeView
	statusView *StatusView
	termWidth  int
	termHeight int
}

func NewLayout(treeView *TreeView, statusView *StatusView) *Layout {
	return &Layout{
		treeView:   treeView,
		statusView: statusView,
		termWidth:  80,
		termHeight: 24,
	}
}

func (l *Layout) SetSize(width, height int) {
	l.termWidth = width
	l.termHeight = height
}

func (l *Layout) Render(term *terminal.Terminal, appState *state.AppState) error {
	statusRect := Rect{
		X:      1,
		Y:      l.termHeight,
		Width:  l.termWidth,
		Height: 1,
	}

	treeRect := Rect{
		X:      1,
		Y:      1,
		Width:  l.termWidth,
		Height: l.termHeight - 1,
	}

	if err := l.treeView.Render(term, treeRect, appState); err != nil {
		return err
	}

	if err := l.statusView.Render(term, statusRect, appState); err != nil {
		return err
	}

	return nil
}
