package hotcoal

import "testing"

func TestItoa(t *testing.T) {
	if "123" != Itoa(123).String() {
		t.Fail()
	}
}
