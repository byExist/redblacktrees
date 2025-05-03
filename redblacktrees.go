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
	key    K
	value  V
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

// Tree represents the root of a red-black tree.
type Tree[K cmp.Ordered, V any] struct {
	Root *Node[K, V]
}

// New returns a new empty Red-Black Tree.
func New[K cmp.Ordered, V any]() *Tree[K, V] {
	return &Tree[K, V]{}
}

// Clear sets the tree root to nil, effectively clearing the tree.
func Clear[K cmp.Ordered, V any](t *Tree[K, V]) {
	t.Root = nil
}

// Insert inserts a new key-value pair into the red-black tree.
// Returns true if inserted, false if replaced.
func Insert[K cmp.Ordered, V any](t *Tree[K, V], key K, value V) bool {
	z := &Node[K, V]{key: key, value: value, color: red, size: 1}
	y := (*Node[K, V])(nil)
	x := t.Root

	for x != nil {
		y = x
		x.size++
		if key < x.key {
			x = x.left
		} else if key > x.key {
			x = x.right
		} else {
			x.value = value
			// restore sizes on the path back up
			for y != nil {
				updateSize(y)
				y = y.parent
			}
			return false
		}
	}

	z.parent = y
	if y == nil {
		t.Root = z
	} else if key < y.key {
		y.left = z
	} else {
		y.right = z
	}
	insertFixup(t, z)
	return true
}

// Delete removes a node with the given key from the red-black tree.
func Delete[K cmp.Ordered, V any](t *Tree[K, V], key K) bool {
	z := t.Root
	for z != nil {
		if key < z.key {
			z = z.left
		} else if key > z.key {
			z = z.right
		} else {
			break
		}
	}
	if z == nil {
		return false
	}

	y := z
	yOriginalColor := y.color
	var x *Node[K, V]

	if z.left == nil {
		x = z.right
		transplant(t, z, z.right)
	} else if z.right == nil {
		x = z.left
		transplant(t, z, z.left)
	} else {
		y = minimum(z.right)
		yOriginalColor = y.color
		x = y.right
		if y.parent == z {
			if x != nil {
				x.parent = y
			}
		} else {
			transplant(t, y, y.right)
			y.right = z.right
			if y.right != nil {
				y.right.parent = y
			}
		}
		transplant(t, z, y)
		y.left = z.left
		if y.left != nil {
			y.left.parent = y
		}
		y.color = z.color
		updateSize(y)
	}
	fixSizeUpward(z.parent)
	if yOriginalColor == black {
		deleteFixup(t, x, z.parent)
	}
	return true
}

// Search finds a node with the given key in the red-black tree.
func Search[K cmp.Ordered, V any](t *Tree[K, V], key K) (*Node[K, V], bool) {
	x := t.Root
	for x != nil {
		if key < x.key {
			x = x.left
		} else if key > x.key {
			x = x.right
		} else {
			return x, true
		}
	}
	return nil, false
}

// Min returns the node with the minimum key in the tree.
func Min[K cmp.Ordered, V any](t *Tree[K, V]) (*Node[K, V], bool) {
	if t.Root == nil {
		return nil, false
	}
	return minimum(t.Root), true
}

// Max returns the node with the maximum key in the tree.
func Max[K cmp.Ordered, V any](t *Tree[K, V]) (*Node[K, V], bool) {
	if t.Root == nil {
		return nil, false
	}
	n := t.Root
	for n.right != nil {
		n = n.right
	}
	return n, true
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

// Floor returns the node with the greatest key less than or equal to the given key.
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

// Higher returns the node with the smallest key greater than the given key.
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

// Lower returns the node with the greatest key less than the given key.
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

// Predecessor returns the in-order predecessor node of n, if any.
func Predecessor[K cmp.Ordered, V any](n *Node[K, V]) (*Node[K, V], bool) {
	if n.left != nil {
		x := n.left
		for x.right != nil {
			x = x.right
		}
		return x, true
	}
	p := n.parent
	for p != nil && n == p.left {
		n = p
		p = p.parent
	}
	return p, p != nil
}

// Successor returns the in-order successor node of n, if any.
func Successor[K cmp.Ordered, V any](n *Node[K, V]) (*Node[K, V], bool) {
	if n.right != nil {
		x := n.right
		for x.left != nil {
			x = x.left
		}
		return x, true
	}
	p := n.parent
	for p != nil && n == p.right {
		n = p
		p = p.parent
	}
	return p, p != nil
}

// InOrder returns an iterator for in-order traversal of the tree.
func InOrder[K cmp.Ordered, V any](t *Tree[K, V]) iter.Seq[Node[K, V]] {
	return func(yield func(Node[K, V]) bool) {
		var stack []*Node[K, V]
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

// Range returns an iterator over nodes with keys in [from, to).
func Range[K cmp.Ordered, V any](t *Tree[K, V], from, to K) iter.Seq[Node[K, V]] {
	return func(yield func(Node[K, V]) bool) {
		var stack []*Node[K, V]
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

// Kth returns the node with the given 0-based rank (k).
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

// Len returns the number of nodes in the tree.
func Len[K cmp.Ordered, V any](t *Tree[K, V]) int {
	if t.Root == nil {
		return 0
	}
	return t.Root.size
}

func updateSize[K cmp.Ordered, V any](n *Node[K, V]) {
	if n == nil {
		return
	}
	n.size = 1
	if n.left != nil {
		n.size += n.left.size
	}
	if n.right != nil {
		n.size += n.right.size
	}
}

func fixSizeUpward[K cmp.Ordered, V any](n *Node[K, V]) {
	for n != nil {
		updateSize(n)
		n = n.parent
	}
}

func insertFixup[K cmp.Ordered, V any](t *Tree[K, V], z *Node[K, V]) {
	for isRed(z.parent) {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if isRed(y) {
				setColor(z.parent, black)
				setColor(y, black)
				setColor(z.parent.parent, red)
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					rotateLeft(t, z)
				}
				setColor(z.parent, black)
				setColor(z.parent.parent, red)
				rotateRight(t, z.parent.parent)
			}
		} else {
			y := z.parent.parent.left
			if isRed(y) {
				setColor(z.parent, black)
				setColor(y, black)
				setColor(z.parent.parent, red)
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					rotateRight(t, z)
				}
				setColor(z.parent, black)
				setColor(z.parent.parent, red)
				rotateLeft(t, z.parent.parent)
			}
		}
	}
	setColor(t.Root, black)
}

func deleteFixup[K cmp.Ordered, V any](t *Tree[K, V], x, parent *Node[K, V]) {
	for x != t.Root && !isRed(x) && parent != nil {
		if x == parent.left {
			w := parent.right
			if isRed(w) {
				setColor(w, black)
				setColor(parent, red)
				rotateLeft(t, parent)
				w = parent.right
			}
			if w == nil || (!isRed(w.left) && !isRed(w.right)) {
				setColor(w, red)
				x = parent
				parent = x.parent
			} else {
				if !isRed(w.right) {
					if w.left != nil {
						setColor(w.left, black)
					}
					setColor(w, red)
					rotateRight(t, w)
					w = parent.right
				}
				setColor(w, parent.color)
				setColor(parent, black)
				setColor(w.right, black)
				rotateLeft(t, parent)
				x = t.Root
				break
			}
		} else {
			w := parent.left
			if isRed(w) {
				setColor(w, black)
				setColor(parent, red)
				rotateRight(t, parent)
				w = parent.left
			}
			if w == nil || (!isRed(w.left) && !isRed(w.right)) {
				setColor(w, red)
				x = parent
				parent = x.parent
			} else {
				if !isRed(w.left) {
					if w.right != nil {
						setColor(w.right, black)
					}
					setColor(w, red)
					rotateLeft(t, w)
					w = parent.left
				}
				setColor(w, parent.color)
				setColor(parent, black)
				if w != nil && w.left != nil {
					setColor(w.left, black)
				}
				rotateRight(t, parent)
				x = t.Root
				break
			}
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

func isRed[K cmp.Ordered, V any](n *Node[K, V]) bool {
	return n != nil && n.color == red
}

func setColor[K cmp.Ordered, V any](n *Node[K, V], c color) {
	if n != nil {
		n.color = c
	}
}

func rotateLeft[K cmp.Ordered, V any](t *Tree[K, V], x *Node[K, V]) {
	y := x.right
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
	y.left = x.right
	if x.right != nil {
		x.right.parent = y
	}
	x.parent = y.parent
	if y.parent == nil {
		t.Root = x
	} else if y == y.parent.right {
		y.parent.right = x
	} else {
		y.parent.left = x
	}
	x.right = y
	y.parent = x
	updateSize(y)
	updateSize(x)
}

func minimum[K cmp.Ordered, V any](n *Node[K, V]) *Node[K, V] {
	for n.left != nil {
		n = n.left
	}
	return n
}
