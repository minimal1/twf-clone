package filetree

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
