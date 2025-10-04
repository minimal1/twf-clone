package state

import "github.com/minimal1/twf-clone/internal/filetree"

type AppState struct {
	cursor    *CursorState
	selection *SelectionState
	view      *ViewState
	config    *ConfigState
}

func NewAppState() *AppState {
	return &AppState{
		cursor:    NewCursorState(),
		selection: NewSelectionState(),
		view:      NewViewState(),
		config:    NewConfigState(),
	}
}

func (as *AppState) Cursor() *CursorState {
	return as.cursor
}

func (as *AppState) Selection() *SelectionState {
	return as.selection
}

func (as *AppState) View() *ViewState {
	return as.view
}

func (as *AppState) Config() *ConfigState {
	return as.config
}

func (as *AppState) Initialize(rootNode *filetree.TreeNode) error {
	if rootNode != nil {
		as.cursor.SetCurrentNode(rootNode)
	}

	return nil
}
