package hotcoal

import (
	"testing"
)

func TestAllowlist(t *testing.T) {
	allowlist := Allowlist("foo", "bar", "tar")

	for _, el := range []string{"foo", "bar", "tar"} {
		hs, err := allowlist.Validate(el)
		if el != hs.String() || err != nil {
			t.Fail()
		}

		hs, err = allowlist.V(el)
		if el != hs.String() || err != nil {
			t.Fail()
		}

		hs = allowlist.MustValidate(el)
		if el != hs.String() {
			t.Fail()
		}

		hs = allowlist.MV(el)
		if el != hs.String() {
			t.Fail()
		}
	}

	el := "baz"

	hs, err := allowlist.Validate(el)
	if "" != hs.String() || err == nil {
		t.Fail()
	}

	hs, err = allowlist.V(el)
	if "" != hs.String() || err == nil {
		t.Fail()
	}

	func() {
		var hs hotcoalString

		defer func() {
			p := recover()
			if "" != hs.String() || p == nil {
				t.Fail()
			}
		}()

		hs = allowlist.MustValidate(el)
	}()

	func() {
		var hs hotcoalString

		defer func() {
			p := recover()
			if "" != hs.String() || p == nil {
				t.Fail()
			}
		}()

		hs = allowlist.MV(el)
	}()
}
