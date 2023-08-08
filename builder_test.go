package hotcoal

import "testing"

func TestBuilder(t *testing.T) {
  var b Builder

  if b.Grow(64); b.Len() != 0 || b.Cap() != 64 {
    t.Fail()
  }

  b.Write("foo")
  b.Write("bar").Write("tar")

  if "foobartar" != b.String() {
    t.Fail()
  }

  if "foobartar" != b.HotcoalString().String() {
    t.Fail()
  }

  if b.Len() != 9 || b.Cap() != 64 {
    t.Fail()
  }

  b.Reset()

  if b.Len() != 0 || b.Cap() != 0 || "" != b.String() {
    t.Fail()
  }
}
