// Package avl TODO
package avl

import (
	"iter"
	"slices"
	"sync"

	"github.com/vinicius-lino-figueiredo/bst"
)

// NewBST TODO
func NewBST[K any, V any](unique bool, creationSize int, comparer bst.Comparer[K, V]) bst.BST[K, V] {
	if unique {
		creationSize = 1
	} else if creationSize <= 0 {
		creationSize = 8
	}
	return &Root[K, V]{
		unique:       unique,
		creationSize: creationSize,
		comparer:     comparer,
		nodePool:     sync.Pool{New: func() any { return &node[K, V]{} }},
		Node: node[K, V]{
			values: make([]V, 0, creationSize),
		},
	}
}

// Root TODO.
type Root[K any, V any] struct {
	Node         node[K, V]
	initialized  bool
	nodeCount    int
	unique       bool
	creationSize int
	nodePool     sync.Pool
	comparer     bst.Comparer[K, V]
}

type node[K any, V any] struct {
	height  int
	key     K
	values  []V
	parent  *node[K, V]
	lower   *node[K, V]
	greater *node[K, V]
}

// Greater implements [bst.Node].
func (n *node[K, V]) Greater() bst.Node[K, V] {
	return n.greater
}

// Key implements [bst.Node].
func (n *node[K, V]) Key() K {
	return n.key
}

// Lower implements [bst.Node].
func (n *node[K, V]) Lower() bst.Node[K, V] {
	return n.lower
}

// Parent implements [bst.Node].
func (n *node[K, V]) Parent() bst.Node[K, V] {
	return n.parent
}

// Values implements [bst.Node].
func (n *node[K, V]) Values() []V {
	return n.values
}

// Insert implements bst.BST.
func (r *Root[K, V]) Insert(key K, value V) error {
	if !r.initialized {
		r.Node.height = 1
		r.Node.key = key
		r.initialized = true
		r.Node.values = append(r.Node.values, value)
		r.nodeCount++
		return nil
	}
	node := &r.Node
Loop:
	for {
		comparison, err := r.comparer.CompareKeys(key, node.key)
		if err != nil {
			return err
		}
		switch {
		case comparison > 0:
			if node.greater == nil {
				node.greater = r.createEmptyNode(key, node)
				r.nodeCount++
				r.updateHeight(node)
				node = node.greater
				break Loop
			}
			node = node.greater
		case comparison < 0:
			if node.lower == nil {
				node.lower = r.createEmptyNode(key, node)
				r.nodeCount++
				r.updateHeight(node)
				node = node.lower
				break Loop
			}
			node = node.lower
		default:
			if r.unique {
				return bst.ErrUniqueViolated{Key: key}
			}
			break Loop
		}
	}
	node.values = append(node.values, value)

	r.balance(node.parent)

	return nil
}

func (r *Root[K, V]) balance(node *node[K, V]) {
	for node != nil {
		switch r.balanceFactor(node) {
		case -2:
			if r.balanceFactor(node.greater) > 0 { // Right-Left
				r.rotateRight(node.greater)
			}
			r.rotateLeft(node)
		default:
		case +2:
			if r.balanceFactor(node.lower) < 0 { // Left-Right
				r.rotateLeft(node.lower)
			}
			r.rotateRight(node)
		}

		node.height = max(r.height(node.lower), r.height(node.greater)) + 1

		node = node.parent
	}
}

func (r *Root[K, V]) updateHeight(node *node[K, V]) {
	for node != nil {
		node.height = max(r.height(node.lower), r.height(node.greater)) + 1
		node = node.parent
	}
}

func (r *Root[K, V]) height(node *node[K, V]) int {
	if node != nil {
		return node.height
	}
	return 0
}

//	func (r *Root[K, V]) rotate(node *node[K, V], bf int) *node[K, V] {
//		switch bf {
//		case -1, 0, +1:
//			return node
//		case +2:
//			return r.rotateLeftleft(node)
//		case -2:
//			if node.lower.greater == nil {
//				return r.rotateRight(node)
//			}
//			return r.rotateLeftRight(node)
//		default:
//		}
//		return node
//	}
//
//	func (r *Root[K, V]) rotateRightLeft(node *node[K, V]) *node[K, V] {
//		node = r.rotateRight(node.greater)
//		node = r.rotateLeftleft(node.parent)
//		return node
//	}
//
//	func (r *Root[K, V]) rotateLeftRight(node *node[K, V]) *node[K, V] {
//		return node
//	}
func (r *Root[K, V]) rotateRight(n *node[K, V]) *node[K, V] {
	newGreater := n.lower

	r.swapData(newGreater, n)
	n.greater, n.lower = n.lower, n.greater

	n.lower, newGreater.lower, newGreater.greater = newGreater.lower, newGreater.greater, n.lower

	newGreater.parent = n
	if n.lower != nil {
		n.lower.parent = n
	}
	if newGreater.greater != nil {
		newGreater.greater.parent = newGreater
	}

	r.updateHeight(n)
	r.updateHeight(newGreater)
	return nil
}

func (r *Root[K, V]) rotateLeft(n *node[K, V]) *node[K, V] {
	newLower := n.greater

	r.swapData(newLower, n)
	n.lower, n.greater = n.greater, n.lower

	n.greater, newLower.greater, newLower.lower = newLower.greater, newLower.lower, n.greater

	newLower.parent = n
	if n.greater != nil {
		n.greater.parent = n
	}
	if newLower.lower != nil {
		newLower.lower.parent = newLower
	}

	r.updateHeight(n)
	r.updateHeight(newLower)
	return nil
}

func (r *Root[K, V]) swapData(a *node[K, V], b *node[K, V]) {
	a.key, b.key = b.key, a.key
	a.values, b.values = b.values, a.values
	a.height, b.height = b.height, a.height
}

// func (r *root[k, v]) rotateleftleft(node *node[k, v]) *node[k, v] {
// 	// swapping values
// 	node.key, node.greater.key = node.greater.key, node.key
// 	node.values, node.greater.values = node.greater.values, node.values
// 	node.gh, node.greater.gh = node.greater.gh, node.gh
// 	node.lh, node.greater.lh = node.lh+1, node.lh
//
// 	// swapping positions
// 	node.lower, node.greater = node.greater, node.lower
// 	// rotating
// 	node.lower.lower, node.greater = node.greater, node.lower.greater
// 	// fixing references
// 	node.greater.parent, node.lower.greater = node, nil
// 	return node
// }

func (r *Root[K, V]) createEmptyNode(key K, parent *node[K, V]) *node[K, V] {
	node := r.nodePool.Get().(*node[K, V])
	node.key = key
	node.values = make([]V, 0, r.creationSize)
	node.parent = parent
	node.height = 1
	return node
}

// Search implements bst.BST.
func (r *Root[K, V]) Search(key K) (bst.Node[K, V], error) {
	res, err := r.search(key)
	if err != nil || res == nil {
		return nil, err
	}
	return res, nil
}

func (r *Root[K, V]) search(key K) (*node[K, V], error) {
	if !r.initialized {
		return nil, nil
	}
	node := &r.Node
	for {
		comparison, err := r.comparer.CompareKeys(key, node.key)
		if err != nil {
			return nil, err
		}
		switch {
		case comparison > 0 && node.greater != nil:
			node = node.greater
			continue
		case comparison < 0 && node.lower != nil:
			node = node.lower
			continue
		case comparison == 0:
			return node, nil
		}
		return nil, nil
	}
}

// Query implements bst.BST.
func (r *Root[K, V]) Query(query bst.Query[K]) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		switch {
		case query.GreaterThan != nil:
			switch query.LowerThan {
			case nil:
				_ = r.queryGreater(&r.Node, query.GreaterThan, yield)
			default:
				_ = r.doubleQuery(&r.Node, query, yield)
			}
		case query.LowerThan != nil:
			_ = r.queryLower(&r.Node, query.LowerThan, yield)
		default:
		}
	}
}

func (r *Root[K, V]) doubleQuery(node *node[K, V], query bst.Query[K], yield func(V, error) bool) bool {
	ltComp, err := r.comparer.CompareKeys(node.key, query.LowerThan.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}

	switch {
	case ltComp > 0:
		return r.treatAboveMax(node, query, yield)
	case ltComp < 0:
		return r.treatBelowMax(node, query, yield)
	default:
		return r.treatEqualMax(node, query, yield)
	}
}

func (r *Root[K, V]) treatAboveMax(node *node[K, V], query bst.Query[K], yield func(V, error) bool) bool {
	if node.lower == nil {
		return true
	}
	gtComp, err := r.comparer.CompareKeys(node.key, query.GreaterThan.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}

	if gtComp < 0 {
		return true
	}

	return r.doubleQuery(node.lower, query, yield)
}

func (r *Root[K, V]) treatBelowMax(node *node[K, V], query bst.Query[K], yield func(V, error) bool) bool {
	gtComp, err := r.comparer.CompareKeys(node.key, query.GreaterThan.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}
	switch {
	case gtComp < 0: // node lower than min
		if r.Node.greater != nil {
			return r.doubleQuery(node.greater, query, yield)
		}
	case gtComp == 0: // node equal to min
		if query.GreaterThan.IncludeEqual && !r.yieldValues(node, yield) {
			return false
		}
	default:
		if node.lower != nil && !r.queryGreater(node.lower, query.GreaterThan, yield) {
			return false
		}
		if !r.yieldValues(node, yield) {
			return false
		}
	}
	if node.greater != nil {
		return r.doubleQuery(node.greater, query, yield)
	}
	return true
}

func (r *Root[K, V]) treatEqualMax(node *node[K, V], query bst.Query[K], yield func(V, error) bool) bool {
	gtComp, err := r.comparer.CompareKeys(node.key, query.GreaterThan.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}
	switch {
	case gtComp > 0:
		if node.lower != nil && !r.queryGreater(node.lower, query.GreaterThan, yield) {
			return false
		}
		if query.LowerThan.IncludeEqual {
			return r.yieldValues(node, yield)
		}
	case gtComp < 0:
	default:
		if query.GreaterThan.IncludeEqual && query.LowerThan.IncludeEqual {
			return r.yieldValues(node, yield)
		}
	}
	return true
}

func (r *Root[K, V]) queryGreater(node *node[K, V], bound *bst.Bound[K], yield func(V, error) bool) bool {
	comp, err := r.comparer.CompareKeys(node.key, bound.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}
	switch {
	case comp > 0:
		if node.lower != nil && !r.queryGreater(node.lower, bound, yield) {
			return false
		}
		if !r.yieldValues(node, yield) {
			return false
		}
	case comp < 0:
	default:
		if bound.IncludeEqual && !r.yieldValues(node, yield) {
			return false
		}
	}
	if node.greater != nil {
		return r.queryGreater(node.greater, bound, yield)
	}
	return true
}

func (r *Root[K, V]) queryLower(node *node[K, V], bound *bst.Bound[K], yield func(V, error) bool) bool {
	comp, err := r.comparer.CompareKeys(node.key, bound.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}

	switch {
	case comp < 0:
		if node.lower != nil && !r.queryLower(node.lower, bound, yield) {
			return false
		}
		if !r.yieldValues(node, yield) {
			return false
		}
		if node.greater != nil && !r.queryLower(node.greater, bound, yield) {
			return false
		}
	case comp > 0:
		if node.lower != nil {
			return r.queryLower(node.lower, bound, yield)
		}
	default:
		if node.lower != nil && !r.queryLower(node.lower, bound, yield) {
			return false
		}
		if bound.IncludeEqual && !r.yieldValues(node, yield) {
			return false
		}
	}
	return true
}

func (r *Root[K, V]) yieldValues(node *node[K, V], yield func(V, error) bool) bool {
	for _, v := range node.values {
		if !yield(v, nil) {
			return false
		}
	}
	return true
}

// Delete implements bst.BST.
func (r *Root[K, V]) Delete(key K, value *V) error {
	if !r.initialized {
		return nil
	}
	node, err := r.search(key)
	if err != nil || node == nil {
		return err
	}
	if value != nil {
		if err = r.deleteValue(node, value); err != nil || len(node.values) > 0 {
			return err
		}
	}
	r.nodeCount--

Switch:
	switch {
	case node.lower != nil:
		if node.greater != nil {
			r.deleteDoubleChildrenNode(node)
			break
		}
		r.takePlace(node, node.lower)
	case node.greater != nil:
		r.takePlace(node, node.greater)
	default:
		node.values = node.values[:0]
		switch node {
		case &r.Node:
			r.initialized = false
			break Switch
		case node.parent.lower:
			node.parent.lower = nil
		default:
			node.parent.greater = nil
		}
		node.height = 0

		r.updateHeight(node.parent)
		r.balance(node.parent)

		node.parent = nil
		r.nodePool.Put(node)
		return nil
	}
	r.updateHeight(node)
	r.balance(node)

	return nil
}

func (r *Root[K, V]) balanceFactor(node *node[K, V]) int {
	if node == nil {
		return 0
	}
	return r.height(node.lower) - r.height(node.greater)
}

func (r *Root[K, V]) takePlace(node, victim *node[K, V]) {
	node.key, node.values = victim.key, victim.values
	node.greater, node.lower = victim.greater, victim.lower
	if node.greater != nil {
		node.greater.parent = node
	}
	if node.lower != nil {
		node.lower.parent = node
	}

	victim.parent, victim.greater, victim.lower = nil, nil, nil
	victim.values = victim.values[:0]
	r.nodePool.Put(victim)
}

func (r *Root[K, V]) deleteDoubleChildrenNode(n *node[K, V]) {
	var closestNode *node[K, V]
	switch n.lower.height - n.greater.height {
	case 1:
		closestNode = r.getMax(n.lower)

		// cloning closest value
		n.key = closestNode.key
		n.values = closestNode.values
		if closestNode != n.lower {
			closestNode.parent.greater = closestNode.lower
			if closestNode.lower != nil {
				closestNode.lower.parent = n.lower
			}
		} else {
			n.lower = closestNode.lower
			if n.lower != nil {
				n.lower.parent = n
			}
		}
	default:
		closestNode = r.getMin(n.greater)

		// cloning closest value
		n.key = closestNode.key
		n.values = closestNode.values
		if closestNode != n.greater {
			closestNode.parent.lower = closestNode.greater
			if closestNode.greater != nil {
				closestNode.greater.parent = n.greater
			}
		} else {
			n.greater = closestNode.greater
			if n.greater != nil {
				n.greater.parent = n
			}
		}
	}
	closestNode.greater = nil
	closestNode.lower = nil
	closestNode.parent = nil
	r.nodePool.Put(closestNode)
}

func (r *Root[K, V]) deleteValue(node *node[K, V], value *V) error {
	for n, v := range node.values {
		found, err := r.comparer.CompareValues(*value, v)
		if err != nil {
			return err
		}
		if found {
			node.values = slices.Delete(node.values, n, n+1)
			return nil
		}
	}
	return nil
}

// GetAll implements bst.BST.
func (r *Root[K, V]) GetAll() iter.Seq[V] {
	return func(yield func(V) bool) {
		if !r.initialized {
			return
		}
		_ = r.getAll(&r.Node, yield)
	}
}

func (r *Root[K, V]) getAll(node *node[K, V], yield func(V) bool) bool {
	if node.lower != nil {
		if !r.getAll(node.lower, yield) {
			return false
		}
	}
	for _, value := range node.values {
		if !yield(value) {
			return false
		}
	}
	if node.greater == nil {
		return true
	}
	return r.getAll(node.greater, yield)
}

// GetMax implements bst.BST.
func (r *Root[K, V]) GetMax() bst.Node[K, V] {
	return r.getMax(&r.Node)
}

func (r *Root[K, V]) getMax(node *node[K, V]) *node[K, V] {
	for node.greater != nil {
		node = node.greater
	}
	return node
}

// GetMin implements bst.BST.
func (r *Root[K, V]) GetMin() bst.Node[K, V] {
	return r.getMin(&r.Node)
}

func (r *Root[K, V]) getMin(node *node[K, V]) *node[K, V] {
	for node.lower != nil {
		node = node.lower
	}
	return node
}

// GetNumberOfKeys implements bst.BST.
func (r *Root[K, V]) GetNumberOfKeys() int {
	return r.nodeCount
}

// Update implements bst.BST.
func (r *Root[K, V]) Update(key K, old V, nw V) error {
	node, err := r.search(key)
	if err != nil || node == nil {
		return err
	}
	for n, value := range node.values {
		equals, err := r.comparer.CompareValues(value, old)
		if err != nil {
			return err
		}
		if equals {
			node.values[n] = nw
			break
		}
	}
	return nil
}
