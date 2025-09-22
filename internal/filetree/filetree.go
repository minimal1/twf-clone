package filetree

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileTree interface {
	LoadRoot(path string) error
	GetRoot() *TreeNode
	GetCurrentNode() *TreeNode
	SetCurrentNode(node *TreeNode) error
	ExpandNode(node *TreeNode) error
	CollapseNode(node *TreeNode) error
	RefreshNode(node *TreeNode) error
}

type FileTreeImpl struct {
	root        *TreeNode
	currentNode *TreeNode
}

func NewFileTree() *FileTreeImpl {
	return &FileTreeImpl{}
}

func (ft *FileTreeImpl) LoadRoot(path string) error {
	node, err := NewTreeNode(path)

	if err != nil {
		return err
	}

	ft.root = node
	ft.currentNode = node

	return nil
}

func (ft *FileTreeImpl) GetRoot() *TreeNode {
	return ft.root
}

func (ft *FileTreeImpl) GetCurrentNode() *TreeNode {
	return ft.currentNode
}

func (ft *FileTreeImpl) SetCurrentNode(node *TreeNode) error {
	if node == nil {
		return fmt.Errorf("node cannot be nil")
	}

	ft.currentNode = node
	return nil
}

func (ft *FileTreeImpl) ExpandNode(node *TreeNode) error {
	if !node.CanExpand() {
		return fmt.Errorf("failed to expand node, caused by this node can't expand")
	}

	if node.Expanded {
		return nil
	}

	if node.Loaded {
		node.Expanded = true
		return nil
	}

	if err := ft.loadChildren(node); err != nil {
		return err
	}
	node.Expanded = true

	return nil
}

func (ft *FileTreeImpl) loadChildren(node *TreeNode) error {
	entries, err := os.ReadDir(node.Path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", node.Path, err)
	}

	for _, entry := range entries {
		childPath := filepath.Join(node.Path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		child := NewTreeNodeFromInfo(childPath, info)
		node.AddChild(child)
	}

	node.Loaded = true
	return nil
}

func (ft *FileTreeImpl) CollapseNode(node *TreeNode) error {
	if !node.CanExpand() {
		return fmt.Errorf("failed to collapse node, caused by this node can't collapse")
	}

	node.Expanded = false
	return nil
}

func (ft *FileTreeImpl) RefreshNode(node *TreeNode) error {
	if !node.IsDir {
		return nil
	}

	node.Children = []*TreeNode{}
	node.Loaded = false
	node.Expanded = false

	return ft.ExpandNode(node)
}
