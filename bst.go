// Package bst TODO
package bst

import (
	"fmt"
	"iter"
)

// ErrUniqueViolated TODO
type ErrUniqueViolated struct {
	Key any
}

func (e ErrUniqueViolated) Error() string {
	return fmt.Sprintf("constraint violated: %v is not unique", e.Key)
}

// Bound TODO
type Bound[K any] struct {
	Value        K
	IncludeEqual bool
}

// Query TODO
type Query[K any] struct {
	GreaterThan *Bound[K]
	LowerThan   *Bound[K]
}

// BST TODO
type BST[K any, V any] interface {
	Insert(key K, value V) error

	Search(key K) (*Node[K, V], error)
	Query(query Query[K]) iter.Seq2[V, error]
	GetMax() *Node[K, V]
	GetMin() *Node[K, V]
	GetNumberOfKeys() int
	GetAll() iter.Seq[V]

	Update(key K, old V, nw V) error

	Delete(key K, value *V) error
}

// Node TODO
type Node[K any, V any] struct {
	Values  []V
	Key     K
	Lower   *Node[K, V]
	Greater *Node[K, V]
	Parent  *Node[K, V]
}

// Comparer TODO
type Comparer[K any, V any] interface {
	CompareKeys(a K, b K) (int, error)
	CompareValues(a V, b V) (bool, error)
}
