package views

import (
	"strings"

	"github.com/minimal1/twf-clone/internal/filetree"
	"github.com/minimal1/twf-clone/internal/state"
	"github.com/minimal1/twf-clone/internal/terminal"
)

type TreeView struct {
	walker *filetree.Walker
}

func NewTreeView(walker *filetree.Walker) *TreeView {
	return &TreeView{
		walker: walker,
	}
}

func (tv *TreeView) Render(term *terminal.Terminal, rect Rect, appState *state.AppState) error {
	visibleNodes := tv.walker.GetVisibleNodes()

	scrollOffset := appState.View().GetScrollOffset()

	startIdx := scrollOffset
	endIdx := scrollOffset + rect.Height

	for i := startIdx; i < endIdx && i < len(visibleNodes); i++ {
		if i < 0 {
			continue
		}

		node := visibleNodes[i]
		y := rect.Y + (i - startIdx)

		indent := strings.Repeat("  ", node.Depth())

		color := terminal.ColorWhite
		if node == appState.Cursor().GetCurrentNode() {
			color = terminal.ColorYellow
		}
		if appState.Selection().IsSelected(node) {
			color = terminal.ColorGreen
		}

		text := indent + node.GetDisplayName()
		term.WriteColoredAt(y, rect.X, text, color)
	}

	return nil
}

func (tv *TreeView) GetMinSize() (width, height int) {
	return 20, 10
}
