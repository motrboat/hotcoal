package hotcoal

import "fmt"

// allowlistT holds an allowlist of items, which is used to validate string variables
// such as column names or table names, guarding against SQL injection
type allowlistT struct {
	items map[hotcoalString]unitT
}

type unitT = struct{}

var unit = unitT{}

// Allowlist creates an allowlistT, which is used to validate validate string variables
// such as column names or table names, guarding against SQL injection
func Allowlist(firstAllowlistItem hotcoalString, otherAllowlistItems ...hotcoalString) allowlistT {
	ret := allowlistT{
		items: map[hotcoalString]unitT{
			firstAllowlistItem: unit,
		},
	}

	for _, el := range otherAllowlistItems {
		ret.items[el] = unit
	}

	return ret
}

// The Validate method validates a string variable against the allowlist and returns a hotcoalString.
// If the value is not in the allowlist, it returns an error.
func (a allowlistT) Validate(value string) (hotcoalString, error) {
	if _, ok := a.items[hotcoalString(value)]; ok {
		return hotcoalString(value), nil
	}

	return "", fmt.Errorf("Hotcoal validation error - value %#v is not in allowlist %#v", value, a.items)
}

// The V method is an shorthand for Validate
func (a allowlistT) V(value string) (hotcoalString, error) {
	return a.Validate(value)
}

// The Validate method validates a string variable against the allowlist and returns a hotcoalString.
// If the value is not in the allowlist, it panics.
func (a allowlistT) MustValidate(value string) hotcoalString {
	ret, err := a.Validate(value)
	if err != nil {
		panic(err)
	}

	return ret
}

// The MV method is an shorthand for MustValidate
func (a allowlistT) MV(value string) hotcoalString {
	return a.MustValidate(value)
}
