# redblacktrees [![GoDoc](https://pkg.go.dev/badge/github.com/byExist/redblacktrees)](https://pkg.go.dev/github.com/byExist/redblacktrees) [![Go Report Card](https://goreportcard.com/badge/github.com/byExist/redblacktrees)](https://goreportcard.com/report/github.com/byExist/redblacktrees)

A generic Red-Black Tree for Go with rank, range, and k-th queries.

---

## ✨ Features

- Generic Red-Black Tree using Go generics
- O(log n) insert, delete, and search
- Rank, k-th element, ceiling, floor, range, predecessor/successor
- In-order iterator
- Tree size maintained for fast queries

---

## ✅ Use When

- You need **fast insertions and deletions**
- You want **balanced performance across reads and writes**
- You need **ordered map-like behavior** with efficient range queries

---

## 🚫 Avoid If

- You need **maximum search performance** → try [AVL Tree](https://github.com/byExist/avltrees)
- You need **concurrent** access (not thread-safe)

---

## 📦 Installation

```bash
go get github.com/byExist/redblacktrees
```

---

## 🚀 Quick Start

```go
package main

import (
	"fmt"
	rbts "github.com/byExist/redblacktrees"
)

func main() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 3, "three")
	rbts.Insert(tree, 1, "one")
	rbts.Insert(tree, 2, "two")

	if node, found := rbts.Search(tree, 2); found {
		fmt.Println("Found:", node.Value)
	}

	for node := range rbts.InOrder(tree) {
		fmt.Printf("%d: %s\n", node.Key, node.Value)
	}

	fmt.Println("Rank of 3:", rbts.Rank(tree, 3))

	if node, ok := rbts.Kth(tree, 0); ok {
		fmt.Println("0-th smallest:", node.Key)
	}
}
```

---

## 📊 Performance

Benchmarked on Apple M1 Pro:

| Operation            | Time (ns/op) | Memory (B/op) | Allocations |
|---------------------|--------------|----------------|-------------|
| Insert (Random)     | 791.1        | 64 B           | 1           |
| Insert (Sequential) | 101.5        | 64 B           | 1           |
| Search (Hit)        | 10.60        | 0 B            | 0           |
| Search (Miss)       | 12.21        | 0 B            | 0           |
| Delete (Random)     | 2.36         | 0 B            | 0           |

---

## 📚 Documentation

Full API reference: [pkg.go.dev/github.com/byExist/redblacktrees](https://pkg.go.dev/github.com/byExist/redblacktrees)

---

## 🪪 License

MIT License. See [LICENSE](LICENSE).