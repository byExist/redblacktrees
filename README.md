# redblacktrees [![GoDoc](https://pkg.go.dev/badge/github.com/byExist/redblacktrees.svg)](https://pkg.go.dev/github.com/byExist/redblacktrees) [![Go Report Card](https://goreportcard.com/badge/github.com/byExist/redblacktrees)](https://goreportcard.com/report/github.com/byExist/redblacktrees)

## What is "redblacktrees"?

`redblacktrees` is a generic Red-Black Tree implementation in Go, designed to provide efficient insertion, deletion, and searching operations with logarithmic time complexity. It supports advanced operations such as rank, k-th element, range queries, and more.

## Features

- Generic support using Go generics
- Efficient Insert, Delete, and Search operations (O(log n))
- Support for advanced queries:
  - Rank
  - k-th smallest element
  - Range queries
  - Ceiling
  - Floor
  - Higher
  - Lower
  - Predecessor / Successor
- In-order traversal using iter.Seq
- Tree size maintained at each node for fast queries

## Installation

To install, use the following command:

```bash
go get github.com/byExist/redblacktrees
```

## Quick Start

```go
package main

import (
	"fmt"
	rbts "github.com/byExist/redblacktrees"
)

func main() {
	tree := rbts.New[int, string]()

	// Insert elements
	rbts.Insert(tree, 3, "three")
	rbts.Insert(tree, 1, "one")
	rbts.Insert(tree, 2, "two")

	// Search
	if node, found := rbts.Search(tree, 2); found {
		fmt.Println("Found:", node.Value)
	}

	// Delete
	rbts.Delete(tree, 1)

	// In-order traversal
	for node := range rbts.InOrder(tree) {
		fmt.Printf("%d: %s\n", node.Key, node.Value)
	}

	// Rank
	fmt.Println("Rank of key 3:", rbts.Rank(tree, 3))

	// K-th smallest
	if node, ok := rbts.Kth(tree, 0); ok {
		fmt.Println("0-th smallest:", node.Key)
	}
}
```

## Usage

The `redblacktrees` package provides a robust implementation of red-black trees with support for generic types and useful query functions. It is ideal for ordered maps and fast range queries where performance and balance are critical.

## API Overview

### Constructors

- `New[K cmp.Ordered, V any]() *Tree[K, V]`


### Core Functions
- `Insert(t *Tree[K, V], key K, value V) bool`
- `Delete(t *Tree[K, V], key K) bool`
- `Search(t *Tree[K, V], key K) (*Node[K, V], bool)`
- `Len(t *Tree[K, V]) int`
- `Clear(t *Tree[K, V])`
- `Min(t *Tree[K, V]) (*Node[K, V], bool)`
- `Max(t *Tree[K, V]) (*Node[K, V], bool)`
- `Ceiling(t *Tree[K, V], key K) (*Node[K, V], bool)`
- `Floor(t *Tree[K, V], key K) (*Node[K, V], bool)`
- `Higher(t *Tree[K, V], key K) (*Node[K, V], bool)`
- `Lower(t *Tree[K, V], key K) (*Node[K, V], bool)`
- `Rank(t *Tree[K, V], key K) int`
- `Kth(t *Tree[K, V], k int) (*Node[K, V], bool)`
- `InOrder(t *Tree[K, V]) iter.Seq[Node[K, V]]`
- `Range(t *Tree[K, V], from, to K) iter.Seq[Node[K, V]]`
- `Predecessor(n *Node[K, V]) (*Node[K, V], bool)`
- `Successor(n *Node[K, V]) (*Node[K, V], bool)`

### Node

The `Node[K, V]` type represents a single node in the red-black tree. It provides access to the key and value stored in the node.

- `Key() K`: Returns the key of the node.
- `Value() V`: Returns the value stored in the node.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.