package e

import "fmt"

// Wrap takes a custom message and an existing error, and combines them into a single error.
// The returned error is formatted as "custom message: original error".
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %s", msg, err)
}
