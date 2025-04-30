package redblacktrees

import (
	"cmp"
	"iter"
)

type color bool

const (
	red   color = true
	black color = false
)

// Node represents a node in a red-black tree.
type Node[K cmp.Ordered, V any] struct {
	key   K
	value V

	color  color
	left   *Node[K, V]
	right  *Node[K, V]
	parent *Node[K, V]
	size   int
}

// Key returns the key of the node.
func (n *Node[K, V]) Key() K {
	return n.key
}

// Value returns the value of the node.
func (n *Node[K, V]) Value() V {
	return n.value
}

// Tree represents a red-black tree.
type Tree[K cmp.Ordered, V any] struct {
	Root *Node[K, V]
}

// New returns a new empty red-black tree.
func New[K cmp.Ordered, V any]() *Tree[K, V] {
	return &Tree[K, V]{}
}

// Len returns the number of nodes in the tree.
func Len[K cmp.Ordered, V any](t *Tree[K, V]) int {
	if t.Root == nil {
		return 0
	}
	return t.Root.size
}

// Clear removes all nodes from the tree.
func Clear[K cmp.Ordered, V any](t *Tree[K, V]) {
	t.Root = nil
}

// Insert adds a key-value pair into the tree or updates the value if the key already exists.
func Insert[K cmp.Ordered, V any](t *Tree[K, V], key K, value V) bool {
	newNode := &Node[K, V]{key: key, value: value, color: red, size: 1}
	var parent *Node[K, V]
	curr := t.Root

	for curr != nil {
		parent = curr
		if key < curr.key {
			curr = curr.left
		} else if key > curr.key {
			curr = curr.right
		} else {
			curr.value = value
			return false
		}
	}

	newNode.parent = parent
	if parent == nil {
		t.Root = newNode
	} else if key < parent.key {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	for p := newNode; p != nil; p = p.parent {
		updateSize(p)
	}

	fixInsert(t, newNode)
	return true
}

// Delete removes a key (and its associated value) from the tree if it exists.
func Delete[K cmp.Ordered, V any](t *Tree[K, V], key K) bool {
	z, found := Search(t, key)
	if !found {
		return false
	}

	var y = z
	originalColor := y.color
	var x *Node[K, V]
	var fixFrom *Node[K, V] = z.parent

	if z.left == nil {
		x = z.right
		transplant(t, z, z.right)
	} else if z.right == nil {
		x = z.left
		transplant(t, z, z.left)
	} else {
		y = minNode(z.right)
		originalColor = y.color
		x = y.right
		if y.parent == z {
			if x != nil {
				x.parent = y
			}
			fixFrom = y
		} else {
			transplant(t, y, y.right)
			y.right = z.right
			y.right.parent = y
			fixFrom = y.parent
		}
		transplant(t, z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}
	if originalColor == black && x != nil {
		fixDelete(t, x)
	}
	for p := fixFrom; p != nil; p = p.parent {
		updateSize(p)
	}
	return true
}

// Search looks for a key in the tree and returns the node and whether it was found.
func Search[K cmp.Ordered, V any](t *Tree[K, V], key K) (*Node[K, V], bool) {
	curr := t.Root
	for curr != nil {
		if key < curr.key {
			curr = curr.left
		} else if key > curr.key {
			curr = curr.right
		} else {
			return curr, true
		}
	}
	return nil, false
}

// Min returns the node with the smallest key in the tree, or false if the tree is empty.
func Min[K cmp.Ordered, V any](t *Tree[K, V]) (*Node[K, V], bool) {
	if t.Root == nil {
		return nil, false
	}
	return minNode(t.Root), true
}

// Max returns the node with the largest key in the tree, or false if the tree is empty.
func Max[K cmp.Ordered, V any](t *Tree[K, V]) (*Node[K, V], bool) {
	if t.Root == nil {
		return nil, false
	}
	return maxNode(t.Root), true
}

// Ceiling returns the node with the smallest key greater than or equal to the given key.
func Ceiling[K cmp.Ordered, V any](t *Tree[K, V], key K) (*Node[K, V], bool) {
	curr := t.Root
	var result *Node[K, V]
	for curr != nil {
		if key == curr.key {
			return curr, true
		} else if key < curr.key {
			result = curr
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return result, result != nil
}

// Floor returns the node with the largest key less than or equal to the given key.
func Floor[K cmp.Ordered, V any](t *Tree[K, V], key K) (*Node[K, V], bool) {
	curr := t.Root
	var result *Node[K, V]
	for curr != nil {
		if key == curr.key {
			return curr, true
		} else if key < curr.key {
			curr = curr.left
		} else {
			result = curr
			curr = curr.right
		}
	}
	return result, result != nil
}

// Higher returns the node with the smallest key strictly greater than the given key.
func Higher[K cmp.Ordered, V any](t *Tree[K, V], key K) (*Node[K, V], bool) {
	curr := t.Root
	var result *Node[K, V]
	for curr != nil {
		if key < curr.key {
			result = curr
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return result, result != nil
}

// Lower returns the node with the largest key strictly less than the given key.
func Lower[K cmp.Ordered, V any](t *Tree[K, V], key K) (*Node[K, V], bool) {
	curr := t.Root
	var result *Node[K, V]
	for curr != nil {
		if key <= curr.key {
			curr = curr.left
		} else {
			result = curr
			curr = curr.right
		}
	}
	return result, result != nil
}

// Rank returns the number of nodes with keys less than the given key.
func Rank[K cmp.Ordered, V any](t *Tree[K, V], key K) int {
	rank := 0
	curr := t.Root
	for curr != nil {
		if key < curr.key {
			curr = curr.left
		} else {
			leftSize := 0
			if curr.left != nil {
				leftSize = curr.left.size
			}
			if key == curr.key {
				rank += leftSize
				break
			}
			rank += leftSize + 1
			curr = curr.right
		}
	}
	return rank
}

// Kth returns the node that is the k-th smallest (0-based) key.
func Kth[K cmp.Ordered, V any](t *Tree[K, V], k int) (*Node[K, V], bool) {
	curr := t.Root
	for curr != nil {
		leftSize := 0
		if curr.left != nil {
			leftSize = curr.left.size
		}
		if k < leftSize {
			curr = curr.left
		} else if k > leftSize {
			k -= leftSize + 1
			curr = curr.right
		} else {
			return curr, true
		}
	}
	return nil, false
}

// InOrder returns an iterator that traverses the tree in-order.
func InOrder[K cmp.Ordered, V any](t *Tree[K, V]) iter.Seq[Node[K, V]] {
	return func(yield func(Node[K, V]) bool) {
		stack := []*Node[K, V]{}
		curr := t.Root

		for curr != nil || len(stack) > 0 {
			for curr != nil {
				stack = append(stack, curr)
				curr = curr.left
			}
			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if !yield(*n) {
				return
			}
			curr = n.right
		}
	}
}

// Range returns an iterator that traverses nodes with keys in [from, to).
func Range[K cmp.Ordered, V any](t *Tree[K, V], from, to K) iter.Seq[Node[K, V]] {
	return func(yield func(Node[K, V]) bool) {
		stack := []*Node[K, V]{}
		curr := t.Root

		for curr != nil || len(stack) > 0 {
			for curr != nil {
				stack = append(stack, curr)
				curr = curr.left
			}
			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if n.key >= from && n.key < to {
				if !yield(*n) {
					return
				}
			}
			if n.key >= to {
				curr = nil
			} else {
				curr = n.right
			}
		}
	}
}

// Predecessor returns the node with the largest key smaller than the given node. Returns false if no such node exists.
func Predecessor[K cmp.Ordered, V any](n *Node[K, V]) (*Node[K, V], bool) {
	if n.left != nil {
		return maxNode(n.left), true
	}
	y := n.parent
	for y != nil && n == y.left {
		n = y
		y = y.parent
	}
	if y == nil {
		return nil, false
	}
	return y, true
}

// Successor returns the node with the smallest key greater than the given node. Returns false if no such node exists.
func Successor[K cmp.Ordered, V any](n *Node[K, V]) (*Node[K, V], bool) {
	if n.right != nil {
		return minNode(n.right), true
	}
	y := n.parent
	for y != nil && n == y.right {
		n = y
		y = y.parent
	}
	if y == nil {
		return nil, false
	}
	return y, true
}

func fixInsert[K cmp.Ordered, V any](t *Tree[K, V], z *Node[K, V]) {
	for z != t.Root && z.parent.color == red {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y != nil && y.color == red {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					leftRotate(t, z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				rightRotate(t, z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if y != nil && y.color == red {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					rightRotate(t, z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				leftRotate(t, z.parent.parent)
			}
		}
	}
	t.Root.color = black
}

func fixDelete[K cmp.Ordered, V any](t *Tree[K, V], x *Node[K, V]) {
	for x != t.Root && colorOf(x) == black {
		if x == x.parent.left {
			w := x.parent.right
			if colorOf(w) == red {
				w.color = black
				x.parent.color = red
				leftRotate(t, x.parent)
				w = x.parent.right
			}
			if colorOf(w.left) == black && colorOf(w.right) == black {
				w.color = red
				x = x.parent
			} else {
				if colorOf(w.right) == black {
					w.left.color = black
					w.color = red
					rightRotate(t, w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = black
				w.right.color = black
				leftRotate(t, x.parent)
				x = t.Root
			}
		} else {
			w := x.parent.left
			if colorOf(w) == red {
				w.color = black
				x.parent.color = red
				rightRotate(t, x.parent)
				w = x.parent.left
			}
			if colorOf(w.right) == black && colorOf(w.left) == black {
				w.color = red
				x = x.parent
			} else {
				if colorOf(w.left) == black {
					w.right.color = black
					w.color = red
					leftRotate(t, w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = black
				w.left.color = black
				rightRotate(t, x.parent)
				x = t.Root
			}
		}
	}
	if x != nil {
		x.color = black
	}
}

func leftRotate[K cmp.Ordered, V any](t *Tree[K, V], x *Node[K, V]) {
	y := x.right
	if y == nil {
		return
	}
	setRight(x, y.left)
	y.parent = x.parent
	if x.parent == nil {
		t.Root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	setLeft(y, x)
	updateSize(x)
	updateSize(y)
}

func rightRotate[K cmp.Ordered, V any](t *Tree[K, V], y *Node[K, V]) {
	x := y.left
	if x == nil {
		return
	}
	setLeft(y, x.right)
	x.parent = y.parent
	if y.parent == nil {
		t.Root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}
	setRight(x, y)
	updateSize(y)
	updateSize(x)
}

// colorOf returns the color of a node, defaulting to black if the node is nil.
func colorOf[K cmp.Ordered, V any](n *Node[K, V]) color {
	if n == nil {
		return black
	}
	return n.color
}

// transplant replaces one subtree as a child of its parent with another subtree.
func transplant[K cmp.Ordered, V any](t *Tree[K, V], u, v *Node[K, V]) {
	if u.parent == nil {
		t.Root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

// updateSize recalculates the size of a subtree rooted at the given node.
func updateSize[K cmp.Ordered, V any](n *Node[K, V]) {
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}

// minNode returns the node with the minimum key starting from the given node.
func minNode[K cmp.Ordered, V any](n *Node[K, V]) *Node[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}

// maxNode returns the node with the maximum key starting from the given node.
func maxNode[K cmp.Ordered, V any](n *Node[K, V]) *Node[K, V] {
	for n.right != nil {
		n = n.right
	}
	return n
}

// setLeft sets the left child of a parent node and updates the child's parent pointer.
func setLeft[K cmp.Ordered, V any](parent, child *Node[K, V]) {
	parent.left = child
	if child != nil {
		child.parent = parent
	}
}

// setRight sets the right child of a parent node and updates the child's parent pointer.
func setRight[K cmp.Ordered, V any](parent, child *Node[K, V]) {
	parent.right = child
	if child != nil {
		child.parent = parent
	}
}
