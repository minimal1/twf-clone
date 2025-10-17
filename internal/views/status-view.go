package views

import (
	"fmt"

	"github.com/minimal1/twf-clone/internal/state"
	"github.com/minimal1/twf-clone/internal/terminal"
)

type StatusView struct{}

func (sv *StatusView) Render(term *terminal.Terminal, rect Rect, appState *state.AppState) error {
	currentNode := appState.Cursor().GetCurrentNode()
	if currentNode == nil {
		return nil
	}

	promptMsg := appState.View().GetPrompt()
	if promptMsg != "" {
		term.WriteColoredAt(rect.Y, rect.X, promptMsg, terminal.ColorBlue)
		return nil
	}

	path := currentNode.Path
	selectedCount := len(appState.Selection().GetSelectedNodes())

	leftText := fmt.Sprintf(" %s", path)
	term.WriteColoredAt(rect.Y, rect.X, leftText, terminal.ColorCyan)

	if selectedCount > 0 {
		rightText := fmt.Sprintf("Selected: %d", selectedCount)
		rightX := rect.X + rect.Width - len(rightText)
		term.WriteColoredAt(rect.Y, rightX, rightText, terminal.ColorYellow)
	}
	return nil
}

func (sv *StatusView) GetMinSize() (width, height int) {
	return 40, 1
}
