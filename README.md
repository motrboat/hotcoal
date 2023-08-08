# Hotcoal

## Secure your handcrafted SQL against injection

Many Golang apps do not use an ORM, and instead prefer to use raw SQL for their database operations.

Even though writing raw SQL can be sometimes more flexible, efficient and clean than using an ORM, it can make you vulnerable to SQL injection.

While most database technologies (PostreSQL, MySQL, Sqlite etc.) and Golang SQL libraries (`database/sql`, `sqlx` etc.) allow you to use prepared statements, which take parameters and sanitize them, there are some limitations:
 * you can use parameters for values, but not for column names or table names
 * for more complex stuff, a good ORM can generate SQL dynamically for you, and do it safely; but you aren't using an ORM, so you start handcrafting SQL ...

This is where Hotcoal comes in.

## How it works

Hotcoal provides a **minimal** API which helps you guard your handcrafted SQL against SQL injection.

1. Handcraft your SQL using Hotcoal instead of plain strings. This way all your parameters are validated against the allowlists you provide.

2. Convert the result to a plain string and use it with your SQL library.

You can use Hotcoal with any SQL library you like.

Hotcoal does not replace prepared queries with parameters, but complements it.

## Examples

We use this example table:

```sql
CREATE TABLE "users" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"first_name"	TEXT NOT NULL,
	"middle_name"	TEXT NOT NULL,
	"last_name"	TEXT NOT NULL,
	"nickname"	TEXT,
	"email"	TEXT NOT NULL
);
```

### Validate column name

We use Hotcoal to validate the column name and create a SQL query string.

String constants can be converted to a hotcoalString directly, without validation, using `Wrap`.

But string variables cannot be converted to a hotcoalString directly. This helps protect against SQL injection.

Instead, you must validate the string variable against an allowlist, using `Allowlist` and `Validate`. <br>
If you try to convert a string variable to a hotcoalString without validation,
using `Wrap`, it doesn't compile.


```golang
import (
  "database/sql"
  "github.com/motrboat/hotcoal"  
)

func queryCount(db *sql.DB, columnName string, value string) *sql.Row {

  // validating the columnName variable against an allowlist
  // it returns a hotcoalString
  validatedColumnName, err := hotcoal.
    Allowlist("first_name", "middle_name", "last_name").
    Validate(columnName)

  if err != nil {
    return nil, err
  }  

  // handcraft SQL using hotcoalStrings
  query := hotcoal.Wrap("SELECT COUNT(*) FROM users WHERE ") +
    validatedColumnName +
    hotcoal.Wrap(" = ?;")

  // if you try to use a regular string variable, which has not been validated,
  // it's vulnerable to SQL injection, so it doesn't compile:
  //       query := hotcoal.Wrap("SELECT COUNT(*) FROM users WHERE ") +
  //         columnName +
  //         hotcoal.Wrap(" = ?;")

  // convert the hotcoalString back to a regular string only when it is passed to the DB library
  // use prepared query with parameter to sanitize the value
  row := db.QueryRow(query.String(), value)

  return row
}
```

### Validate handcrafted SQL

To handcraft more complicated SQL, you can also use the hotcoal equivalents
of `strings` functions `Join`, `Replace` and `ReplaceAll`.
These can be called with hotcoalStrings, or string constants.
But not string variables - those must be validated first.

```golang
import (
  "database/sql"
  "github.com/motrboat/hotcoal"  
)

type Filter struct {
  ColumnName string
  Value string
}

func queryCount(db *sql.DB, filters []Filter) *sql.Row {
  query := hotcoal.Wrap("SELECT COUNT(*) FROM users WHERE {{FILTERS}};")
  allowlist := hotcoal.Allowlist("first_name", "middle_name", "last_name", "nickname")

  sqlArr := hotcoal.Slice{}
  values := []string{}

  for _, filter := range filters {
    validatedColumnName, err := allowlist.Validate(filter.ColumnName)
    if err != nil {
      return nil, err
    }

    sqlArr = append(sqlArr, validatedColumnName + hotcoal.Wrap(" = ?"))
    values = append(values, filter.Value)
  }

  query = hotcoal.ReplaceAll(
    query,
    "{{FILTERS}}",
    hotcoal.Join(sqlArr, " OR "),
  )

  row := db.QueryRow(query.String(), values...)

  return row
}
```

### String builder

Hotcoal also offers a version of `strings.Builder` using `hotcoalStrings`. It minimizes memory copying and is more efficient.

```golang
import (
  "database/sql"
  "github.com/motrboat/hotcoal"  
)

type Filter struct {
  ColumnName string
  Value string
}

func queryCount(db *sql.DB, filters []Filter) *sql.Row {
  allowlist := hotcoal.Allowlist("first_name", "middle_name", "last_name", "nickname")

  builder := hotcoal.Builder{}
  values := []string{}

  builder.Write("SELECT COUNT(*) FROM users WHERE ")

  for i, filter := range filters {
    if i != 0 {
      builder.Write(" OR ")
    }

    validatedColumnName, err := allowlist.Validate(filter.ColumnName)
    if err != nil {
      return nil, err
    }

    builder.Write(validatedColumnName + hotcoal.Wrap(" = ?"))
    values = append(values, filter.Value)
  }

  builder.Write(";")

  row := db.QueryRow(builder.String(), values...)

  return row
}
```

## Documentation

- [type hotcoalString](<#hotcoalString>)
  - [func Wrap\(s hotcoalString\) hotcoalString](<#Wrap>)
  - [func W\(s hotcoalString\) hotcoalString](<#W>)
  - [func \(s hotcoalString\) String\(\) string](<#hotcoalString.String>)
- [type Slice](<#Slice>)
- [type allowlistT](<#allowlistT>)
  - [func Allowlist\(firstAllowlistItem hotcoalString, otherAllowlistItems ...hotcoalString\) allowlistT](<#Allowlist>)
  - [func \(a allowlistT\) Validate\(value string\) \(hotcoalString, error\)](<#allowlistT.Validate>)
  - [func \(a allowlistT\) V\(value string\) \(hotcoalString, error\)](<#allowlistT.V>)
  - [func \(a allowlistT\) MustValidate\(value string\) hotcoalString](<#allowlistT.MustValidate>)
  - [func \(a allowlistT\) MV\(value string\) hotcoalString](<#allowlistT.MV>)
- [func Join\(elems \[\]hotcoalString, sep hotcoalString\) hotcoalString](<#Join>)
- [func Replace\(s, old, new hotcoalString, n int\) hotcoalString](<#Replace>)
- [func ReplaceAll\(s, old, new hotcoalString\) hotcoalString](<#ReplaceAll>)
- [type Builder](<#Builder>)
  - [func \(b \*Builder\) Cap\(\) int](<#Builder.Cap>)
  - [func \(b \*Builder\) Grow\(n int\)](<#Builder.Grow>)
  - [func \(b \*Builder\) Len\(\) int](<#Builder.Len>)
  - [func \(b \*Builder\) Reset\(\)](<#Builder.Reset>)
  - [func \(b \*Builder\) Write\(s hotcoalString\) \*Builder](<#Builder.Write>)
  - [func \(b \*Builder\) HotcoalString\(\) hotcoalString](<#Builder.HotcoalString>)
  - [func \(b \*Builder\) String\(\) string](<#Builder.String>)

<a name="hotcoalString"></a>
### type [hotcoalString](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L6>)

hotcoalString is an abstract data type, which is used for handcrafting SQL, protecting against SQL injection

```go
type hotcoalString string
```

<a name="Wrap"></a>
#### func [Wrap](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L25>)

```go
func Wrap(s hotcoalString) hotcoalString
```

The Wrap function converts an untyped string constant to a hotcoalString. You can only use it with an untyped string constant, not with a string variable. For the latter, please use an Allowlist to validate the variable and guard against SQL injection.

<a name="W"></a>
#### func [W](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L30>)

```go
func W(s hotcoalString) hotcoalString
```

The W function is a shorthand for Wrap

<a name="hotcoalString.String"></a>
#### func \(hotcoalString\) [String](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L17>)

```go
func (s hotcoalString) String() string
```

The String method converts a hotcoalString to a plain string. Please do all your SQL handcrafting using hotcoalStrings, and convert the result to a plain string only when you pass it to the SQL library.

<a name="Slice"></a>
### type [Slice](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L11>)

Slice is an alias for a slice of hotcoalStrings. Since hotcoalString is not exported, we export this alias, which allows you to create slices.

```go
type Slice = []hotcoalString
```

<a name="allowlistT"></a>
### type [allowlistT](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L7-L9>)

allowlistT holds an allowlist of items, which is used to validate string variables such as column names or table names, guarding against SQL injection

```go
type allowlistT struct {
    items map[hotcoalString]unitT
}
```

<a name="Allowlist"></a>
#### func [Allowlist](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L17>)

```go
func Allowlist(firstAllowlistItem hotcoalString, otherAllowlistItems ...hotcoalString) allowlistT
```

Allowlist creates an allowlistT, which is used to validate validate string variables such as column names or table names, guarding against SQL injection

<a name="allowlistT.Validate"></a>
#### func \(allowlistT\) [Validate](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L33>)

```go
func (a allowlistT) Validate(value string) (hotcoalString, error)
```

The Validate method validates a string variable against the allowlist and returns a hotcoalString. If the value is not in the allowlist, it returns an error.

<a name="allowlistT.V"></a>
#### func \(allowlistT\) [V](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L42>)

```go
func (a allowlistT) V(value string) (hotcoalString, error)
```

The V method is an shorthand for Validate

<a name="allowlistT.MustValidate"></a>
#### func \(allowlistT\) [MustValidate](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L48>)

```go
func (a allowlistT) MustValidate(value string) hotcoalString
```

The MustValidate method validates a string variable against the allowlist and returns a hotcoalString. If the value is not in the allowlist, it panics.

<a name="allowlistT.MV"></a>
#### func \(allowlistT\) [MV](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L58>)

```go
func (a allowlistT) MV(value string) hotcoalString
```

The MV method is an shorthand for MustValidate

<a name="Join"></a>
### func [Join](<https://github.com/motrboat/hotcoal/blob/main/strings.go#L9>)

```go
func Join(elems []hotcoalString, sep hotcoalString) hotcoalString
```

Join concatenates the elements of its first argument to create a single hotcoalString. The separator hotcoalString sep is placed between elements in the resulting hotcoalString.

Under the hood, it uses strings.Join https://pkg.go.dev/strings#Join

<a name="Replace"></a>
### func [Replace](<https://github.com/motrboat/hotcoal/blob/main/strings.go#L30>)

```go
func Replace(s, old, new hotcoalString, n int) hotcoalString
```

Replace returns a copy of the hotcoalString s with the first n non\-overlapping instances of old replaced by new. If old is empty, it matches at the beginning of the hotcoalString and after each UTF\-8 sequence, yielding up to k\+1 replacements for a k\-rune hotcoalString. If n \< 0, there is no limit on the number of replacements.

Under the hood, it uses strings.Replace https://pkg.go.dev/strings#Replace

<a name="ReplaceAll"></a>
### func [ReplaceAll](<https://github.com/motrboat/hotcoal/blob/main/strings.go#L48>)

```go
func ReplaceAll(s, old, new hotcoalString) hotcoalString
```

ReplaceAll returns a copy of the string s with all non\-overlapping instances of old replaced by new. If old is empty, it matches at the beginning of the string and after each UTF\-8 sequence, yielding up to k\+1 replacements for a k\-rune string.

Under the hood, it uses strings.ReplaceAll https://pkg.go.dev/strings#ReplaceAll

<a name="Builder"></a>
### type [Builder](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L13-L15>)

A Builder is used to efficiently build a hotcoalString using the Write method. It minimizes memory copying. The zero value is ready to use. Do not copy a non\-zero Builder.

Under the hood, it uses strings.Builder https://pkg.go.dev/strings#Builder

```go
type Builder struct {
    stringBuilder strings.Builder
}
```

<a name="Builder.Cap"></a>
#### func \(\*Builder\) [Cap](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L20>)

```go
func (b *Builder) Cap() int
```

Cap returns the capacity of the builder's underlying byte slice. It is the total space allocated for the hotcoalString being built and includes any bytes already written.

<a name="Builder.Grow"></a>
#### func \(\*Builder\) [Grow](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L27>)

```go
func (b *Builder) Grow(n int)
```

Grow grows b's capacity, if necessary, to guarantee space for another n bytes. After Grow\(n\), at least n bytes can be written to b without another allocation. If n is negative, Grow panics.

<a name="Builder.Len"></a>
#### func \(\*Builder\) [Len](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L32>)

```go
func (b *Builder) Len() int
```

Len returns the number of accumulated bytes; b.Len\(\) == len\(b.String\(\)\).

<a name="Builder.Reset"></a>
#### func \(\*Builder\) [Reset](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L37>)

```go
func (b *Builder) Reset()
```

Reset resets the Builder to be empty.

<a name="Builder.Write"></a>
#### func \(\*Builder\) [Write](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L42>)

```go
func (b *Builder) Write(s hotcoalString) *Builder
```

Write appends the contents of s to b's buffer. It returns b.

<a name="Builder.HotcoalString"></a>
#### func \(\*Builder\) [HotcoalString](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L54>)

```go
func (b *Builder) HotcoalString() hotcoalString
```

String returns the accumulated string as a hotcoalString.

<a name="Builder.String"></a>
#### func \(\*Builder\) [String](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L59>)

```go
func (b *Builder) String() string
```

String returns the accumulated string as a plain string.

## Disclaimer

Hotcoal comes without any warranty.

