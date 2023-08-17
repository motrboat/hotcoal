package hotcoal

import "testing"

func TestString(t *testing.T) {
	if "foo" != Wrap("foo").String() {
		t.Fail()
	}

	if "foo" != W("foo").String() {
		t.Fail()
	}
}
