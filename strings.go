package hotcoal

import "strings"

// Join concatenates the elements of its first argument to create a single hotcoalString. The separator
// hotcoalString sep is placed between elements in the resulting hotcoalString.
//
// Under the hood, it uses strings.Join https://pkg.go.dev/strings#Join
func Join(elems []hotcoalString, sep hotcoalString) hotcoalString {
	stringElems := make([]string, 0, len(elems))
	for _, el := range elems {
		stringElems = append(stringElems, string(el))
	}

	stringSep := string(sep)

	result := strings.Join(stringElems, stringSep)

	return hotcoalString(result)
}

// Replace returns a copy of the hotcoalString s with the first n
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the hotcoalString
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune hotcoalString.
// If n < 0, there is no limit on the number of replacements.
//
// Under the hood, it uses strings.Replace https://pkg.go.dev/strings#Replace
func Replace(s, old, new hotcoalString, n int) hotcoalString {
	result := strings.Replace(
		string(s),
		string(old),
		string(new),
		n,
	)

	return hotcoalString(result)
}

// ReplaceAll returns a copy of the string s with all
// non-overlapping instances of old replaced by new.
// If old is empty, it matches at the beginning of the string
// and after each UTF-8 sequence, yielding up to k+1 replacements
// for a k-rune string.
//
// Under the hood, it uses strings.ReplaceAll https://pkg.go.dev/strings#ReplaceAll
func ReplaceAll(s, old, new hotcoalString) hotcoalString {
	result := strings.ReplaceAll(
		string(s),
		string(old),
		string(new),
	)

	return hotcoalString(result)
}
