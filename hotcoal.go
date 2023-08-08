// Package hotcoal secures your handcrafted SQL against injection
package hotcoal

// hotcoalString is an abstract data type, which is used for handcrafting SQL,
// protecting against SQL injection
type hotcoalString string

// Slice is an alias for a slice of hotcoalStrings.
// Since hotcoalString is not exported, we export this alias,
// which allows you to create slices.
type Slice = []hotcoalString

// The String method converts a hotcoalString to a plain string.
// Please do all your SQL handcrafting using hotcoalStrings,
// and convert the result to a plain string only when you pass it to the
// SQL library.
func (s hotcoalString) String() string {
	return string(s)
}

// The Wrap function converts an untyped string constant to a hotcoalString.
// You can only use it with an untyped string constant, not with a string variable.
// For the latter, please use an Allowlist to validate the variable and guard
// against SQL injection.
func Wrap(s hotcoalString) hotcoalString {
	return s
}

// The W function is a shorthand for Wrap
func W(s hotcoalString) hotcoalString {
	return s
}
