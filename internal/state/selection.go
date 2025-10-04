package state

import (
	"slices"

	"github.com/minimal1/twf-clone/internal/filetree"
)

type ClipboardType int

const (
	ClipboardCopy ClipboardType = iota
	ClipboardCut
)

type SelectionState struct {
	selectedNodes []*filetree.TreeNode
	marks         map[string]*filetree.TreeNode
	clipboard     []*filetree.TreeNode
	clipboardType ClipboardType
}

func NewSelectionState() *SelectionState {
	return &SelectionState{
		selectedNodes: make([]*filetree.TreeNode, 0),
		marks:         make(map[string]*filetree.TreeNode),
		clipboard:     make([]*filetree.TreeNode, 0),
		clipboardType: ClipboardCopy,
	}
}

func (ss *SelectionState) IsSelected(node *filetree.TreeNode) bool {
	return slices.Contains(ss.selectedNodes, node)
}

func (ss *SelectionState) ToggleSelection(node *filetree.TreeNode) {
	index := slices.Index(ss.selectedNodes, node)

	if index == -1 {
		ss.selectedNodes = append(ss.selectedNodes, node)
	} else {
		ss.selectedNodes = append(ss.selectedNodes[:index], ss.selectedNodes[index+1:]...)
	}
}

func (ss *SelectionState) GetSelectedNodes() []*filetree.TreeNode {
	return ss.selectedNodes
}

func (ss *SelectionState) ClearSelection() {
	ss.selectedNodes = make([]*filetree.TreeNode, 0)
}

func (ss *SelectionState) SetMark(mark string, node *filetree.TreeNode) {
	ss.marks[mark] = node
}

func (ss *SelectionState) GetMark(mark string) *filetree.TreeNode {
	return ss.marks[mark]
}

func (ss *SelectionState) ClearMarks() {
	ss.marks = make(map[string]*filetree.TreeNode)
}

func (ss *SelectionState) Copy(nodes []*filetree.TreeNode) {
	ss.clipboardType = ClipboardCopy
	ss.clipboard = make([]*filetree.TreeNode, len(nodes))
	copy(ss.clipboard, nodes)
}

func (ss *SelectionState) Cut(nodes []*filetree.TreeNode) {
	ss.clipboardType = ClipboardCut
	ss.clipboard = make([]*filetree.TreeNode, len(nodes))
	copy(ss.clipboard, nodes)
}

func (ss *SelectionState) GetClipboard() []*filetree.TreeNode {
	return ss.clipboard
}

func (ss *SelectionState) ClearClipboard() {
	ss.clipboard = make([]*filetree.TreeNode, 0)
}
