package hotcoal

import "testing"

func TestJoin(t *testing.T) {
  if "foo-bar-tar" != Join(Slice{"foo", "bar", "tar"}, "-").String() {
    t.Fail()
  }
}

func TestReplace(t *testing.T) {
  if "bbbaa" != W("aaaaa").Replace("a", "b", 3).String() {
    t.Fail()
  }
}

func TestReplaceAll(t *testing.T) {
  if "bbbbb" != W("aaaaa").ReplaceAll("a", "b").String() {
    t.Fail()
  }
}
