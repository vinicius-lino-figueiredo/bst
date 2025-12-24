// Package unbalanced TODO
package unbalanced

import (
	"iter"
	"math/rand"
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
		nodePool:     sync.Pool{New: func() any { return &bst.Node[K, V]{} }},
		Node: bst.Node[K, V]{
			Values: make([]V, 0, creationSize),
		},
	}
}

// Root TODO.
type Root[K any, V any] struct {
	bst.Node[K, V]
	initialized  bool
	nodeCount    int
	unique       bool
	creationSize int
	nodePool     sync.Pool
	comparer     bst.Comparer[K, V]
}

// Insert implements bst.BST.
func (r *Root[K, V]) Insert(key K, value V) error {
	if !r.initialized {
		r.Key = key
		r.initialized = true
		r.Values = append(r.Values, value)
		r.nodeCount++
		return nil
	}
	node := &r.Node
Loop:
	for {
		comparison, err := r.comparer.CompareKeys(key, node.Key)
		if err != nil {
			return err
		}
		switch {
		case comparison > 0:
			if node.Greater == nil {
				node.Greater = r.createEmptyNode(key, node)
				r.nodeCount++
				node = node.Greater
				break Loop
			}
			node = node.Greater
		case comparison < 0:
			if node.Lower == nil {
				node.Lower = r.createEmptyNode(key, node)
				r.nodeCount++
				node = node.Lower
				break Loop
			}
			node = node.Lower
		default:
			if r.unique {
				return bst.ErrUniqueViolated{Key: key}
			}
			break Loop
		}
	}
	node.Values = append(node.Values, value)
	return nil
}

func (r *Root[K, V]) createEmptyNode(key K, parent *bst.Node[K, V]) *bst.Node[K, V] {
	node := r.nodePool.Get().(*bst.Node[K, V])
	node.Key = key
	node.Values = make([]V, 0, r.creationSize)
	node.Parent = parent
	return node
}

// Search implements bst.BST.
func (r *Root[K, V]) Search(key K) (*bst.Node[K, V], error) {
	node := &r.Node
	for {
		comparison, err := r.comparer.CompareKeys(key, node.Key)
		if err != nil {
			return nil, err
		}
		switch {
		case comparison > 0 && node.Greater != nil:
			node = node.Greater
			continue
		case comparison < 0 && node.Lower != nil:
			node = node.Lower
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

func (r *Root[K, V]) doubleQuery(node *bst.Node[K, V], query bst.Query[K], yield func(V, error) bool) bool {
	ltComp, err := r.comparer.CompareKeys(node.Key, query.LowerThan.Value)
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

func (r *Root[K, V]) treatAboveMax(node *bst.Node[K, V], query bst.Query[K], yield func(V, error) bool) bool {
	if node.Lower == nil {
		return true
	}
	gtComp, err := r.comparer.CompareKeys(node.Key, query.GreaterThan.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}

	if gtComp < 0 {
		return true
	}

	return r.doubleQuery(node.Lower, query, yield)
}

func (r *Root[K, V]) treatBelowMax(node *bst.Node[K, V], query bst.Query[K], yield func(V, error) bool) bool {
	gtComp, err := r.comparer.CompareKeys(node.Key, query.GreaterThan.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}
	switch {
	case gtComp < 0: // node lower than min
		if r.Greater != nil {
			return r.doubleQuery(node.Greater, query, yield)
		}
	case gtComp == 0: // node equal to min
		if query.GreaterThan.IncludeEqual && !r.yieldValues(node, yield) {
			return false
		}
	default:
		if !r.yieldValues(node, yield) {
			return false
		}
	}
	if node.Greater != nil {
		return r.doubleQuery(node.Greater, query, yield)
	}
	return true
}

func (r *Root[K, V]) treatEqualMax(node *bst.Node[K, V], query bst.Query[K], yield func(V, error) bool) bool {
	gtComp, err := r.comparer.CompareKeys(node.Key, query.GreaterThan.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}
	switch {
	case gtComp > 0:
		if node.Lower != nil && !r.queryGreater(node.Lower, query.GreaterThan, yield) {
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

func (r *Root[K, V]) queryGreater(node *bst.Node[K, V], bound *bst.Bound[K], yield func(V, error) bool) bool {
	comp, err := r.comparer.CompareKeys(node.Key, bound.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}
	switch {
	case comp > 0:
		if node.Lower != nil && !r.queryGreater(node.Lower, bound, yield) {
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
	if node.Greater != nil {
		return r.queryGreater(node.Greater, bound, yield)
	}
	return true
}

func (r *Root[K, V]) queryLower(node *bst.Node[K, V], bound *bst.Bound[K], yield func(V, error) bool) bool {
	comp, err := r.comparer.CompareKeys(node.Key, bound.Value)
	if err != nil {
		yield(*new(V), err)
		return false
	}

	switch {
	case comp < 0:
		if node.Lower != nil && !r.queryLower(node.Lower, bound, yield) {
			return false
		}
		if !r.yieldValues(node, yield) {
			return false
		}
		if node.Greater != nil && !r.queryLower(node.Greater, bound, yield) {
			return false
		}
	case comp > 0:
		if node.Lower != nil {
			return r.queryLower(node.Lower, bound, yield)
		}
	default:
		if node.Lower != nil && !r.queryLower(node.Lower, bound, yield) {
			return false
		}
		if bound.IncludeEqual && !r.yieldValues(node, yield) {
			return false
		}
	}
	return true
}

func (r *Root[K, V]) yieldValues(node *bst.Node[K, V], yield func(V, error) bool) bool {
	for _, v := range node.Values {
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
	node, err := r.Search(key)
	if err != nil || node == nil {
		return err
	}
	if value != nil {
		if err = r.deleteValue(node, value); err != nil || len(node.Values) > 0 {
			return err
		}
	}
	r.nodeCount--

	switch {
	case node.Lower != nil:
		if node.Greater != nil {
			r.deleteDoubleChildrenNode(node)
			return nil
		}
		r.takePlace(node, node.Lower)
	case node.Greater != nil:
		r.takePlace(node, node.Greater)
	default:
		node.Values = node.Values[:0]
		switch node {
		case &r.Node:
			r.initialized = false
			return nil
		case node.Parent.Lower:
			node.Parent.Lower = nil
		default:
			node.Parent.Greater = nil
		}
		node.Parent = nil
		r.nodePool.Put(node)
	}

	return nil
}

func (r *Root[K, V]) takePlace(node, victim *bst.Node[K, V]) {
	node.Key, node.Values = victim.Key, victim.Values
	node.Greater, node.Lower = victim.Greater, victim.Lower
	if node.Greater != nil {
		node.Greater.Parent = node
	}
	if node.Lower != nil {
		node.Lower.Parent = node
	}

	victim.Parent, victim.Greater, victim.Lower = nil, nil, nil
	victim.Values = victim.Values[:0]
	r.nodePool.Put(victim)
}

func (r *Root[K, V]) deleteDoubleChildrenNode(node *bst.Node[K, V]) {
	if rand.Float32() > 0.5 {
		closestNode := r.getMax(node.Lower)

		// cloning closest value
		node.Key = closestNode.Key
		node.Values = closestNode.Values
		if closestNode != node.Lower {
			closestNode.Parent.Greater = closestNode.Lower
		} else {
			node.Lower = nil
		}
		closestNode.Greater = nil
		closestNode.Lower = nil
		closestNode.Parent = nil
		r.nodePool.Put(closestNode)
	} else {
		closestNode := r.getMin(node.Greater)

		// cloning closest value
		node.Key = closestNode.Key
		node.Values = closestNode.Values
		if closestNode != node.Greater {
			closestNode.Parent.Lower = closestNode.Greater
		} else {
			node.Greater = nil
		}
		r.nodePool.Put(closestNode)
	}
}

func (r *Root[K, V]) deleteValue(node *bst.Node[K, V], value *V) error {
	for n, v := range node.Values {
		found, err := r.comparer.CompareValues(*value, v)
		if err != nil {
			return err
		}
		if found {
			node.Values = slices.Delete(node.Values, n, n+1)
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

func (r *Root[K, V]) getAll(node *bst.Node[K, V], yield func(V) bool) bool {
	if node.Lower != nil {
		if !r.getAll(node.Lower, yield) {
			return false
		}
	}
	for _, value := range node.Values {
		if !yield(value) {
			return false
		}
	}
	if node.Greater == nil {
		return true
	}
	return r.getAll(node.Greater, yield)
}

// GetMax implements bst.BST.
func (r *Root[K, V]) GetMax() *bst.Node[K, V] {
	return r.getMax(&r.Node)
}

func (r *Root[K, V]) getMax(node *bst.Node[K, V]) *bst.Node[K, V] {
	for node.Greater != nil {
		node = node.Greater
	}
	return node
}

// GetMin implements bst.BST.
func (r *Root[K, V]) GetMin() *bst.Node[K, V] {
	return r.getMin(&r.Node)
}

func (r *Root[K, V]) getMin(node *bst.Node[K, V]) *bst.Node[K, V] {
	for node.Lower != nil {
		node = node.Lower
	}
	return node
}

// GetNumberOfKeys implements bst.BST.
func (r *Root[K, V]) GetNumberOfKeys() int {
	return r.nodeCount
}

// Update implements bst.BST.
func (r *Root[K, V]) Update(key K, old V, nw V) error {
	node, err := r.Search(key)
	if err != nil || node == nil {
		return err
	}
	for n, value := range node.Values {
		equals, err := r.comparer.CompareValues(value, old)
		if err != nil {
			return err
		}
		if equals {
			node.Values[n] = nw
			break
		}
	}
	return nil
}
