// Package unbalanced TODO
package unbalanced

import (
	"iter"
	"math/rand"
	"slices"
	"sync"

	"github.com/vinicius-lino-figueiredo/bst/domain"
)

// RecursiveLimit TODO
var RecursiveLimit = 100000

// NewBST TODO
func NewBST[K any, V any](unique bool, creationSize int, comparer domain.Comparer[K, V]) domain.BST[K, V] {
	if unique {
		creationSize = 1
	} else if creationSize <= 0 {
		creationSize = 8
	}
	return newBST(unique, creationSize, comparer)
}

func newBST[K any, V any](unique bool, creationSize int, comparer domain.Comparer[K, V]) *Root[K, V] {
	return &Root[K, V]{
		unique:       unique,
		creationSize: creationSize,
		comparer:     comparer,
		nodePool:     sync.Pool{New: func() any { return &domain.Node[K, V]{} }},
		Node: domain.Node[K, V]{
			Values: make([]V, 0, creationSize),
		},
	}
}

// Root TODO.
type Root[K any, V any] struct {
	domain.Node[K, V]
	initialized  bool
	nodeCount    int
	unique       bool
	creationSize int
	nodePool     sync.Pool
	comparer     domain.Comparer[K, V]
}

// Delete implements domain.BST.
func (r *Root[K, V]) Delete(key K, value *V) error {
	if !r.initialized {
		return nil
	}
	node, err := r.search(key)
	if err != nil || node == nil {
		return err
	}
	if value != nil {
		if err = r.deleteValue(node, value); err != nil || len(node.Values) > 0 {
			return err
		}
	}
	node.Values = node.Values[:0]

	switch {
	case node.Lower != nil:
		if node.Greater != nil {
			r.deleteDoubleChildrenNode(node)
			return nil
		}
		node.Key = node.Lower.Key
		node.Values = node.Lower.Values
		node.Lower = node.Lower.Lower
		node.Lower, node = node.Lower.Lower, node.Lower
	case node.Greater != nil:
		node.Key = node.Greater.Key
		node.Values = node.Greater.Values
		node.Lower = node.Greater.Lower
		node.Greater, node = node.Greater.Greater, node.Greater
	default:
		r.initialized = false
		return nil
	}

	node.Parent = nil
	node.Lower = nil
	node.Lower = nil
	r.nodePool.Put(node)

	return nil
}

func (r *Root[K, V]) deleteDoubleChildrenNode(node *domain.Node[K, V]) {
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
	r.nodeCount--
}

func (r *Root[K, V]) deleteValue(node *domain.Node[K, V], value *V) error {
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

// GetMax implements domain.BST.
func (r *Root[K, V]) GetMax() *domain.Node[K, V] {
	return r.getMax(&r.Node)
}

func (r *Root[K, V]) getMax(node *domain.Node[K, V]) *domain.Node[K, V] {
	for node.Greater != nil {
		node = node.Greater
	}
	return node
}

// GetMin implements domain.BST.
func (r *Root[K, V]) GetMin() *domain.Node[K, V] {
	return r.getMin(&r.Node)
}

func (r *Root[K, V]) getMin(node *domain.Node[K, V]) *domain.Node[K, V] {
	for node.Lower != nil {
		node = node.Lower
	}
	return node
}

// GetNumberOfKeys implements domain.BST.
func (r *Root[K, V]) GetNumberOfKeys() int {
	if !r.initialized {
		return 0
	}
	return r.nodeNumberOfKeys(&r.Node)
}

func (r *Root[K, V]) nodeNumberOfKeys(node *domain.Node[K, V]) int {
	res := 0
	if node.Lower != nil {
		res += r.nodeNumberOfKeys(node.Lower)
	}
	if node.Greater != nil {
		res += r.nodeNumberOfKeys(node.Greater)
	}
	return res
}

// Insert implements domain.BST.
func (r *Root[K, V]) Insert(key K, value V) error {
	if !r.initialized {
		r.Key = key
		r.initialized = true
		r.Values = append(r.Values, value)
		r.nodeCount++
		return nil
	}
	node := &r.Node
	for {
		comparison, err := r.comparer.CompareKeys(key, node.Key)
		if err != nil {
			return err
		}
		switch {
		case comparison > 0:
			if node.Greater == nil {
				node.Greater = r.createEmptyNode(key, node)
			}
			node = node.Greater
		case comparison < 0:
			if node.Lower == nil {
				node.Lower = r.createEmptyNode(key, node)
			}
			node = node.Lower
		default:
			if len(node.Values) != 0 {
				if r.unique {
					return domain.ErrUniqueViolated{Key: key}
				}
			} else {
				r.nodeCount++
			}
			node.Values = append(node.Values, value)
			return nil
		}
	}
}

func (r *Root[K, V]) createEmptyNode(key K, parent *domain.Node[K, V]) *domain.Node[K, V] {
	node := r.nodePool.Get().(*domain.Node[K, V])
	node.Key = key
	node.Values = make([]V, 0, r.creationSize)
	node.Parent = parent
	return node
}

func (r *Root[K, V]) doubleQueryRecursive(node *domain.Node[K, V], query domain.Query[K], res *[]V) error {
	lowerComp, err := r.comparer.CompareKeys(node.Key, query.LowerThan.Value)
	if err != nil {
		return err
	}
	greaterComp, err := r.comparer.CompareKeys(node.Key, query.GreaterThan.Value)
	if err != nil {
		return err
	}

	switch {
	case lowerComp > 0: // Node is greater than max
		if node.Lower == nil {
			break
		}

		// Node is within min bound but is not equal to it
		if greaterComp > 0 {
			return r.doubleQueryRecursive(node.Lower, query, res)
		}
	case lowerComp < 0: // Node is lower than max
		switch {

		// Node is equal to min and min bound is inclusive
		case greaterComp == 0:
			if query.GreaterThan.IncludeEqual {
				*res = append(*res, node.Values...)
				if node.Greater != nil {
					return r.singleQueryRecursive(node.Greater, query.LowerThan, 1, res)
				}
			}
			if node.Greater != nil {
				return r.singleQueryRecursive(node.Greater, query.LowerThan, 1, res)
			}

		// Node is lower than min but has a greater child node
		case greaterComp < 0 && node.Greater != nil:
			return r.doubleQueryRecursive(node.Greater, query, res)

		case greaterComp > 0: // Node is greater than min
			if node.Lower != nil {
				if err := r.doubleQueryRecursive(node.Lower, query, res); err != nil {
					return err
				}
			}
			*res = append(*res, node.Values...)

			if node.Greater != nil {
				return r.doubleQueryRecursive(node.Greater, query, res)
			}

		}
	default: // Node is equal to max
		if !query.LowerThan.IncludeEqual {
			if node.Lower != nil && (greaterComp > 0 || (greaterComp == 0 && query.GreaterThan.IncludeEqual)) {
				return r.singleQueryRecursive(node.Lower, query.GreaterThan, -1, res)
			}
			return nil
		}
		switch {
		case greaterComp == 0 && query.GreaterThan.IncludeEqual: // Node is equal to min
			*res = append(*res, node.Values...)
		case greaterComp > 0: // Node is greater than min
			return r.singleQueryRecursive(node.Lower, query.GreaterThan, -1, res)
		case greaterComp < 0: // Node is lower than min
		}
	}
	return nil
}

func (r *Root[K, V]) singleQueryRecursive(node *domain.Node[K, V], bound *domain.Bound[K], multiplier int, res *[]V) error {
	comp, err := r.comparer.CompareKeys(bound.Value, node.Key)
	if err != nil {
		return err
	}
	comp *= multiplier
	switch {
	case comp > 0: // node is within bound
		if multiplier < 0 { // bound is greaterThan
			if node.Lower != nil {
				if err := r.singleQueryRecursive(node.Lower, bound, multiplier, res); err != nil {
					return nil
				}
			}
			*res = append(*res, node.Values...)
			if node.Greater != nil {
				return r.singleQueryRecursive(node.Greater, bound, multiplier, res)
			}
			break
		}
		if node.Lower != nil {
			r.addAnyChildAndSelf(node.Lower, res)
		}
		*res = append(*res, node.Values...)
		if node.Greater != nil {
			if err := r.singleQueryRecursive(node.Greater, bound, multiplier, res); err != nil {
				return nil
			}
		}
	case comp < 0:
		next := node.Lower
		if multiplier < 0 {
			next = node.Greater
		}
		if next != nil {
			return r.singleQueryRecursive(next, bound, multiplier, res)
		}
	case comp == 0:
		switch {
		case multiplier < 0:
			if bound.IncludeEqual {
				*res = append(*res, node.Values...)
			}
			if node.Greater != nil {
				r.addAnyChildAndSelf(node.Greater, res)
			}
		default:
			if node.Lower != nil {
				r.addAnyChildAndSelf(node.Lower, res)
			}
			if bound.IncludeEqual {
				*res = append(*res, node.Values...)
			}
		}
	}
	return nil
}

func (r *Root[K, V]) addAnyChildAndSelf(node *domain.Node[K, V], res *[]V) {
	if node.Lower != nil {
		r.addAnyChildAndSelf(node.Lower, res)
	}
	*res = append(*res, node.Values...)
	if node.Greater != nil {
		r.addAnyChildAndSelf(node.Greater, res)
	}
}

// Query implements domain.BST.
func (r *Root[K, V]) Query(query domain.Query[K]) ([]V, error) {
	res := make([]V, 0, r.nodeCount*r.creationSize)

	if err := r.runQuery(query, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Root[K, V]) runQuery(query domain.Query[K], res *[]V) error {
	if query.LowerThan != nil {
		if query.GreaterThan != nil {
			if r.nodeCount < RecursiveLimit {
				return r.doubleQueryRecursive(&r.Node, query, res)
			}
			return r.doubleQuery(query, res)
		}
		if r.nodeCount < RecursiveLimit {
			return r.singleQueryRecursive(&r.Node, query.LowerThan, 1, res)
		}
		return r.singleQuery(query.LowerThan, 1, res)
	}
	if query.GreaterThan != nil {
		if r.nodeCount < RecursiveLimit {
			return r.singleQueryRecursive(&r.Node, query.GreaterThan, -1, res)
		}
		return r.singleQuery(query.LowerThan, -1, res)
	}
	return nil
}

func (r *Root[K, V]) doubleQuery(query domain.Query[K], res *[]V) error {
	nodes := make([]*domain.Node[K, V], 1, r.nodeCount)
	nodes[0] = &r.Node
	lastLen := 1
	start, end := 0, 1
	for {
		for n, node := range nodes[start:end] {
			comparison, err := r.comparer.CompareKeys(query.LowerThan.Value, node.Key)
			if err != nil {
				return err
			}
			switch {
			case comparison < 0:
				if node.Lower != nil {
					nodes = append(nodes, node.Lower)
				}
			case comparison > 0:
				err := r.addChildrenAndSelf(node, nodes[start+n:cap(nodes)], res, 1, query.GreaterThan)
				if err != nil {
					return err
				}
			case comparison == 0:
				if query.LowerThan.IncludeEqual {
					*res = append(*res, node.Values...)
				}
				if node.Lower != nil {
					nodes = append(nodes, node.Lower)
				}
			}
		}
		diff := len(nodes) - lastLen
		if diff == 0 {
			return nil
		}
		lastLen = len(nodes)
		start = end
		end += diff
	}
}

func (r *Root[K, V]) addChildrenAndSelf(node *domain.Node[K, V], nodes []*domain.Node[K, V], res *[]V, multiplier int, bound *domain.Bound[K]) error {
	nodes = nodes[:1]
	nodes[0] = node
	lastLen := 1
	start, end := 0, 1
	for {
		for _, node := range nodes[start:end] {
			comp, err := r.comparer.CompareKeys(bound.Value, node.Key)
			if err != nil {
				return err
			}
			comp *= multiplier
			switch {
			case comp < 0:
				*res = append(*res, node.Values...)
				if node.Lower != nil {
					nodes = append(nodes, node.Lower)
				}
				if node.Greater != nil {
					nodes = append(nodes, node.Greater)
				}
			case comp == 0 && bound.IncludeEqual:
				println("")
			default:
				continue
			}
		}
		diff := len(nodes) - lastLen
		if diff == 0 {
			return nil
		}
		lastLen = len(nodes)
		start = end
		end += diff
	}

}

func (r *Root[K, V]) singleQuery(query *domain.Bound[K], multiplier int, res *[]V) error {
	nodes := make([]*domain.Node[K, V], 1, r.nodeCount)
	nodes[0] = &r.Node
	lastLen := 1
	start, end := 0, 1
	for {
		for n, node := range nodes[start:end] {
			comparison, err := r.comparer.CompareKeys(query.Value, node.Key)
			if err != nil {
				return err
			}
			comparison *= multiplier
			switch {
			case comparison < 0:
				if node.Lower != nil {
					nodes = append(nodes, node.Lower)
				}
			case comparison > 0:
				r.addChildrenAndSelfNoBound(node, nodes[start+n:cap(nodes)], res)
			case comparison == 0:
				if query.IncludeEqual {
					*res = append(*res, node.Values...)
				}
				if node.Lower != nil {
					nodes = append(nodes, node.Lower)
				}
			}
		}
		diff := len(nodes) - lastLen
		if diff == 0 {
			return nil
		}
		lastLen = len(nodes)
		start = end
		end += diff
	}
}

func (r *Root[K, V]) addChildrenAndSelfNoBound(node *domain.Node[K, V], nodes []*domain.Node[K, V], res *[]V) {
	nodes = nodes[:1]
	nodes[0] = node
	lastLen := 1
	start, end := 0, 1
	for {
		for _, node := range nodes[start:end] {
			*res = append(*res, node.Values...)
			if node.Lower != nil {
				nodes = append(nodes, node.Lower)
			}
			if node.Greater != nil {
				nodes = append(nodes, node.Greater)
			}
		}
		diff := len(nodes) - lastLen
		if diff == 0 {
			return
		}
		lastLen = len(nodes)
		start = end
		end += diff
	}

}

// Search implements domain.BST.
func (r *Root[K, V]) Search(key K) (*domain.Node[K, V], error) {
	node, err := r.search(key)
	return node, err
}

func (r *Root[K, V]) search(key K) (*domain.Node[K, V], error) {
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

// Update implements domain.BST.
func (r *Root[K, V]) Update(key K, old V, nw V) error {
	node, err := r.search(key)
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

// GetAll implements domain.BST.
func (r *Root[K, V]) GetAll() iter.Seq[V] {
	return func(yield func(V) bool) {
		if r.initialized {
			r.yieldNode(&r.Node, yield)
		}
	}
}

func (r *Root[K, V]) yieldNode(node *domain.Node[K, V], yield func(V) bool) {
	if node.Lower != nil {
		r.yieldNode(node.Lower, yield)
	}
	for _, value := range node.Values {
		if !yield(value) {
			return
		}
	}
	if node.Greater != nil {
		r.yieldNode(node.Greater, yield)
	}
}
