package core

// Ptr is a helper function to create a pointer
// to a value.
func Ptr[T any](v T) *T {
	return &v
}
