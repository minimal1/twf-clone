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
