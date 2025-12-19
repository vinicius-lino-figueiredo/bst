// Package comparer TODO
package comparer

import (
	"cmp"

	"github.com/vinicius-lino-figueiredo/bst/domain"
)

// NewComparer TODO
func NewComparer[K cmp.Ordered, V comparable]() domain.Comparer[K, V] {
	return Comparer[K, V]{}
}

// Comparer TODO
type Comparer[K cmp.Ordered, V comparable] struct{}

// CompareKeys implements bst.Comparer.
func (c Comparer[K, V]) CompareKeys(a K, b K) (int, error) {
	return cmp.Compare(a, b), nil
}

// CompareValues implements bst.Comparer.
func (c Comparer[K, V]) CompareValues(a V, b V) (bool, error) {
	return a == b, nil
}
