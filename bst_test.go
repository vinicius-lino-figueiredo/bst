package bst

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/suite"
)

// BinarySearchTreeTestSuite is the test suite for the bst package. All of the
// tests are directly based on the original code
type BinarySearchTreeTestSuite struct {
	suite.Suite
}

// Upon creation, left, right are null, key and data can be set
func (s *BinarySearchTreeTestSuite) TestCreation() {
	bst := NewBinarySearchTree(Options{})
	s.Nil(bst.left)
	s.Nil(bst.right)
	s.Nil(bst.Key())
	s.Len(bst.data, 0)

	bst = NewBinarySearchTree(Options{Key: 6, Value: "ggg"})
	s.Nil(bst.left)
	s.Nil(bst.right)
	s.Equal(6, bst.Key())
	s.Len(bst.data, 1)
	s.Equal("ggg", bst.data[0])
}

// Can get maxkey and minkey descendants
func (s *BinarySearchTreeTestSuite) TestSanityChecks() {
	t := NewBinarySearchTree(Options{Key: 10})
	l := NewBinarySearchTree(Options{Key: 5})
	r := NewBinarySearchTree(Options{Key: 15})
	ll := NewBinarySearchTree(Options{Key: 3})
	lr := NewBinarySearchTree(Options{Key: 8})
	rl := NewBinarySearchTree(Options{Key: 11})
	rr := NewBinarySearchTree(Options{Key: 42})

	t.left = l
	t.right = r
	l.left = ll
	l.right = lr
	r.left = rl
	r.right = rr

	// Getting min and max key descendants
	s.Equal(3, t.GetMinKeyDescendant().Key())
	s.Equal(42, t.GetMaxKeyDescendant().Key())

	s.Equal(3, t.left.GetMinKeyDescendant().Key())
	s.Equal(8, t.left.GetMaxKeyDescendant().Key())

	s.Equal(11, t.right.GetMinKeyDescendant().Key())
	s.Equal(42, t.right.GetMaxKeyDescendant().Key())

	s.Equal(11, t.right.left.GetMinKeyDescendant().Key())
	s.Equal(11, t.right.left.GetMaxKeyDescendant().Key())

	// Getting min and max keys
	s.Equal(3, t.GetMinKey())
	s.Equal(42, t.GetMaxKey())

	s.Equal(3, t.left.GetMinKey())
	s.Equal(8, t.left.GetMaxKey())

	s.Equal(11, t.right.GetMinKey())
	s.Equal(42, t.right.GetMaxKey())

	s.Equal(11, t.right.left.GetMinKey())
}

// Can check a condition against every node in a tree
func (s *BinarySearchTreeTestSuite) TestCheckCondition() {
	t := NewBinarySearchTree(Options{Key: 10})
	l := NewBinarySearchTree(Options{Key: 6})
	r := NewBinarySearchTree(Options{Key: 16})
	ll := NewBinarySearchTree(Options{Key: 4})
	lr := NewBinarySearchTree(Options{Key: 8})
	rl := NewBinarySearchTree(Options{Key: 12})
	rr := NewBinarySearchTree(Options{Key: 42})

	t.left = l
	t.right = r
	l.left = ll
	l.right = lr
	r.left = rl
	r.right = rr

	test := func(key any, _ any) error {
		if key.(int)%2 != 0 {
			return errors.New("Key is not even")
		}
		return nil
	}

	for _, node := range []*BinarySearchTree{l, r, ll, lr, rl, rr} {
		node.SetKey(node.Key().(int) + 1)
		err := t.CheckAllNodesFullfillCondition(test)
		s.Error(err)
		node.SetKey(node.Key().(int) - 1)
	}
}

// Can check that a tree verifies node ordering
func (s *BinarySearchTreeTestSuite) TestVerifyOrdening() {
	t := NewBinarySearchTree(Options{Key: float64(10)})
	l := NewBinarySearchTree(Options{Key: float64(5)})
	r := NewBinarySearchTree(Options{Key: float64(15)})
	ll := NewBinarySearchTree(Options{Key: float64(3)})
	lr := NewBinarySearchTree(Options{Key: float64(8)})
	rl := NewBinarySearchTree(Options{Key: float64(11)})
	rr := NewBinarySearchTree(Options{Key: float64(42)})

	t.left = l
	t.right = r
	l.left = ll
	l.right = lr
	r.left = rl
	r.right = rr

	err := t.CheckNodeOrdering()
	s.NoError(err)

	// Let's be paranoid and check all cases...
	l.SetKey(float64(12))
	s.Error(t.CheckNodeOrdering())

	l.SetKey(float64(5))

	r.SetKey(float64(9))
	s.Error(t.CheckNodeOrdering())
	r.SetKey(float64(15))

	ll.SetKey(float64(6))
	s.Error(t.CheckNodeOrdering())
	ll.SetKey(float64(11))
	s.Error(t.CheckNodeOrdering())
	ll.SetKey(float64(3))

	lr.SetKey(float64(4))
	s.Error(t.CheckNodeOrdering())
	lr.SetKey(float64(11))
	s.Error(t.CheckNodeOrdering())
	lr.SetKey(float64(8))

	rl.SetKey(float64(16))
	s.Error(t.CheckNodeOrdering())
	rl.SetKey(float64(9))
	s.Error(t.CheckNodeOrdering())
	rl.SetKey(float64(11))

	rr.SetKey(float64(12))
	s.Error(t.CheckNodeOrdering())
	rr.SetKey(float64(7))
	s.Error(t.CheckNodeOrdering())
	rr.SetKey(float64(10.5))
	s.Error(t.CheckNodeOrdering())
	rr.SetKey(float64(42))

	s.NoError(t.CheckNodeOrdering())
}

// Checking if a tree's internal pointers (i.e. parents) are correct
func (s *BinarySearchTreeTestSuite) TestInternalPointers() {
	t := NewBinarySearchTree(Options{Key: 10})
	l := NewBinarySearchTree(Options{Key: 5})
	r := NewBinarySearchTree(Options{Key: 15})
	ll := NewBinarySearchTree(Options{Key: 3})
	lr := NewBinarySearchTree(Options{Key: 8})
	rl := NewBinarySearchTree(Options{Key: 11})
	rr := NewBinarySearchTree(Options{Key: 42})

	t.left = l
	t.right = r
	l.left = ll
	l.right = lr
	r.left = rl
	r.right = rr

	s.Error(t.CheckInternalPointers())
	l.parent = t
	s.Error(t.CheckInternalPointers())
	r.parent = t
	s.Error(t.CheckInternalPointers())
	ll.parent = l
	s.Error(t.CheckInternalPointers())
	lr.parent = l
	s.Error(t.CheckInternalPointers())
	rl.parent = r
	s.Error(t.CheckInternalPointers())
	rr.parent = r

	s.NoError(t.CheckInternalPointers())
}

// it('Can get the number of inserted keys', () => {
func (s *BinarySearchTreeTestSuite) TestGetNumberOfKeys() {
	bst := NewBinarySearchTree(Options{})

	s.Equal(0, bst.GetNumberOfKeys())
	bst.Insert(10, nil)
	s.Equal(1, bst.GetNumberOfKeys())
	bst.Insert(5, nil)
	s.Equal(2, bst.GetNumberOfKeys())
	bst.Insert(3, nil)
	s.Equal(3, bst.GetNumberOfKeys())
	bst.Insert(8, nil)
	s.Equal(4, bst.GetNumberOfKeys())
	bst.Insert(15, nil)
	s.Equal(5, bst.GetNumberOfKeys())
	bst.Insert(12, nil)
	s.Equal(6, bst.GetNumberOfKeys())
	bst.Insert(37, nil)
	s.Equal(7, bst.GetNumberOfKeys())
}

func (s *BinarySearchTreeTestSuite) TestInsert() {

	// Insert at the root if its the first insertion
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{})

		bst.Insert(10, "some data")

		s.NoError(bst.CheckIsBST())
		s.Equal(10, bst.Key())
		s.Equal([]any{"some data"}, bst.data)
		s.Nil(bst.left)
		s.Nil(bst.right)
	})

	// Insert on the left if key is less than the root's
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{})

		bst.Insert(10, "some data")
		bst.Insert(7, "some other data")

		s.NoError(bst.CheckIsBST())
		s.Nil(bst.right)
		s.Equal(7, bst.left.Key())
		s.Equal([]any{"some other data"}, bst.left.data)
		s.Nil(bst.left.left)
		s.Nil(bst.left.right)
	})

	// Insert on the right if key is greater than the root's
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{})

		bst.Insert(10, "some data")
		bst.Insert(14, "some other data")

		s.NoError(bst.CheckIsBST())
		s.Nil(bst.left)
		s.Equal(14, bst.right.Key())
		s.Equal([]any{"some other data"}, bst.right.data)
		s.Nil(bst.right.left)
		s.Nil(bst.right.right)
	})

	// Recursive insertion on the left works
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{})

		bst.Insert(10, "some data")
		bst.Insert(7, "some other data")
		bst.Insert(1, "hello")
		bst.Insert(9, "world")

		s.NoError(bst.CheckIsBST())
		s.Nil(bst.right)
		s.Equal(7, bst.left.Key())
		s.Equal([]any{"some other data"}, bst.left.data)

		s.Equal(1, bst.left.left.Key())
		s.Equal([]any{"hello"}, bst.left.left.data)

		s.Equal(9, bst.left.right.Key())
		s.Equal([]any{"world"}, bst.left.right.data)

	})

	// Recursive insertion on the right works
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{})

		bst.Insert(10, "some data")
		bst.Insert(17, "some other data")
		bst.Insert(11, "hello")
		bst.Insert(19, "world")

		s.NoError(bst.CheckIsBST())
		s.Nil(bst.left)
		s.Equal(17, bst.right.Key())
		s.Equal([]any{"some other data"}, bst.right.data)

		s.Equal(11, bst.right.left.Key())
		s.Equal([]any{"hello"}, bst.right.left.data)

		s.Equal(19, bst.right.right.Key())
		s.Equal([]any{"world"}, bst.right.right.data)
	})

	// If uniqueness constraint not enforced, we can insert different data for same key
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{})

		bst.Insert(10, "some data")
		bst.Insert(3, "hello")
		bst.Insert(3, "world")

		s.NoError(bst.CheckIsBST())
		s.Equal(3, bst.left.Key())
		s.Equal([]any{"hello", "world"}, bst.left.data)

		bst.Insert(12, "a")
		bst.Insert(12, "b")

		s.NoError(bst.CheckIsBST())
		s.Equal(12, bst.right.Key())
		s.Equal([]any{"a", "b"}, bst.right.data)
	})

	// If uniqueness constraint is enforced, we cannot insert different data for same key
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{Unique: true})

		bst.Insert(10, "some data")
		bst.Insert(3, "hello")

		err := bst.Insert(3, "world")
		if err != nil {
			var e *ErrViolated
			s.ErrorAs(err, &e)
			s.Equal(3, e.key)
		}

		s.NoError(bst.CheckIsBST())

		s.Equal(3, bst.left.Key())
		s.Equal([]any{"hello"}, bst.left.data)

		bst.Insert(12, "a")

		err = bst.Insert(12, "world")
		if err != nil {
			var e *ErrViolated
			s.ErrorAs(err, &e)
			s.Equal(12, e.key)
		}

		s.NoError(bst.CheckIsBST())
		s.Equal(12, bst.right.Key())
		s.Equal([]any{"a"}, bst.right.data)
	})

	// Can insert 0 or the empty string
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{})

		bst.Insert(0, "some data")

		s.NoError(bst.CheckIsBST())

		s.Equal(0, bst.Key())
		s.Equal([]any{"some data"}, bst.data)
		s.Nil(bst.left)
		s.Nil(bst.right)

		bst = NewBinarySearchTree(Options{})

		bst.Insert("", "some other data")

		s.NoError(bst.CheckIsBST())
		s.Equal("", bst.Key())
		s.Equal([]any{"some other data"}, bst.data)
		s.Nil(bst.left)
		s.Nil(bst.right)
	})
	// Can insert a lot of keys and still get a BST (sanity check)
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{})

		randarr := getRandomArray(100)
		for _, n := range randarr {
			bst.Insert(n, "some data")
		}

		s.NoError(bst.CheckIsBST())
	})

	// All children get a pointer to their parent, the root doesnt
	s.Run("TestInsert_InsertRootIfFirst", func() {
		bst := NewBinarySearchTree(Options{})

		bst.Insert(10, "root")
		bst.Insert(5, "yes")
		bst.Insert(15, "no")

		s.NoError(bst.CheckIsBST())

		s.Nil(bst.parent)
		s.Equal(bst, bst.left.parent)
		s.Equal(bst, bst.right.parent)
	})
} // ==== End of 'Insertion' ==== //

func getRandomArray(n int) []any {
	res := make([]any, n)

	for i := range n {
		res[i] = i
	}

	rand.Shuffle(n, func(i, j int) {
		x, y := res[i], res[j]
		res[i], res[j] = y, x
	})

	return res
}

func TestBinarySearchTreeTestSuite(t *testing.T) {
	suite.Run(t, new(BinarySearchTreeTestSuite))
}
