package nocompile

import "github.com/motrboat/hotcoal"


var x = "foo"

var _ = hotcoal.Wrap(x)                          // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Wrap

var _ = hotcoal.W(x)                             // ERROR

var y = hotcoal.Wrap("bar")                      // OK

var _ = hotcoal.W("bar")                         // OK

var _ = hotcoal.W(hotcoal.W("bar"))              // OK

var _ = hotcoal.Wrap("foo") + "bar"              // This is ok, bar gets converted to a hotcoalString
                                                 // but we don't encourage you to use it, because it's less safe
                                                 // Please wrap all your string constants in hotcoalStrings

var z = hotcoal.Slice{}                          // ok

var _ = append(z, x)                             // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to append

var _ = append(z, y)                             // OK

var _ = hotcoal.W(hotcoal.Join(z, y))            // OK

var _ = hotcoal.Join([]string{}, y)              // ERROR: cannot use []string{} (value of type []string) as []hotcoal.hotcoalString value in argument to hotcoal.Join

var _ = hotcoal.Join(z, x)                       // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Join

var _ = hotcoal.W(hotcoal.Replace(y, y, y, 0))   // OK

var _ = hotcoal.Replace(x, y, y, 0)              // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Replace

var _ = hotcoal.Replace(y, x, y, 0)              // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Replace

var _ = hotcoal.Replace(y, y, x, 0)              // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Replace

var _ = hotcoal.W(hotcoal.ReplaceAll(y, y, y))   // OK

var _ = hotcoal.ReplaceAll(x, y, y)              // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.ReplaceAll

var _ = hotcoal.ReplaceAll(y, x, y)              // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.ReplaceAll

var _ = hotcoal.ReplaceAll(y, y, x)              // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.ReplaceAll

var b hotcoal.Builder

var _ = b.Write(x)                               // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to b.Write

var _ = b.Write(y)                               // OK

var _ = hotcoal.W(b.HotcoalString())             // OK

var _ = hotcoal.W(b.String())                    // ERROR: cannot use b.String() (value of type string) as hotcoal.hotcoalString value in argument to hotcoal.W

var _ = hotcoal.Allowlist(x)                     // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Allowlist

var t = hotcoal.Allowlist(y)                     // OK

var _ = hotcoal.Allowlist(y, x)                  // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Allowlist

var _ = hotcoal.Allowlist(y, y)                  // OK

var _ = hotcoal.Allowlist(y, y, x)               // ERROR: cannot use x (variable of type string) as hotcoal.hotcoalString value in argument to hotcoal.Allowlist

var _ = hotcoal.Allowlist(y, y, y)               // OK

func f() {
  var _, _ = t.Validate("bar")                   // OK

  var _, _ = t.V("bar")                          // OK

  var _ = t.MustValidate("bar")                  // OK

  var _ = t.MV("bar")                            // OK
}
