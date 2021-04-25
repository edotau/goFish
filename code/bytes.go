package code

import ()

// The Slice interface reflects the built-in slice type behavior.
type Slice interface {
	Make(len, cap int) Slice
	Len() int
	Cap() int
	Slice(start, end int) Slice
	Append(src Slice) Slice
	Copy(src Slice) int
}
