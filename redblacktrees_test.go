package redblacktrees_test

import (
	"fmt"
	"testing"

	rbts "github.com/byExist/redblacktrees"
)

func TestNew(t *testing.T) {
	tree := rbts.New[int, string]()
	if tree.Root != nil {
		t.Errorf("New tree should have nil Root")
	}
	if rbts.Len(tree) != 0 {
		t.Errorf("New tree should have size 0, got %d", rbts.Len(tree))
	}
}

func TestInsert(t *testing.T) {
	tree := rbts.New[int, string]()
	inserted := rbts.Insert(tree, 10, "TEN")
	if !inserted {
		t.Errorf("Expected first insert of 10 to return true")
	}
	inserted = rbts.Insert(tree, 10, "ten")
	if inserted {
		t.Errorf("Expected second insert of 10 to return false (overwrite)")
	}
	rbts.Insert(tree, 20, "twenty")
	rbts.Insert(tree, 5, "five")

	if rbts.Len(tree) != 3 {
		t.Errorf("Expected size 3, got %d", rbts.Len(tree))
	}

	node, found := rbts.Search(tree, 10)
	if !found || node.Value() != "ten" {
		t.Errorf("Insert failed for key 10")
	}
}

func TestDelete(t *testing.T) {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "ten")
	rbts.Insert(tree, 20, "twenty")
	rbts.Insert(tree, 5, "five")

	rbts.Delete(tree, 10)
	if rbts.Len(tree) != 2 {
		t.Errorf("Expected size 2 after deletion, got %d", rbts.Len(tree))
	}

	_, found := rbts.Search(tree, 10)
	if found {
		t.Errorf("Key 10 should have been deleted")
	}
}

func TestSearch(t *testing.T) {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "ten")
	rbts.Insert(tree, 20, "twenty")

	node, found := rbts.Search(tree, 10)
	if !found || node.Value() != "ten" {
		t.Errorf("Search failed for existing key 10")
	}

	_, found = rbts.Search(tree, 30)
	if found {
		t.Errorf("Search should fail for non-existent key 30")
	}
}

func TestInOrder(t *testing.T) {
	tree := rbts.New[int, string]()
	values := []int{20, 10, 30, 5, 15, 25, 35}
	for _, v := range values {
		rbts.Insert(tree, v, "")
	}

	prev := -1
	for n := range rbts.InOrder(tree) {
		if prev != -1 && prev >= n.Key() {
			t.Errorf("InOrder traversal is not sorted: prev=%d, current=%d", prev, n.Key())
		}
		prev = n.Key()
	}
}

func TestLowerBound(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30, 40, 50} {
		rbts.Insert(tree, v, "")
	}

	n := rbts.LowerBound(tree, 25)
	if n == nil || n.Key() != 30 {
		t.Errorf("Expected LowerBound(25) to be 30, got %v", n)
	}

	n = rbts.LowerBound(tree, 50)
	if n == nil || n.Key() != 50 {
		t.Errorf("Expected LowerBound(50) to be 50, got %v", n)
	}

	n = rbts.LowerBound(tree, 60)
	if n != nil {
		t.Errorf("Expected LowerBound(60) to be nil, got %v", n)
	}
}

func TestUpperBound(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30, 40, 50} {
		rbts.Insert(tree, v, "")
	}

	n := rbts.UpperBound(tree, 25)
	if n == nil || n.Key() != 30 {
		t.Errorf("Expected UpperBound(25) to be 30, got %v", n)
	}

	n = rbts.UpperBound(tree, 50)
	if n != nil {
		t.Errorf("Expected UpperBound(50) to be nil, got %v", n)
	}
}

func TestRange(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30, 40, 50} {
		rbts.Insert(tree, v, "")
	}

	var collected []int
	for n := range rbts.Range(tree, 15, 45) {
		collected = append(collected, n.Key())
	}

	expected := []int{20, 30, 40}
	if len(collected) != len(expected) {
		t.Errorf("Expected range length %d, got %d", len(expected), len(collected))
	}

	for i, v := range expected {
		if collected[i] != v {
			t.Errorf("Expected %d at position %d, got %d", v, i, collected[i])
		}
	}
}

func TestRank(t *testing.T) {
	tree := rbts.New[int, string]()
	values := []int{10, 20, 30, 40, 50}
	for _, v := range values {
		rbts.Insert(tree, v, "")
	}

	if r := rbts.Rank(tree, 25); r != 2 {
		t.Errorf("Expected Rank(25) = 2, got %d", r)
	}
	if r := rbts.Rank(tree, 10); r != 0 {
		t.Errorf("Expected Rank(10) = 0, got %d", r)
	}
	if r := rbts.Rank(tree, 60); r != 5 {
		t.Errorf("Expected Rank(60) = 5, got %d", r)
	}
}

func TestKth(t *testing.T) {
	tree := rbts.New[int, string]()
	values := []int{10, 20, 30, 40, 50}
	for _, v := range values {
		rbts.Insert(tree, v, "")
	}

	n, ok := rbts.Kth(tree, 0)
	if !ok || n.Key() != 10 {
		t.Errorf("Expected 0th key = 10, got %v", n)
	}

	n, ok = rbts.Kth(tree, 3)
	if !ok || n.Key() != 40 {
		t.Errorf("Expected 3rd key = 40, got %v", n)
	}

	n, ok = rbts.Kth(tree, 5)
	if ok {
		t.Errorf("Expected 5th key to be invalid, got %v", n)
	}
}

func TestPredecessor(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30, 40, 50} {
		rbts.Insert(tree, v, "")
	}

	n, _ := rbts.Search(tree, 30)
	pred := rbts.Predecessor(n)
	if pred == nil || pred.Key() != 20 {
		t.Errorf("Expected predecessor of 30 to be 20, got %v", pred)
	}
}

func TestSuccessor(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30, 40, 50} {
		rbts.Insert(tree, v, "")
	}

	n, _ := rbts.Search(tree, 30)
	succ := rbts.Successor(n)
	if succ == nil || succ.Key() != 40 {
		t.Errorf("Expected successor of 30 to be 40, got %v", succ)
	}
}

func TestMin(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{20, 10, 30} {
		rbts.Insert(tree, v, "")
	}

	m := rbts.Min(tree)
	if m == nil || m.Key() != 10 {
		t.Errorf("Expected Min = 10, got %v", m)
	}
}

func TestMax(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{20, 10, 30} {
		rbts.Insert(tree, v, "")
	}

	m := rbts.Max(tree)
	if m == nil || m.Key() != 30 {
		t.Errorf("Expected Max = 30, got %v", m)
	}
}

func TestLen(t *testing.T) {
	tree := rbts.New[int, string]()
	if rbts.Len(tree) != 0 {
		t.Errorf("Expected size 0, got %d", rbts.Len(tree))
	}
	rbts.Insert(tree, 1, "")
	rbts.Insert(tree, 2, "")
	rbts.Insert(tree, 3, "")
	if rbts.Len(tree) != 3 {
		t.Errorf("Expected size 3, got %d", rbts.Len(tree))
	}
}

func ExampleNew() {
	tree := rbts.New[int, string]()
	fmt.Println(rbts.Len(tree))
	// Output: 0
}

func ExampleInsert() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "ten")
	rbts.Insert(tree, 5, "five")
	rbts.Insert(tree, 15, "fifteen")
	fmt.Println(rbts.Len(tree))
	// Output: 3
}

func ExampleDelete() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "ten")
	rbts.Delete(tree, 10)
	fmt.Println(rbts.Len(tree))
	// Output: 0
}

func ExampleSearch() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 20, "twenty")
	node, found := rbts.Search(tree, 20)
	fmt.Println(found, node.Value())
	// Output: true twenty
}

func ExampleMin() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 20, "")
	rbts.Insert(tree, 10, "")
	min := rbts.Min(tree)
	fmt.Println(min.Key())
	// Output: 10
}

func ExampleMax() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 20, "")
	rbts.Insert(tree, 30, "")
	max := rbts.Max(tree)
	fmt.Println(max.Key())
	// Output: 30
}

func ExampleLowerBound() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "")
	rbts.Insert(tree, 20, "")
	n := rbts.LowerBound(tree, 15)
	fmt.Println(n.Key())
	// Output: 20
}

func ExampleUpperBound() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "")
	rbts.Insert(tree, 20, "")
	n := rbts.UpperBound(tree, 10)
	fmt.Println(n.Key())
	// Output: 20
}

func ExampleInOrder() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 20, "")
	rbts.Insert(tree, 10, "")
	rbts.Insert(tree, 30, "")
	for n := range rbts.InOrder(tree) {
		fmt.Print(n.Key(), " ")
	}
	fmt.Println()
	// Output: 10 20 30
}

func ExampleRange() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "")
	rbts.Insert(tree, 20, "")
	rbts.Insert(tree, 30, "")
	for n := range rbts.Range(tree, 15, 25) {
		fmt.Print(n.Key(), " ")
	}
	fmt.Println()
	// Output: 20
}

func ExamplePredecessor() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 20, "")
	rbts.Insert(tree, 10, "")
	node, _ := rbts.Search(tree, 20)
	pred := rbts.Predecessor(node)
	fmt.Println(pred.Key())
	// Output: 10
}

func ExampleSuccessor() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 20, "")
	rbts.Insert(tree, 30, "")
	node, _ := rbts.Search(tree, 20)
	succ := rbts.Successor(node)
	fmt.Println(succ.Key())
	// Output: 30
}

func ExampleRank() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "")
	rbts.Insert(tree, 20, "")
	rank := rbts.Rank(tree, 15)
	fmt.Println(rank)
	// Output: 1
}

func ExampleKth() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "")
	rbts.Insert(tree, 20, "")
	n, _ := rbts.Kth(tree, 1)
	fmt.Println(n.Key())
	// Output: 20
}

func ExampleLen() {
	tree := rbts.New[int, string]()
	fmt.Println(rbts.Len(tree))
	// Output: 0
}
