// Package bst implements a bst based on the js seald/node-binary-search-tree.
// This version is not optimized for golang, but is rather a direct translation
// of the original code.
package bst

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/vinicius-lino-figueiredo/bst/customutils"
)

// Options is an object used to pass the arguments to create a
// [BinarySearchTree].
type Options struct {
	Unique             bool
	Parent             *BinarySearchTree
	Key                any
	Value              any
	CompareKeys        func(a any, b any) int
	CheckValueEquality func(a, b any) bool
}

// BinarySearchTree represents a simple binary search tree.
type BinarySearchTree struct {
	right              *BinarySearchTree
	left               *BinarySearchTree
	parent             *BinarySearchTree
	key                customutils.Comparer
	data               []any
	unique             bool
	compareKeys        func(a any, b any) int
	checkValueEquality func(a, b any) bool
}

// NewBinarySearchTree creates a new instance of [BinarySearchTree].
func NewBinarySearchTree(options Options) *BinarySearchTree {
	var b BinarySearchTree
	b.left = nil
	b.right = nil
	if options.Parent != nil {
		b.parent = options.Parent
	}
	if options.Key != nil {
		b.key = customutils.NewCaster(options.Key)
	}
	if options.Value != nil {
		b.data = []any{options.Value}
	}
	b.unique = options.Unique || false

	b.compareKeys = options.CompareKeys
	if options.CompareKeys == nil {
		b.compareKeys = customutils.DefaultCompareKeysFunction
	}
	b.checkValueEquality = options.CheckValueEquality
	if options.CheckValueEquality == nil {
		b.checkValueEquality = customutils.DefaultCheckValueEquality
	}
	return &b
}

// Key returns the node key.
func (b *BinarySearchTree) Key() any {
	return customutils.Retrieve(b.key)
}

// SetKey sets the value for the node key.
func (b *BinarySearchTree) SetKey(key any) {
	b.key = customutils.NewCaster(key)
}

// Data returns the node data.
func (b *BinarySearchTree) Data() []any {
	return b.data
}

// Parent returns the node parent, if any.
func (b *BinarySearchTree) Parent() *BinarySearchTree {
	return b.parent
}

// Unique returns true if the tree has unique constraint.
func (b *BinarySearchTree) Unique() bool {
	return b.unique
}

// Right returns the right node, if any.
func (b *BinarySearchTree) Right() *BinarySearchTree {
	return b.right
}

// Left returns the left node, if any.
func (b *BinarySearchTree) Left() *BinarySearchTree {
	return b.left
}

// GetMaxKeyDescendant returns the descendant with max key.
func (b *BinarySearchTree) GetMaxKeyDescendant() *BinarySearchTree {
	if b.right != nil {
		return b.right.GetMaxKeyDescendant()
	}
	return b // removed unnecessary else
}

// GetMaxKey returns the maximum key.
func (b *BinarySearchTree) GetMaxKey() any {
	key := b.GetMaxKeyDescendant().key
	return customutils.Retrieve(key)
}

// GetMinKeyDescendant returns the descendant with min key.
func (b *BinarySearchTree) GetMinKeyDescendant() *BinarySearchTree {
	if b.left != nil {
		return b.left.GetMinKeyDescendant()
	}
	return b
}

// GetMinKey returns the minimum key.
func (b *BinarySearchTree) GetMinKey() any {
	key := b.GetMinKeyDescendant().key
	return customutils.Retrieve(key)
}

// CheckAllNodesFullfillCondition checks that all nodes (incl. leaves) fullfil
// condition given by fn
// test is a function passed every (key, data) and which returns error if the
// condition is not met.
func (b *BinarySearchTree) CheckAllNodesFullfillCondition(test func(key any, data any) error) error {
	if b.key == nil {
		return nil
	}
	var err error
	if err = test(customutils.Retrieve(b.key), b.data); err != nil {
		return err
	}
	if b.left != nil {
		if err = b.left.CheckAllNodesFullfillCondition(test); err != nil {
			return err
		}
	}
	if b.right != nil {
		if err = b.right.CheckAllNodesFullfillCondition(test); err != nil {
			return err
		}
	}
	return nil
}

// CheckNodeOrdering checks that the core BST properties on node ordering are
// verified. Returns error if they aren't.
func (b *BinarySearchTree) CheckNodeOrdering() error {
	if b.key == nil {
		return nil
	}

	var err error
	if b.left != nil {
		err = b.left.CheckAllNodesFullfillCondition(func(k any, _ any) error {
			if b.callCompareKeys(k, b.key) >= 0 {
				return fmt.Errorf("Tree with root %v is not a binary search tree", b.key)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = b.left.CheckNodeOrdering()
		if err != nil {
			return err
		}
	}

	if b.right != nil {
		err = b.right.CheckAllNodesFullfillCondition(func(k any, _ any) error {
			if b.callCompareKeys(k, b.key) <= 0 {
				return fmt.Errorf("Tree with root %v is not a binary search tree", b.key)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = b.right.CheckNodeOrdering()
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckInternalPointers checks that all pointers are coherent in this tree.
func (b *BinarySearchTree) CheckInternalPointers() error {
	var err error
	if b.left != nil {
		if b.left.parent != b {
			return fmt.Errorf(`Parent pointer broken for key %v`, b.Key())
		}
		err = b.left.CheckInternalPointers()
		if err != nil {
			return err
		}
	}

	if b.right != nil {
		if b.right.parent != b {
			return fmt.Errorf(`Parent pointer broken for key %v`, b.Key())
		}
		err = b.right.CheckInternalPointers()
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckIsBST checks that a tree is a BST as defined here (node ordering and
// pointer references).
func (b *BinarySearchTree) CheckIsBST() error {
	var err error
	if err = b.CheckNodeOrdering(); err != nil {
		return err
	}
	if err = b.CheckInternalPointers(); err != nil {
		return err
	}
	if b.parent != nil {
		return errors.New("The root shouldn't have a parent")
	}
	return nil
}

// GetNumberOfKeys gets the number of keys inserted.
func (b *BinarySearchTree) GetNumberOfKeys() int {
	res := 0

	if b.key == nil {
		return 0
	}

	res = 1
	if b.left != nil {
		res += b.left.GetNumberOfKeys()
	}
	if b.right != nil {
		res += b.right.GetNumberOfKeys()
	}

	return res
}

// CreateSimilar creates a BST similar (i.e. same options except for key and
// value) to the current one.
func (b *BinarySearchTree) CreateSimilar(options Options) *BinarySearchTree {
	options.Unique = b.unique
	options.CompareKeys = b.callCompareKeys
	options.CheckValueEquality = b.checkValueEquality

	return NewBinarySearchTree(options)
}

// CreateLeftChild creates the left child of this BST and returns it.
func (b *BinarySearchTree) CreateLeftChild(options Options) *BinarySearchTree {
	leftChild := b.CreateSimilar(options)
	leftChild.parent = b
	b.left = leftChild

	return leftChild
}

// CreateRightChild creates the right child of this BST and returns it.
func (b *BinarySearchTree) CreateRightChild(options Options) *BinarySearchTree {
	rightChild := b.CreateSimilar(options)
	rightChild.parent = b
	b.right = rightChild

	return rightChild
}

// Insert inserts a new element.
func (b *BinarySearchTree) Insert(key any, value any) error {
	// Empty tree, insert as root
	k := customutils.NewCaster(key)
	if b.key == nil {
		b.key = k
		b.data = append(b.data, value)
		return nil
	}

	// Same key as root
	if b.callCompareKeys(b.key, key) == 0 {
		if b.unique {
			return &ErrViolated{key: key}
		}
		b.data = append(b.data, value)
		return nil
	}

	var err error
	if b.callCompareKeys(k, b.key) < 0 {
		// Insert in left subtree
		if b.left != nil {
			if err = b.left.Insert(key, value); err != nil {
				return err
			}
		} else {
			b.CreateLeftChild(Options{Key: key, Value: value})
		}
	} else {
		// Insert in right subtree
		if b.right != nil {
			if err = b.right.Insert(key, value); err != nil {
				return err
			}
		} else {
			b.CreateRightChild(Options{Key: key, Value: value})
		}
	}
	return nil
}

// Search searches for all data corresponding to a key.
func (b *BinarySearchTree) Search(key any) []any {
	if b.callCompareKeys(b.key, key) == 0 {
		return b.data
	}

	k := customutils.NewCaster(key)
	if b.callCompareKeys(k, b.key) < 0 {
		if b.left != nil {
			return b.left.Search(key)
		}
		return []any{}
	}
	if b.right != nil {
		return b.right.Search(key)
	}
	return []any{}
}

// GetLowerBoundMatcher returns a function that tells whether a given key
// matches a lower bound.
func (b *BinarySearchTree) GetLowerBoundMatcher(query map[string]any) func(customutils.Comparer) bool {
	// No lower bound
	if query["$gd"] == nil && query["$gte"] == nil {
		return func(customutils.Comparer) bool { return true }
	}

	if query["$gt"] != nil && query["$gte"] != nil {
		queryGte := customutils.NewCaster(query["$gte"])
		if b.callCompareKeys(queryGte, query["$gt"]) == 0 {
			return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$gt"]) > 0 }
		}

		if b.callCompareKeys(queryGte, query["$gt"]) > 0 {
			return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$gte"]) >= 0 }
		}
		return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$gt"]) > 0 }
	}

	if query["$gt"] != nil {
		return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$gt"]) > 0 }
	}
	return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$gte"]) >= 0 }
}

// GetUpperBoundMatcher returns a function that tells whether a given key
// matches an upper bound.
func (b *BinarySearchTree) GetUpperBoundMatcher(query map[string]any) func(customutils.Comparer) bool {
	// No upper bound
	if query["$lt"] == nil && query["$lte"] == nil {
		return func(customutils.Comparer) bool { return true }
	}

	if query["$lt"] != nil && query["$lte"] != nil {
		queryLte := customutils.NewCaster(query["$lte"])
		if b.callCompareKeys(queryLte, query["$lt"]) == 0 {
			return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$lt"]) < 0 }
		}

		if b.callCompareKeys(queryLte, query["$lt"]) < 0 {
			return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$lte"]) <= 0 }
		}
		return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$lt"]) < 0 }
	}

	if query["$lt"] != nil {
		return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$lt"]) < 0 }
	}
	return func(key customutils.Comparer) bool { return b.callCompareKeys(key, query["$lte"]) <= 0 }
}

// BetweenBounds gets all data for a key between bounds and returns it in key
// order. Param query is a Mongo-style query where keys are "$lt", "$lte", "$gt"
// or "$gte" (other keys are not considered). Params lbm and ubm are matching
// functions calculated at the first recursive step.
func (b *BinarySearchTree) BetweenBounds(query map[string]any, lbm, ubm func(customutils.Comparer) bool) []any {
	res := []any{}

	if b.key == nil {
		return []any{} // Empty tree
	}

	if lbm == nil {
		lbm = b.GetLowerBoundMatcher(query)
	}
	if ubm == nil {
		ubm = b.GetUpperBoundMatcher(query)
	}

	if lbm(b.key) && b.left != nil {
		res = append(res, b.left.BetweenBounds(query, lbm, ubm)...)
	}
	if lbm(b.key) && ubm(b.key) {
		res = append(res, b.data...)
	}
	if ubm(b.key) && b.right != nil {
		res = append(res, b.right.BetweenBounds(query, lbm, ubm)...)
	}

	return res
}

// DeleteIfLeaf deletes the current node if it is a leaf. Returns true if it was
// deleted.
func (b *BinarySearchTree) DeleteIfLeaf() bool {
	if b.left != nil || b.right != nil {
		return false
	}

	// The leaf is itself a root
	if b.parent == nil {
		b.key = nil
		b.data = []any{}
		return true
	}

	if b.parent.left == b {
		b.parent.left = nil
	} else {
		b.parent.right = nil
	}

	return true
}

// DeleteIfOnlyOneChild deletes the current node if it has only one child
// Returns true if it was deleted.
func (b *BinarySearchTree) DeleteIfOnlyOneChild() bool {
	var child *BinarySearchTree

	if b.left != nil && b.right == nil {
		child = b.left
	}
	if b.left == nil && b.right != nil {
		child = b.right
	}
	if child == nil {
		return false
	}

	// Root
	if b.parent == nil {
		b.key = child.key
		b.data = child.data

		b.left = nil
		if child.left != nil {
			b.left = child.left
			child.left.parent = b
		}

		b.right = nil
		if child.right != nil {
			b.right = child.right
			child.right.parent = b
		}

		return true
	}

	if b.parent.left == b {
		b.parent.left = child
		child.parent = b.parent
	} else {
		b.parent.right = child
		child.parent = b.parent
	}

	return true
}

// Delete a key or just a value. Param value is Optional. If not set, the whole
// key is deleted. If set, only this value is deleted.
func (b *BinarySearchTree) Delete(key any, value any) {
	newData := []any{}
	var replaceWith *BinarySearchTree

	if b.key == nil {
		return
	}

	k := customutils.NewCaster(key)
	if b.callCompareKeys(k, b.key) < 0 {
		if b.left != nil {
			b.left.Delete(key, value)
		}
		return
	}

	if b.callCompareKeys(k, b.key) > 0 {
		if b.right != nil {
			b.right.Delete(key, value)
		}
		return
	}

	if !(b.callCompareKeys(k, b.key) == 0) {
		return
	}

	// Delete only a value
	if len(b.data) > 1 && value != nil {
		for _, d := range b.data {
			if !b.checkValueEquality(d, value) {
				newData = append(newData, d)
			}
		}
		b.data = newData
		return
	}

	// Delete the whole node
	if b.DeleteIfLeaf() {
		return
	}

	if b.DeleteIfOnlyOneChild() {
		return
	}

	// We are in the case where the node to delete has two children
	if rand.Float64() >= 0.5 { // Randomize replacement to avoid unbalancing the tree too much
		// Use the in-order predecessor
		replaceWith = b.left.GetMaxKeyDescendant()

		b.key = replaceWith.key
		b.data = replaceWith.data

		if b == replaceWith.parent { // Special case
			b.left = replaceWith.left
			if replaceWith.left != nil {
				replaceWith.left.parent = replaceWith.parent
			}
		} else {
			replaceWith.parent.right = replaceWith.left
			if replaceWith.left != nil {
				replaceWith.left.parent = replaceWith.parent
			}
		}
	} else {
		// Use the in-order successor
		replaceWith = b.right.GetMinKeyDescendant()

		b.key = replaceWith.key
		b.data = replaceWith.data

		if b == replaceWith.parent { // Special case
			b.right = replaceWith.right
			if replaceWith.right != nil {
				replaceWith.right.parent = replaceWith.parent
			}
		} else {
			replaceWith.parent.left = replaceWith.right
			if replaceWith.right != nil {
				replaceWith.right.parent = replaceWith.parent
			}
		}
	}
}

// ExecuteOnEveryNode executes a function on every node of the tree, in key
// order.
func (b *BinarySearchTree) ExecuteOnEveryNode(fn func(*BinarySearchTree)) {
	if b.left != nil {
		b.left.ExecuteOnEveryNode(fn)
	}
	fn(b)
	if b.right != nil {
		b.right.ExecuteOnEveryNode(fn)
	}
}

// String returns the data as a string.
func (b *BinarySearchTree) String() string {
	return b.format(false, "")
}

// format formats the data as string
func (b *BinarySearchTree) format(printData bool, spacing string) string {
	var res string
	res += fmt.Sprintf("%s* %s", spacing, b.key)
	if printData {
		res += fmt.Sprintf("%s* %s", spacing, b.data)
	}

	if b.left == nil && b.right == nil {
		return res
	}

	if b.left != nil {
		res += b.left.format(printData, fmt.Sprintf("%s  ", spacing))
	} else {
		res += fmt.Sprintf("%s  *", spacing)
	}

	if b.right != nil {
		res += b.right.format(printData, fmt.Sprintf("%s  ", spacing))
	} else {
		res += fmt.Sprintf("%s  *", spacing)
	}
	return res
}

func (b *BinarySearchTree) callCompareKeys(a, other any) int {
	return b.compareKeys(customutils.Retrieve(a), customutils.Retrieve(other))
}

// ================================
// Methods used to test the tree
// ================================

// ============================================
// Methods used to actually work on the tree
// ============================================
