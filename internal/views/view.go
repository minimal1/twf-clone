package views

import (
	"github.com/minimal1/twf-clone/internal/state"
	"github.com/minimal1/twf-clone/internal/terminal"
)

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

type View interface {
	Render(term *terminal.Terminal, rect Rect, appState *state.AppState) error
	GetMinSize() (width, height int)
}
