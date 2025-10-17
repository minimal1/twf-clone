package main

import (
	"fmt"
	"os"

	"github.com/minimal1/twf-clone/internal/filetree"
	"github.com/minimal1/twf-clone/internal/state"
	"github.com/minimal1/twf-clone/internal/terminal"
	"github.com/minimal1/twf-clone/internal/views"
)

type App struct {
	term     *terminal.Terminal
	filetree *filetree.FileTreeImpl
	walker   *filetree.Walker
	appState *state.AppState
	running  bool
}

func NewApp(startPath string) (*App, error) {
	term, termErr := terminal.NewTerminal()
	if termErr != nil {
		return nil, termErr
	}

	ft := filetree.NewFileTree()
	ftErr := ft.LoadRoot(startPath)
	if ftErr != nil {
		term.Cleanup()
		return nil, ftErr
	}

	walker := filetree.NewWalker(ft)

	appState := state.NewAppState()
	appState.Initialize(ft.GetRoot())

	return &App{
		term:     term,
		filetree: ft,
		walker:   walker,
		appState: appState,
		running:  false,
	}, nil
}

func (app *App) Cleanup() {
	if app.term != nil {
		app.term.Cleanup()
	}
}

func main() {
	startPath := "."
	if len(os.Args) > 1 {
		startPath = os.Args[1]
	}

	app, err := NewApp(startPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize app: %v\n", err)
		os.Exit(1)
	}

	defer app.Cleanup()

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %V\n", err)
		os.Exit(1)
	}
}

func (app *App) Run() error {
	if err := app.term.EnableRawMode(); err != nil {
		return err
	}

	defer app.term.DisableRawMode()

	app.term.EnterAltScreen()
	defer app.term.ExitAltScreen()

	app.term.ClearScreen()
	app.term.HideCursor()
	defer app.term.ShowCursor()

	width, height, _ := app.term.GetSize()

	treeView := views.NewTreeView(app.walker)
	statusView := views.StatusView{}
	layout := views.NewLayout(treeView, &statusView)
	layout.SetSize(width, height)

	if err := layout.Render(app.term, app.appState); err != nil {
		return err
	}

	app.running = true

	for app.running {
		event, err := app.term.ReadEvent()
		if err != nil {
			return err
		}

		app.handleEvent(event)

		app.adjustScroll(height)

		app.term.ClearScreen()
		if err := layout.Render(app.term, app.appState); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) handleEvent(event terminal.Event) {
	switch e := event.(type) {
	case terminal.KeyPressEvent:
		app.handleKeyPress(e)
	}
}

func (app *App) handleKeyPress(event terminal.KeyPressEvent) {
	if event.Rune != 0 {
		app.handleRuneKey(event.Rune)
		return
	}

	switch event.Key {
	case terminal.KeyEsc, terminal.KeyCtrlC:
		// 입력 모드 취소
		viewState := app.appState.View()
		if viewState.IsWaitingForInput() {
			viewState.SetInputMode(state.InputModeNormal)
			viewState.ClearPrompt()
			return
		}

		app.running = false
	case terminal.KeyArrowDown:
		app.moveDown()
	case terminal.KeyArrowUp:
		app.moveUp()
	case terminal.KeyArrowRight, terminal.KeyEnter:
		app.expandOrEnter()
	case terminal.KeyArrowLeft:
		app.collapseOrParent()
	}
}

func (app *App) handleRuneKey(r rune) {
	viewState := app.appState.View()

	// 입력 모드 처리
	switch viewState.GetInputMode() {
	case state.InputModeWaitingForMark:
		app.setBookmark(string(r))
		viewState.SetInputMode(state.InputModeNormal)
		viewState.ClearPrompt()
		return
	case state.InputModeWaitingForJump:
		app.jumpToBookmark(string(r))
		viewState.SetInputMode(state.InputModeNormal)
		viewState.ClearPrompt()
		return
	}

	// 일반 키 처리
	switch r {
	case 'q':
		app.running = false
	case 'j':
		app.moveDown()
	case 'k':
		app.moveUp()
	case 'l':
		app.expandOrEnter()
	case 'h':
		app.collapseOrParent()
	case ' ':
		app.toggleSelection()
	case 'm':
		viewState.SetInputMode(state.InputModeWaitingForMark)
		viewState.SetPrompt(" Mark: _")
	case '\'':
		viewState.SetInputMode(state.InputModeWaitingForJump)
		viewState.SetPrompt(" Jump to: _")
	}
}

func (app *App) moveDown() {
	currentNode := app.appState.Cursor().GetCurrentNode()
	nextNode := app.walker.GetNextVisibleNode(currentNode)

	if nextNode != nil {
		app.appState.Cursor().SetCurrentNode(nextNode)
	}
}

func (app *App) moveUp() {
	currentNode := app.appState.Cursor().GetCurrentNode()
	prevNode := app.walker.GetPrevVisibleNode(currentNode)

	if prevNode != nil {
		app.appState.Cursor().SetCurrentNode(prevNode)
	}
}

func (app *App) expandOrEnter() {
	currentNode := app.appState.Cursor().GetCurrentNode()

	if currentNode != nil && currentNode.IsDir {
		app.filetree.ExpandNode(currentNode)
	}
}

func (app *App) collapseOrParent() {
	currentNode := app.appState.Cursor().GetCurrentNode()

	if currentNode == nil {
		return
	}

	if currentNode.IsDir && currentNode.Expanded {
		app.filetree.CollapseNode(currentNode)
	} else if currentNode.Parent != nil {
		app.appState.Cursor().SetCurrentNode(currentNode.Parent)
	}
}

func (app *App) toggleSelection() {
	currentNode := app.appState.Cursor().GetCurrentNode()
	if currentNode != nil {
		app.appState.Selection().ToggleSelection(currentNode)
	}
}

func (app *App) adjustScroll(screenHeight int) {
	currentNode := app.appState.Cursor().GetCurrentNode()
	visibleNodes := app.walker.GetVisibleNodes()

	currentIndex := -1
	for i, node := range visibleNodes {
		if node == currentNode {
			currentIndex = i
			break
		}
	}

	treeHeight := screenHeight - 1
	scrollOffset := app.appState.View().GetScrollOffset()

	if currentIndex < scrollOffset {
		app.appState.View().SetScrollOffset(currentIndex)
	}
	if currentIndex >= scrollOffset+treeHeight {
		app.appState.View().SetScrollOffset(currentIndex - treeHeight + 1)
	}

	maxScroll := max(len(visibleNodes)-treeHeight, 0)

	if scrollOffset > maxScroll {
		app.appState.View().SetScrollOffset(maxScroll)
	}
}

func (app *App) setBookmark(mark string) {
	currentNode := app.appState.Cursor().GetCurrentNode()
	if currentNode != nil {
		app.appState.Selection().SetMark(mark, currentNode)
	}
}

func (app *App) jumpToBookmark(mark string) {
	node := app.appState.Selection().GetMark(mark)
	if node != nil {
		app.appState.Cursor().SetCurrentNode(node)
	}
}
