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

- `Insert(tree *Tree[K, V], key K, value V) *node[K, V]`
- `Delete(tree *Tree[K, V], key K)`
- `Search(tree *Tree[K, V], key K) (*node[K, V], bool)`
- `Min(tree *Tree[K, V]) (*node[K, V], bool)`
- `Max(tree *Tree[K, V]) (*node[K, V], bool)`
- `Rank(tree *Tree[K, V], key K) int`
- `Kth(tree *Tree[K, V], k int) (*node[K, V], bool)`
- `Ceiling(tree *Tree[K, V], key K) (*node[K, V], bool)`
- `Floor(tree *Tree[K, V], key K) (*node[K, V], bool)`
- `Higher(tree *Tree[K, V], key K) (*node[K, V], bool)`
- `Lower(tree *Tree[K, V], key K) (*node[K, V], bool)`
- `Predecessor(n *node[K, V]) (*node[K, V], bool)`
- `Successor(n *node[K, V]) (*node[K, V], bool)`
- `InOrder(tree *Tree[K, V]) iter.Seq[node[K, V]]`
- `Range(tree *Tree[K, V], from, to K) iter.Seq[node[K, V]]`

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.