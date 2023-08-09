package hotcoal

import (
	"fmt"
	"strings"
)

// A Builder is used to efficiently build a hotcoalString using the Write method.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Builder.
//
// Under the hood, it uses strings.Builder https://pkg.go.dev/strings#Builder
type Builder struct {
	stringBuilder strings.Builder
}

// Cap returns the capacity of the builder's underlying byte slice. It is the
// total space allocated for the hotcoalString being built and includes any bytes
// already written.
func (b *Builder) Cap() int {
	return b.stringBuilder.Cap()
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *Builder) Grow(n int) {
	b.stringBuilder.Grow(n)
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *Builder) Len() int {
	return b.stringBuilder.Len()
}

// Reset resets the Builder to be empty.
func (b *Builder) Reset() {
	b.stringBuilder.Reset()
}

// Write appends the contents of s to b's buffer. It returns b, you can chain method calls.
func (b *Builder) Write(s hotcoalString) *Builder {
	_, err := b.stringBuilder.WriteString(string(s))

	if err != nil {
		// this shouldn't happen, the strings.Builder returns nil error
		panic(fmt.Sprintf("Hotcoal Builder.Write received error: %#v", err))
	}

	return b
}

// String returns the accumulated string as a hotcoalString.
func (b *Builder) HotcoalString() hotcoalString {
	return hotcoalString(b.String())
}

// String returns the accumulated string as a plain string.
func (b *Builder) String() string {
	return b.stringBuilder.String()
}
