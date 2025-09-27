package filetree

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Walker struct {
	tree *FileTreeImpl
}

func NewWalker(tree *FileTreeImpl) *Walker {
	return &Walker{tree: tree}
}

func (w *Walker) GetVisibleNodes() []*TreeNode {
	var visible []*TreeNode
	if w.tree.root == nil {
		return visible
	}

	w.collectVisible(w.tree.root, &visible)
	return visible
}

func (w *Walker) collectVisible(node *TreeNode, visible *[]*TreeNode) {
	*visible = append(*visible, node)

	if node.Expanded && node.IsDir {
		for _, child := range node.Children {
			w.collectVisible(child, visible)
		}
	}
}

func (w *Walker) GetNextVisibleNode(current *TreeNode) *TreeNode {
	visible := w.GetVisibleNodes()

	for i, node := range visible {
		if node == current && i+1 < len(visible) {
			return visible[i+1]
		}
	}

	return nil
}

func (w *Walker) GetPrevVisibleNode(current *TreeNode) *TreeNode {
	visible := w.GetVisibleNodes()

	for i, node := range visible {
		if node == current && i > 0 {
			return visible[i-1]
		}
	}

	return nil
}

func (w *Walker) FindByName(pattern string) []*TreeNode {
	var results []*TreeNode
	pattern = strings.ToLower(pattern)

	if w.tree.root == nil {
		return results
	}

	w.searchRecursive(w.tree.root, func(node *TreeNode) bool {
		return strings.Contains(strings.ToLower(node.Name), pattern)
	}, &results)

	return results
}

func (w *Walker) FindByExtension(ext string) []*TreeNode {
	var results []*TreeNode

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	ext = strings.ToLower(ext)

	if w.tree.root == nil {
		return results
	}

	w.searchRecursive(w.tree.root, func(node *TreeNode) bool {
		if node.IsDir {
			return false
		}

		return strings.HasSuffix(strings.ToLower(node.Name), ext)
	}, &results)

	return results
}

func (w *Walker) FilterHidden(showHidden bool) []*TreeNode {
	var results []*TreeNode

	if w.tree.root == nil {
		return results
	}

	w.searchRecursive(w.tree.root, func(node *TreeNode) bool {
		isHidden := strings.HasPrefix(node.Name, ".")
		return showHidden || !isHidden
	}, &results)

	return results
}

func (w *Walker) FindByPattern(pattern string) []*TreeNode {
	var results []*TreeNode

	if w.tree.root == nil {
		return results
	}

	w.searchRecursive(w.tree.root, func(node *TreeNode) bool {
		matched, _ := filepath.Match(pattern, node.Name)
		return matched
	}, &results)

	return results
}

func (w *Walker) searchRecursive(node *TreeNode, predicate func(node *TreeNode) bool, results *[]*TreeNode) {
	if predicate(node) {
		*results = append(*results, node)
	}

	if node.IsDir && node.Loaded {
		for _, child := range node.Children {
			w.searchRecursive(child, predicate, results)
		}
	}
}

type WalkFunc func(node *TreeNode) error

func (w *Walker) Walk(fn WalkFunc) error {
	if w.tree.root == nil {
		return nil
	}

	return w.walkRecursive(w.tree.root, fn)
}

func (w *Walker) walkRecursive(node *TreeNode, fn WalkFunc) error {
	if err := fn(node); err != nil {
		return err
	}

	if node.IsDir && node.Loaded {
		for _, child := range node.Children {
			if err := w.walkRecursive(child, fn); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Walker) WalkFrom(startNode *TreeNode, fn WalkFunc) error {
	if startNode == nil {
		return fmt.Errorf("start node cannot be nil")
	}

	return w.walkRecursive(startNode, fn)
}

func (w *Walker) CollectWhere(predicate func(*TreeNode) bool) []*TreeNode {
	var results []*TreeNode

	w.Walk(func(node *TreeNode) error {
		if predicate(node) {
			results = append(results, node)
		}

		return nil
	})

	return results
}

func (w *Walker) CollectAll() []*TreeNode {
	return w.CollectWhere(func(node *TreeNode) bool {
		return true
	})
}
