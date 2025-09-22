package filetree

import (
	"fmt"
	"os"
	"time"
)

type TreeNode struct {
	Path    string
	Name    string
	IsDir   bool
	Size    int64
	ModTime time.Time

	Parent   *TreeNode
	Children []*TreeNode

	Expanded bool
	Loaded   bool

	Selected bool
}

func NewTreeNode(path string) (*TreeNode, error) {
	info, err := os.Stat(path)

	if err != nil {
		return nil, fmt.Errorf("failded to stat %s: %w", path, err)
	}

	return &TreeNode{
		Path:    path,
		Name:    info.Name(),
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime(),

		Children: []*TreeNode{},
	}, nil
}

func NewTreeNodeFromInfo(path string, info os.FileInfo) *TreeNode {
	return &TreeNode{
		Path:    path,
		Name:    info.Name(),
		IsDir:   info.IsDir(),
		Size:    info.Size(),
		ModTime: info.ModTime(),

		Children: []*TreeNode{},
	}
}

func (n *TreeNode) AddChild(child *TreeNode) {
	n.Children = append(n.Children, child)
	child.Parent = n
}

func (n *TreeNode) RemoveChild(child *TreeNode) bool {
	for i, v := range n.Children {
		if v == child {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			return true
		}
	}

	return false
}

func (n *TreeNode) GetChildByName(name string) *TreeNode {
	for _, v := range n.Children {
		if v.Name == name {
			return v
		}
	}

	return nil
}

func (n *TreeNode) IsRoot() bool {
	return n.Parent == nil
}

func (n *TreeNode) Depth() int {
	depth := 0
	current := n

	for current.Parent != nil {
		depth++
		current = current.Parent
	}

	return depth
}

func (n *TreeNode) CanExpand() bool {
	return n.IsDir
}

func (n *TreeNode) String() string {
	filetype := 'F'
	if n.IsDir {
		filetype = 'D'
	}

	return fmt.Sprintf("[%c] %s", filetype, n.Path)
}

func (n *TreeNode) GetDisplayName() string {
	if !n.IsDir {
		return n.Name
	}

	if n.Expanded {
		return "▼ " + n.Name
	}

	return "▶ " + n.Name
}
