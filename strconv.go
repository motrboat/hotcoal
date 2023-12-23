package hotcoal

import "strconv"

// The Itoa function converts an int to a hotcoalString.
func Itoa(i int) hotcoalString {
	return hotcoalString(strconv.Itoa(i))
}
