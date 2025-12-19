package unbalanced_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vinicius-lino-figueiredo/bst"
	"github.com/vinicius-lino-figueiredo/bst/adapter/comparer"
	"github.com/vinicius-lino-figueiredo/bst/adapter/unbalanced"
)

type BSTTestSuite struct {
	suite.Suite
	b *unbalanced.Root[string, int]
}

func (s *BSTTestSuite) SetupTest() {
	comparer := comparer.NewComparer[string, int]()
	s.b = unbalanced.NewBST(false, 0, comparer).(*unbalanced.Root[string, int])

	s.NoError(s.b.Insert("Leo", 76))
	s.NoError(s.b.Insert("Alice", 42))
	s.NoError(s.b.Insert("Marcus", 87))
	s.NoError(s.b.Insert("Luna", 15))
	s.NoError(s.b.Insert("Felix", 63))
	s.NoError(s.b.Insert("Nina", 91))
	s.NoError(s.b.Insert("Oscar", 28))
	s.NoError(s.b.Insert("Maya", 54))
	s.NoError(s.b.Insert("Alice", 23))
	s.NoError(s.b.Insert("Iris", 33))
	s.NoError(s.b.Insert("Hugo", 88))
	s.NoError(s.b.Insert("Zara", 19))
	s.NoError(s.b.Insert("Felix", 55))
	s.NoError(s.b.Insert("Kai", 45))
	s.NoError(s.b.Insert("Nora", 67))
	s.NoError(s.b.Insert("Theo", 92))
	s.NoError(s.b.Insert("Luna", 38))
	s.NoError(s.b.Insert("Mila", 11))
	s.NoError(s.b.Insert("Oscar", 72))
	s.NoError(s.b.Insert("Maya", 49))
}

func (s *BSTTestSuite) TestInsert() {

	comparer := comparer.NewComparer[string, int]()

	b := unbalanced.NewBST(false, 0, comparer).(*unbalanced.Root[string, int])

	s.NoError(b.Insert("Jack", 10))
	Jack := &b.Node
	s.NotNil(Jack)
	s.Equal("Jack", Jack.Key)
	s.Equal([]int{10}, Jack.Values)

	s.NoError(b.Insert("Diango", 9090))
	Diango := Jack.Lower
	s.NotNil(Diango)
	s.Equal("Diango", Diango.Key)
	s.Equal([]int{9090}, Diango.Values)

	s.NoError(b.Insert("Isaac", 2345678))
	Isaac := Diango.Greater
	s.NotNil(Isaac)
	s.Equal("Isaac", Isaac.Key)
	s.Equal([]int{2345678}, Isaac.Values)

	s.NoError(b.Insert("Kojiro", 12))
	Kojiro := Jack.Greater
	s.NotNil(Kojiro)
	s.Equal("Kojiro", Kojiro.Key)
	s.Equal([]int{12}, Kojiro.Values)

	s.NoError(b.Insert("Sam", 69))
	Sam := Kojiro.Greater
	s.NotNil(Sam)
	s.Equal("Sam", Sam.Key)
	s.Equal([]int{69}, Sam.Values)

	s.NoError(b.Insert("Jonathan", 20))
	Jonathan := Kojiro.Lower
	s.NotNil(Jonathan)
	s.Equal("Jonathan", Jonathan.Key)
	s.Equal([]int{20}, Jonathan.Values)

	s.NoError(b.Insert("Stan", 100))
	Stan := Sam.Greater
	s.NotNil(Stan)
	s.Equal("Stan", Stan.Key)
	s.Equal([]int{100}, Stan.Values)
}

func (s *BSTTestSuite) TestInsertNonUnique() {
	b := unbalanced.NewBST(true, 8, comparer.NewComparer[string, int]())

	s.NoError(b.Insert("unique", 10))
	s.ErrorAs(b.Insert("unique", 11), &bst.ErrUniqueViolated{})

}

func (s *BSTTestSuite) TestSearch() {
	comparer := comparer.NewComparer[string, int]()
	b := unbalanced.NewBST(false, 0, comparer).(*unbalanced.Root[string, int])

	s.NoError(b.Insert("Alice", 42))
	s.NoError(b.Insert("Marcus", 87))
	s.NoError(b.Insert("Luna", 15))
	s.NoError(b.Insert("Felix", 63))
	s.NoError(b.Insert("Nina", 91))
	s.NoError(b.Insert("Oscar", 28))
	s.NoError(b.Insert("Maya", 54))
	s.NoError(b.Insert("Leo", 76))
	s.NoError(b.Insert("Alice", 23))
	s.NoError(b.Insert("Iris", 33))
	s.NoError(b.Insert("Hugo", 88))
	s.NoError(b.Insert("Zara", 19))
	s.NoError(b.Insert("Felix", 55))
	s.NoError(b.Insert("Kai", 45))
	s.NoError(b.Insert("Nora", 67))
	s.NoError(b.Insert("Theo", 92))
	s.NoError(b.Insert("Luna", 38))
	s.NoError(b.Insert("Mila", 11))
	s.NoError(b.Insert("Oscar", 72))
	s.NoError(b.Insert("Maya", 49))

	Alice, err := b.Search("Alice")
	s.NoError(err)
	s.Equal("Alice", Alice.Key)
	s.Equal([]int{42, 23}, Alice.Values)

	Felix, err := b.Search("Felix")
	s.NoError(err)
	s.Equal("Felix", Felix.Key)
	s.Equal([]int{63, 55}, Felix.Values)

	Hugo, err := b.Search("Hugo")
	s.NoError(err)
	s.Equal("Hugo", Hugo.Key)
	s.Equal([]int{88}, Hugo.Values)

	Iris, err := b.Search("Iris")
	s.NoError(err)
	s.Equal("Iris", Iris.Key)
	s.Equal([]int{33}, Iris.Values)

	Kai, err := b.Search("Kai")
	s.NoError(err)
	s.Equal("Kai", Kai.Key)
	s.Equal([]int{45}, Kai.Values)

	Leo, err := b.Search("Leo")
	s.NoError(err)
	s.Equal("Leo", Leo.Key)
	s.Equal([]int{76}, Leo.Values)

	Luna, err := b.Search("Luna")
	s.NoError(err)
	s.Equal("Luna", Luna.Key)
	s.Equal([]int{15, 38}, Luna.Values)

	Marcus, err := b.Search("Marcus")
	s.NoError(err)
	s.Equal("Marcus", Marcus.Key)
	s.Equal([]int{87}, Marcus.Values)

	Maya, err := b.Search("Maya")
	s.NoError(err)
	s.Equal("Maya", Maya.Key)
	s.Equal([]int{54, 49}, Maya.Values)

	Mila, err := b.Search("Mila")
	s.NoError(err)
	s.Equal("Mila", Mila.Key)
	s.Equal([]int{11}, Mila.Values)

	Nina, err := b.Search("Nina")
	s.NoError(err)
	s.Equal("Nina", Nina.Key)
	s.Equal([]int{91}, Nina.Values)

	Nora, err := b.Search("Nora")
	s.NoError(err)
	s.Equal("Nora", Nora.Key)
	s.Equal([]int{67}, Nora.Values)

	Oscar, err := b.Search("Oscar")
	s.NoError(err)
	s.Equal("Oscar", Oscar.Key)
	s.Equal([]int{28, 72}, Oscar.Values)

	Theo, err := b.Search("Theo")
	s.NoError(err)
	s.Equal("Theo", Theo.Key)
	s.Equal([]int{92}, Theo.Values)

	Zara, err := b.Search("Zara")
	s.NoError(err)
	s.Equal("Zara", Zara.Key)
	s.Equal([]int{19}, Zara.Values)

	Invalid, err := b.Search("Invalid")
	s.NoError(err)
	s.Nil(Invalid)
}

func (s *BSTTestSuite) TestDelete() {
	comparer := comparer.NewComparer[string, int]()
	b := unbalanced.NewBST(false, 0, comparer).(*unbalanced.Root[string, int])

	s.NoError(b.Insert("Alice", 42))
	s.NoError(b.Insert("Marcus", 87))
	s.NoError(b.Insert("Luna", 15))
	s.NoError(b.Insert("Felix", 63))
	s.NoError(b.Insert("Nina", 91))
	s.NoError(b.Insert("Oscar", 28))
	s.NoError(b.Insert("Maya", 54))
	s.NoError(b.Insert("Leo", 76))
	s.NoError(b.Insert("Alice", 23))
	s.NoError(b.Insert("Iris", 33))
	s.NoError(b.Insert("Hugo", 88))
	s.NoError(b.Insert("Zara", 19))
	s.NoError(b.Insert("Leo", 77))
	s.NoError(b.Insert("Felix", 55))
	s.NoError(b.Insert("Kai", 45))
	s.NoError(b.Insert("Nora", 67))
	s.NoError(b.Insert("Theo", 92))
	s.NoError(b.Insert("Luna", 38))
	s.NoError(b.Insert("Mila", 11))
	s.NoError(b.Insert("Oscar", 72))
	s.NoError(b.Insert("Maya", 49))

	node, err := b.Search("Felix")
	s.NoError(err)
	s.Equal([]int{63, 55}, node.Values)
	Luna := node.Parent
	s.Equal("Luna", Luna.Key)
	s.Equal("Felix", Luna.Lower.Key)

	// Remove a single item
	value := 63
	s.NoError(b.Delete("Felix", &value))

	// Key still exists, but the delete value is gone
	node, err = b.Search("Felix")
	s.NoError(err)
	s.Equal([]int{55}, node.Values)
	s.Equal("Luna", Luna.Key)
	s.Equal("Felix", Luna.Lower.Key)

	// Deleting other value. Empty Nodes are removed too
	value = 55
	s.NoError(b.Delete("Felix", &value))
	s.Equal("Leo", Luna.Lower.Key)
	node, err = b.Search("Felix")
	s.NoError(err)
	s.Nil(node)

	// Another node took its place
	node, err = b.Search("Leo")
	s.NoError(err)
	s.Equal([]int{76, 77}, node.Values)
	s.Equal("Luna", Luna.Key)
	s.Equal("Leo", Luna.Lower.Key)

	// Deleting it without a value removes the whole node
	s.NoError(b.Delete("Leo", nil))
	node, err = b.Search("Leo")
	s.NoError(err)
	s.Nil(node)
	s.Equal("Iris", Luna.Lower.Key)
}

func (s *BSTTestSuite) TestDeleteRoot() {
	b := unbalanced.NewBST(false, 0, comparer.NewComparer[string, int]())

	s.NoError(b.Insert("b", 2))
	s.NoError(b.Insert("a", 1))
	s.NoError(b.Insert("c", 3))

}

func (s *BSTTestSuite) TestSimpleQuery() {
	data, err := s.b.Query(bst.Query[string]{
		GreaterThan: &bst.Bound[string]{
			Value: "Zara", IncludeEqual: false,
		}},
	)
	s.NoError(err)
	s.Equal([]int{}, data)

	data, err = s.b.Query(bst.Query[string]{
		GreaterThan: &bst.Bound[string]{
			Value: "Zara", IncludeEqual: true,
		}},
	)
	s.NoError(err)
	s.Equal([]int{19}, data)

	data, err = s.b.Query(bst.Query[string]{
		GreaterThan: &bst.Bound[string]{
			Value: "Theo", IncludeEqual: true,
		}},
	)
	s.NoError(err)
	s.Equal([]int{92, 19}, data)

	data, err = s.b.Query(bst.Query[string]{
		GreaterThan: &bst.Bound[string]{
			Value: "Theo", IncludeEqual: false,
		}},
	)
	s.NoError(err)
	s.Equal([]int{19}, data)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan: &bst.Bound[string]{
			Value: "Alice", IncludeEqual: false,
		}},
	)
	s.NoError(err)
	s.Equal([]int{}, data)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan: &bst.Bound[string]{
			Value: "Alice", IncludeEqual: true,
		}},
	)
	s.NoError(err)
	s.Equal([]int{42, 23}, data)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan: &bst.Bound[string]{
			Value: "Felix", IncludeEqual: true,
		}},
	)
	s.NoError(err)
	s.Equal([]int{42, 23, 63, 55}, data)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan: &bst.Bound[string]{
			Value: "Felix", IncludeEqual: false,
		}},
	)
	s.NoError(err)
	s.Equal([]int{42, 23}, data)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan: &bst.Bound[string]{
			Value: "Zhonya", IncludeEqual: false,
		}},
	)
	s.NoError(err)
	s.Equal([]int{42, 23, 63, 55, 88, 33, 45, 76, 15, 38, 87, 54, 49, 11, 91, 67, 28, 72, 92, 19}, data)

}

func (s *BSTTestSuite) TestMutualExcludingQueryBounds() {
	// bounds that cancel each other out
	data, err := s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: false},
		GreaterThan: &bst.Bound[string]{Value: "Zara", IncludeEqual: false}},
	)
	s.NoError(err)
	s.Len(data, 0)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: true},
		GreaterThan: &bst.Bound[string]{Value: "Zara", IncludeEqual: false}},
	)
	s.NoError(err)
	s.Len(data, 0)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: false},
		GreaterThan: &bst.Bound[string]{Value: "Zara", IncludeEqual: true}},
	)
	s.NoError(err)
	s.Len(data, 0)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: true},
		GreaterThan: &bst.Bound[string]{Value: "Zara", IncludeEqual: true}},
	)
	s.NoError(err)
	s.Len(data, 0)

	// same value, non-inclusive query
	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: false},
		GreaterThan: &bst.Bound[string]{Value: "Iris", IncludeEqual: false}},
	)
	s.NoError(err)
	s.Len(data, 0)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: true},
		GreaterThan: &bst.Bound[string]{Value: "Iris", IncludeEqual: false}},
	)
	s.NoError(err)
	s.Len(data, 0)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: false},
		GreaterThan: &bst.Bound[string]{Value: "Iris", IncludeEqual: true}},
	)
	s.NoError(err)
	s.Len(data, 0)

	// same value, inclusive query
	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: true},
		GreaterThan: &bst.Bound[string]{Value: "Iris", IncludeEqual: true}},
	)
	s.NoError(err)
	s.Equal([]int{33}, data)
}

func (s *BSTTestSuite) TestQuery() {
	data, err := s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: false},
		GreaterThan: &bst.Bound[string]{Value: "Alice", IncludeEqual: false}},
	)
	s.NoError(err)
	s.Equal([]int{63, 55, 88}, data)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: true},
		GreaterThan: &bst.Bound[string]{Value: "Alice", IncludeEqual: false}},
	)
	s.NoError(err)
	s.Equal([]int{63, 55, 88, 33}, data)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: false},
		GreaterThan: &bst.Bound[string]{Value: "Alice", IncludeEqual: true}},
	)
	s.NoError(err)
	s.Equal([]int{42, 23, 63, 55, 88}, data)

	data, err = s.b.Query(bst.Query[string]{
		LowerThan:   &bst.Bound[string]{Value: "Iris", IncludeEqual: true},
		GreaterThan: &bst.Bound[string]{Value: "Alice", IncludeEqual: true}},
	)
	s.NoError(err)
	s.Equal([]int{42, 23, 63, 55, 88, 33}, data)
}

func (s *BSTTestSuite) TestUpdate() {
	s.NoError(s.b.Update("Leo", 76, 10))

	s.Equal([]int{10}, s.b.Node.Values)

	s.NoError(s.b.Update("Leo", 12, 1000))

	s.Equal([]int{10}, s.b.Node.Values)
}

func TestBSTTestSuite(t *testing.T) {
	suite.Run(t, new(BSTTestSuite))
}
