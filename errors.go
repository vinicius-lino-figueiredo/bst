package bst

import "fmt"

// ErrViolated represents a duplicate key when the bst is unique.
type ErrViolated struct {
	key any
}

// Error implements error.
func (e ErrViolated) Error() string {
	return fmt.Sprintf("Can't insert key %v, it violates the unique constraint", e.key)
}
