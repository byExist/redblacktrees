package redblacktrees_test

import (
	"fmt"
	"math/rand"
	"testing"

	rbts "github.com/byExist/redblacktrees"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tree := rbts.New[int, string]()
	assert.Nil(t, tree.Root, "New tree should have nil Root")
	assert.Equal(t, 0, rbts.Len(tree), "New tree should have size 0")
}

func TestLen(t *testing.T) {
	tree := rbts.New[int, string]()
	assert.Equal(t, 0, rbts.Len(tree))
	rbts.Insert(tree, 1, "")
	rbts.Insert(tree, 2, "")
	rbts.Insert(tree, 3, "")
	assert.Equal(t, 3, rbts.Len(tree))
}

func TestInsert(t *testing.T) {
	tree := rbts.New[int, string]()

	inserted := rbts.Insert(tree, 10, "TEN")
	assert.True(t, inserted, "Expected first insert of 10 to return true")
	inserted = rbts.Insert(tree, 10, "ten")
	assert.False(t, inserted, "Expected second insert of 10 to return false (overwrite)")
	rbts.Insert(tree, 20, "twenty")
	rbts.Insert(tree, 5, "five")

	assert.Equal(t, 3, rbts.Len(tree), "Expected size 3")

	node, found := rbts.Search(tree, 10)
	require.True(t, found, "Insert failed for key 10")
	assert.Equal(t, "ten", node.Value())
}

func TestDelete(t *testing.T) {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "ten")
	rbts.Insert(tree, 20, "twenty")
	rbts.Insert(tree, 5, "five")

	rbts.Delete(tree, 10)
	assert.Equal(t, 2, rbts.Len(tree), "Expected size 2 after deletion")

	_, found := rbts.Search(tree, 10)
	assert.False(t, found, "Key 10 should have been deleted")
}

func TestSearch(t *testing.T) {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 10, "ten")
	rbts.Insert(tree, 20, "twenty")

	node, found := rbts.Search(tree, 10)
	require.True(t, found, "Search failed for existing key 10")
	assert.Equal(t, "ten", node.Value())

	_, found = rbts.Search(tree, 30)
	assert.False(t, found, "Search should fail for non-existent key 30")
}

func TestMin(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{20, 10, 30} {
		rbts.Insert(tree, v, "")
	}

	m, ok := rbts.Min(tree)
	require.True(t, ok)
	assert.Equal(t, 10, m.Key())
}

func TestMax(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{20, 10, 30} {
		rbts.Insert(tree, v, "")
	}

	m, ok := rbts.Max(tree)
	require.True(t, ok)
	assert.Equal(t, 30, m.Key())
}

func TestCeiling(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30} {
		rbts.Insert(tree, v, "")
	}

	n, ok := rbts.Ceiling(tree, 5)
	require.True(t, ok)
	assert.Equal(t, 10, n.Key())

	n, ok = rbts.Ceiling(tree, 20)
	require.True(t, ok)
	assert.Equal(t, 20, n.Key())

	_, ok = rbts.Ceiling(tree, 40)
	assert.False(t, ok)
}

func TestFloor(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30} {
		rbts.Insert(tree, v, "")
	}

	_, ok := rbts.Floor(tree, 5)
	assert.False(t, ok)

	n, ok := rbts.Floor(tree, 15)
	require.True(t, ok)
	assert.Equal(t, 10, n.Key())

	n, ok = rbts.Floor(tree, 20)
	require.True(t, ok)
	assert.Equal(t, 20, n.Key())
}

func TestHigher(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30} {
		rbts.Insert(tree, v, "")
	}

	n, ok := rbts.Higher(tree, 10)
	require.True(t, ok)
	assert.Equal(t, 20, n.Key())

	_, ok = rbts.Higher(tree, 35)
	assert.False(t, ok)
}

func TestLower(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30} {
		rbts.Insert(tree, v, "")
	}

	_, ok := rbts.Lower(tree, 5)
	assert.False(t, ok)

	n, ok := rbts.Lower(tree, 15)
	require.True(t, ok)
	assert.Equal(t, 10, n.Key())

	n, ok = rbts.Lower(tree, 30)
	require.True(t, ok)
	assert.Equal(t, 20, n.Key())
}

func TestPredecessor(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30, 40, 50} {
		rbts.Insert(tree, v, "")
	}

	n, found := rbts.Search(tree, 30)
	require.True(t, found)
	pred, ok := rbts.Predecessor(n)
	require.True(t, ok)
	assert.Equal(t, 20, pred.Key())
}

func TestSuccessor(t *testing.T) {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30, 40, 50} {
		rbts.Insert(tree, v, "")
	}

	n, found := rbts.Search(tree, 30)
	require.True(t, found)
	succ, ok := rbts.Successor(n)
	require.True(t, ok)
	assert.Equal(t, 40, succ.Key())
}

func TestInOrder(t *testing.T) {
	tree := rbts.New[int, string]()
	values := []int{20, 10, 30, 5, 15, 25, 35}
	for _, v := range values {
		rbts.Insert(tree, v, "")
	}

	prev := -1
	for n := range rbts.InOrder(tree) {
		if prev != -1 {
			assert.Less(t, prev, n.Key(), "InOrder traversal is not sorted")
		}
		prev = n.Key()
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
	assert.Equal(t, len(expected), len(collected), "Expected range length")

	for i, v := range expected {
		assert.Equal(t, v, collected[i], "Expected value at position")
	}
}

func TestRank(t *testing.T) {
	tree := rbts.New[int, string]()
	values := []int{10, 20, 30, 40, 50}
	for _, v := range values {
		rbts.Insert(tree, v, "")
	}

	assert.Equal(t, 2, rbts.Rank(tree, 25))
	assert.Equal(t, 0, rbts.Rank(tree, 10))
	assert.Equal(t, 5, rbts.Rank(tree, 60))
}

func TestKth(t *testing.T) {
	tree := rbts.New[int, string]()
	values := []int{10, 20, 30, 40, 50}
	for _, v := range values {
		rbts.Insert(tree, v, "")
	}

	n, ok := rbts.Kth(tree, 0)
	require.True(t, ok)
	assert.Equal(t, 10, n.Key())

	n, ok = rbts.Kth(tree, 3)
	require.True(t, ok)
	assert.Equal(t, 40, n.Key())

	_, ok = rbts.Kth(tree, 5)
	assert.False(t, ok)
}

func ExampleNew() {
	tree := rbts.New[int, string]()
	fmt.Println(rbts.Len(tree))
	// Output: 0
}

func ExampleLen() {
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
	min, ok := rbts.Min(tree)
	if ok {
		fmt.Println(min.Key())
	}
	// Output: 10
}

func ExampleMax() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 20, "")
	rbts.Insert(tree, 30, "")
	max, ok := rbts.Max(tree)
	if ok {
		fmt.Println(max.Key())
	}
	// Output: 30
}

func ExampleCeiling() {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30} {
		rbts.Insert(tree, v, "")
	}

	n, _ := rbts.Ceiling(tree, 15)
	fmt.Println(n.Key())
	n, _ = rbts.Ceiling(tree, 20)
	fmt.Println(n.Key())
	// Output:
	// 20
	// 20
}

func ExampleFloor() {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30} {
		rbts.Insert(tree, v, "")
	}

	n, _ := rbts.Floor(tree, 25)
	fmt.Println(n.Key())
	n, _ = rbts.Floor(tree, 20)
	fmt.Println(n.Key())
	// Output:
	// 20
	// 20
}

func ExampleHigher() {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30} {
		rbts.Insert(tree, v, "")
	}

	n, _ := rbts.Higher(tree, 15)
	fmt.Println(n.Key())
	n, _ = rbts.Higher(tree, 20)
	fmt.Println(n.Key())
	// Output:
	// 20
	// 30
}

func ExampleLower() {
	tree := rbts.New[int, string]()
	for _, v := range []int{10, 20, 30} {
		rbts.Insert(tree, v, "")
	}

	n, _ := rbts.Lower(tree, 15)
	fmt.Println(n.Key())
	n, _ = rbts.Lower(tree, 20)
	fmt.Println(n.Key())
	// Output:
	// 10
	// 10
}

func ExamplePredecessor() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 20, "")
	rbts.Insert(tree, 10, "")
	node, _ := rbts.Search(tree, 20)
	pred, ok := rbts.Predecessor(node)
	if ok {
		fmt.Println(pred.Key())
	}
	// Output: 10
}

func ExampleSuccessor() {
	tree := rbts.New[int, string]()
	rbts.Insert(tree, 20, "")
	rbts.Insert(tree, 30, "")
	node, _ := rbts.Search(tree, 20)
	succ, ok := rbts.Successor(node)
	if ok {
		fmt.Println(succ.Key())
	}
	// Output: 30
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

func BenchmarkInsertRandom(b *testing.B) {
	r := rand.New(rand.NewSource(42))
	for i := 0; i < b.N; i++ {
		tree := rbts.New[int, string]()
		for j := 0; j < 1000; j++ {
			key := r.Intn(1_000_000)
			rbts.Insert(tree, key, "value")
		}
	}
}

func BenchmarkInsertSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tree := rbts.New[int, string]()
		for j := 0; j < 1000; j++ {
			rbts.Insert(tree, j, "value")
		}
	}
}

func BenchmarkSearchHit(b *testing.B) {
	tree := rbts.New[int, string]()
	for i := 0; i < 1000; i++ {
		rbts.Insert(tree, i, "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbts.Search(tree, i%1000)
	}
}

func BenchmarkSearchMiss(b *testing.B) {
	tree := rbts.New[int, string]()
	for i := 0; i < 1000; i++ {
		rbts.Insert(tree, i, "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbts.Search(tree, 1_000_000+i)
	}
}

func BenchmarkDeleteRandom(b *testing.B) {
	r := rand.New(rand.NewSource(42))
	tree := rbts.New[int, string]()
	keys := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = r.Intn(1_000_000)
		rbts.Insert(tree, keys[i], "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rbts.Delete(tree, keys[i%1000])
	}
}
