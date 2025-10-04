package state

import (
	"time"

	"github.com/minimal1/twf-clone/internal/filetree"
)

type Position struct {
	Row int
	Col int
}

type Navigation struct {
	Node      *filetree.TreeNode
	Position  Position
	Timestamp time.Time
}

type CursorState struct {
	currentNode *filetree.TreeNode
	position    Position
	history     []Navigation
	maxHistory  int
}

func NewCursorState() *CursorState {
	return &CursorState{
		position: Position{
			Row: 2,
			Col: 2,
		},
		history:    make([]Navigation, 0),
		maxHistory: 50,
	}
}

func (cs *CursorState) GetCurrentNode() *filetree.TreeNode {
	return cs.currentNode
}

func (cs *CursorState) SetCurrentNode(node *filetree.TreeNode) {
	cs.currentNode = node
}

func (cs *CursorState) GetPosition() Position {
	return cs.position
}

func (cs *CursorState) SetPosition(pos Position) {
	cs.position = pos
}

func (cs *CursorState) MoveTo(node *filetree.TreeNode, pos Position) {
	if cs.currentNode != nil {
		cs.addToHistory(cs.currentNode, cs.position)
	}

	cs.SetCurrentNode(node)
	cs.SetPosition(pos)
}

func (cs *CursorState) CanGoBack() bool {
	return len(cs.history) > 0
}

func (cs *CursorState) GoBack() (*filetree.TreeNode, Position, bool) {
	if !cs.CanGoBack() {
		return nil, Position{}, false
	}

	lastIndex := len(cs.history) - 1
	nav := cs.history[lastIndex]

	cs.history = cs.history[:lastIndex]

	cs.currentNode = nav.Node
	cs.position = nav.Position

	return nav.Node, nav.Position, true
}

func (cs *CursorState) addToHistory(node *filetree.TreeNode, pos Position) {
	nav := Navigation{
		Node:      node,
		Position:  pos,
		Timestamp: time.Now(),
	}

	cs.history = append(cs.history, nav)

	if len(cs.history) > cs.maxHistory {
		cs.history = cs.history[1:]
	}
}
