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

type Node[K cmp.Ordered, V any] struct {
	key    K
	value  V
	color  color
	left   *Node[K, V]
	right  *Node[K, V]
	parent *Node[K, V]
	size   int
}

func (n *Node[K, V]) Key() K {
	return n.key
}

func (n *Node[K, V]) Value() V {
	return n.value
}

type Tree[K cmp.Ordered, V any] struct {
	Root *Node[K, V]
}

func New[K cmp.Ordered, V any]() *Tree[K, V] {
	return &Tree[K, V]{}
}

func Len[K cmp.Ordered, V any](t *Tree[K, V]) int {
	if t.Root == nil {
		return 0
	}
	return t.Root.size
}

func Clear[K cmp.Ordered, V any](t *Tree[K, V]) {
	t.Root = nil
}

func Insert[K cmp.Ordered, V any](t *Tree[K, V], key K, value V) bool {
	curr := t.Root
	var parent *Node[K, V]

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

	newNode := &Node[K, V]{
		key:    key,
		value:  value,
		color:  red,
		parent: parent,
		size:   1,
	}

	if parent == nil {
		t.Root = newNode
	} else if key < parent.key {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	for n := newNode.parent; n != nil; n = n.parent {
		updateSize(n)
	}

	fixInsert(t, newNode)
	t.Root.color = black
	return true
}

func Delete[K cmp.Ordered, V any](t *Tree[K, V], key K) bool {
	node, found := Search(t, key)
	if !found {
		return false
	}

	var y *Node[K, V] = node
	originalColor := y.color
	var x *Node[K, V]
	var xParent *Node[K, V]

	if node.left == nil {
		x = node.right
		transplant(t, node, node.right)
		xParent = node.parent
	} else if node.right == nil {
		x = node.left
		transplant(t, node, node.left)
		xParent = node.parent
	} else {
		y = minNode(node.right)
		originalColor = y.color
		x = y.right
		if y.parent == node {
			if x != nil {
				x.parent = y
			}
			xParent = y
		} else {
			transplant(t, y, y.right)
			y.right = node.right
			y.right.parent = y
			xParent = y.parent
		}
		transplant(t, node, y)
		y.left = node.left
		y.left.parent = y
		y.color = node.color
	}

	for p := xParent; p != nil; p = p.parent {
		updateSize(p)
	}

	if originalColor == black {
		fixDelete(t, x, xParent)
	}

	return true
}

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

func Min[K cmp.Ordered, V any](t *Tree[K, V]) (*Node[K, V], bool) {
	if t.Root == nil {
		return nil, false
	}
	curr := t.Root
	for curr.left != nil {
		curr = curr.left
	}
	return curr, true
}

func Max[K cmp.Ordered, V any](t *Tree[K, V]) (*Node[K, V], bool) {
	if t.Root == nil {
		return nil, false
	}
	curr := t.Root
	for curr.right != nil {
		curr = curr.right
	}
	return curr, true
}

func InOrder[K cmp.Ordered, V any](t *Tree[K, V]) iter.Seq[*Node[K, V]] {
	return func(yield func(*Node[K, V]) bool) {
		stack := []*Node[K, V]{}
		curr := t.Root
		for curr != nil || len(stack) > 0 {
			for curr != nil {
				stack = append(stack, curr)
				curr = curr.left
			}
			curr = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if !yield(curr) {
				return
			}
			curr = curr.right
		}
	}
}

func Range[K cmp.Ordered, V any](t *Tree[K, V], low, high K) iter.Seq[*Node[K, V]] {
	return func(yield func(*Node[K, V]) bool) {
		stack := []*Node[K, V]{}
		curr := t.Root
		for curr != nil || len(stack) > 0 {
			for curr != nil {
				stack = append(stack, curr)
				curr = curr.left
			}
			curr = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if curr.key > high {
				return
			}
			if curr.key >= low && curr.key <= high {
				if !yield(curr) {
					return
				}
			}

			curr = curr.right
		}
	}
}

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
	if result != nil {
		return result, true
	}
	return nil, false
}

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
	if result != nil {
		return result, true
	}
	return nil, false
}

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
	if result != nil {
		return result, true
	}
	return nil, false
}

func Lower[K cmp.Ordered, V any](t *Tree[K, V], key K) (*Node[K, V], bool) {
	curr := t.Root
	var result *Node[K, V]
	for curr != nil {
		if key > curr.key {
			result = curr
			curr = curr.right
		} else {
			curr = curr.left
		}
	}
	if result != nil {
		return result, true
	}
	return nil, false
}

func Rank[K cmp.Ordered, V any](t *Tree[K, V], key K) int {
	rank := 0
	curr := t.Root
	for curr != nil {
		if key < curr.key {
			curr = curr.left
		} else if key > curr.key {
			leftSize := 0
			if curr.left != nil {
				leftSize = curr.left.size
			}
			rank += leftSize + 1
			curr = curr.right
		} else {
			if curr.left != nil {
				rank += curr.left.size
			}
			break
		}
	}
	return rank
}

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

func Predecessor[K cmp.Ordered, V any](n *Node[K, V]) (*Node[K, V], bool) {
	if n.left != nil {
		p := n.left
		for p.right != nil {
			p = p.right
		}
		return p, true
	}
	p := n.parent
	for p != nil && n == p.left {
		n = p
		p = p.parent
	}
	if p != nil {
		return p, true
	}
	return nil, false
}

func Successor[K cmp.Ordered, V any](n *Node[K, V]) (*Node[K, V], bool) {
	if n.right != nil {
		p := n.right
		for p.left != nil {
			p = p.left
		}
		return p, true
	}
	p := n.parent
	for p != nil && n == p.right {
		n = p
		p = p.parent
	}
	if p != nil {
		return p, true
	}
	return nil, false
}

func (n *Node[K, V]) grandparent() *Node[K, V] {
	if n.parent == nil {
		return nil
	}
	return n.parent.parent
}

func (n *Node[K, V]) uncle() *Node[K, V] {
	g := n.grandparent()
	if g == nil {
		return nil
	}
	if n.parent == g.left {
		return g.right
	}
	return g.left
}

func isRed[K cmp.Ordered, V any](n *Node[K, V]) bool {
	return n != nil && n.color == red
}

func isBlack[K cmp.Ordered, V any](n *Node[K, V]) bool {
	return n == nil || n.color == black
}

func setColor[K cmp.Ordered, V any](n *Node[K, V], c color) {
	if n != nil {
		n.color = c
	}
}

func updateSize[K cmp.Ordered, V any](n *Node[K, V]) {
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}

func rotateLeft[K cmp.Ordered, V any](t *Tree[K, V], x *Node[K, V]) {
	y := x.right
	if y == nil {
		return
	}
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		t.Root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y

	updateSize(x)
	updateSize(y)
}

func rotateRight[K cmp.Ordered, V any](t *Tree[K, V], y *Node[K, V]) {
	x := y.left
	if x == nil {
		return
	}
	y.left = x.right
	if x.right != nil {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == nil {
		t.Root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}
	x.right = y
	y.parent = x

	updateSize(y)
	updateSize(x)
}

func fixInsert[K cmp.Ordered, V any](t *Tree[K, V], n *Node[K, V]) {
	for n != t.Root && isRed(n.parent) {
		if n.parent == n.grandparent().left {
			uncle := n.uncle()
			if isRed(uncle) {
				setColor(n.parent, black)
				setColor(uncle, black)
				setColor(n.grandparent(), red)
				n = n.grandparent()
			} else {
				if n == n.parent.right {
					n = n.parent
					rotateLeft(t, n)
				}
				setColor(n.parent, black)
				setColor(n.grandparent(), red)
				rotateRight(t, n.grandparent())
			}
		} else {
			uncle := n.uncle()
			if isRed(uncle) {
				setColor(n.parent, black)
				setColor(uncle, black)
				setColor(n.grandparent(), red)
				n = n.grandparent()
			} else {
				if n == n.parent.left {
					n = n.parent
					rotateRight(t, n)
				}
				setColor(n.parent, black)
				setColor(n.grandparent(), red)
				rotateLeft(t, n.grandparent())
			}
		}
	}
	setColor(t.Root, black)
}

func fixDelete[K cmp.Ordered, V any](t *Tree[K, V], x, parent *Node[K, V]) {
	for x != t.Root && isBlack(x) {
		if parent != nil && x == parent.left {
			sibling := parent.right
			if isRed(sibling) {
				setColor(sibling, black)
				setColor(parent, red)
				rotateLeft(t, parent)
				sibling = parent.right
			}
			if isBlack(sibling.left) && isBlack(sibling.right) {
				setColor(sibling, red)
				x = parent
				parent = x.parent
			} else {
				if isBlack(sibling.right) {
					setColor(sibling.left, black)
					setColor(sibling, red)
					rotateRight(t, sibling)
					sibling = parent.right
				}
				setColor(sibling, parent.color)
				setColor(parent, black)
				setColor(sibling.right, black)
				rotateLeft(t, parent)
				x = t.Root
			}
		} else if parent != nil {
			sibling := parent.left
			if isRed(sibling) {
				setColor(sibling, black)
				setColor(parent, red)
				rotateRight(t, parent)
				sibling = parent.left
			}
			if isBlack(sibling.left) && isBlack(sibling.right) {
				setColor(sibling, red)
				x = parent
				parent = x.parent
			} else {
				if isBlack(sibling.left) {
					setColor(sibling.right, black)
					setColor(sibling, red)
					rotateLeft(t, sibling)
					sibling = parent.left
				}
				setColor(sibling, parent.color)
				setColor(parent, black)
				setColor(sibling.left, black)
				rotateRight(t, parent)
				x = t.Root
			}
		} else {
			break
		}
	}
	if x != nil {
		setColor(x, black)
	}
}

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

func minNode[K cmp.Ordered, V any](n *Node[K, V]) *Node[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}
