package user

import "fmt"

// ErrNotFound ...
type ErrNotFound struct {
	Email string
}

// Error ...
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("user %s not found", e.Email)
}
