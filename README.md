# Hotcoal

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/motrboat/hotcoal.svg)](https://pkg.go.dev/github.com/motrboat/hotcoal)
[![Go Report Card](https://goreportcard.com/badge/github.com/motrboat/hotcoal)](https://goreportcard.com/report/github.com/motrboat/hotcoal)
![Coverage](https://img.shields.io/badge/Coverage-97.2%25-brightgreen)

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

func queryCount(db *sql.DB, tableName string, filters []Filter) (*sql.Row, error) {
  validatedTableName, err := hotcoal.
    Allowlist("users", "customers").
    Validate(tableName)
  if err != nil {
    return nil, err
  }

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

  query := hotcoal.Wrap("SELECT COUNT(*) FROM {{TABLE}} WHERE {{FILTERS}};").
    ReplaceAll(
      "{{TABLE}}",
      validatedTableName,
    ).
    ReplaceAll(
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

- [type hotcoalString](<#type-hotcoalstring>)
  - [func Wrap\(s hotcoalString\) hotcoalString](<#func-wrap>)
  - [func W\(s hotcoalString\) hotcoalString](<#func-w>)
  - [func \(s hotcoalString\) String\(\) string](<#func-hotcoalstring-string>)
- [type Slice](<#type-slice>)
- [type allowlistT](<#type-allowlistt>)
  - [func Allowlist\(firstAllowlistItem hotcoalString, otherAllowlistItems ...hotcoalString\) allowlistT](<#func-allowlist>)
  - [func \(a allowlistT\) Validate\(value string\) \(hotcoalString, error\)](<#func-allowlistt-validate>)
  - [func \(a allowlistT\) V\(value string\) \(hotcoalString, error\)](<#func-allowlistt-v>)
  - [func \(a allowlistT\) MustValidate\(value string\) hotcoalString](<#func-allowlistt-mustvalidate>)
  - [func \(a allowlistT\) MV\(value string\) hotcoalString](<#func-allowlistt-mv>)
- [func Join\(elems \[\]hotcoalString, sep hotcoalString\) hotcoalString](<#func-join>)
- [func \(s hotcoalString\) Replace\(old, new hotcoalString, n int\) hotcoalString](<#func-hotcoalstring-replace>)
- [func \(s hotcoalString\) ReplaceAll\(old, new hotcoalString\) hotcoalString](<#func-hotcoalstring-replaceall>)
- [type Builder](<#type-builder>)
  - [func \(b \*Builder\) Cap\(\) int](<#func-builder-cap>)
  - [func \(b \*Builder\) Grow\(n int\)](<#func-builder-grow>)
  - [func \(b \*Builder\) Len\(\) int](<#func-builder-len>)
  - [func \(b \*Builder\) Reset\(\)](<#func-builder-reset>)
  - [func \(b \*Builder\) Write\(s hotcoalString\) \*Builder](<#func-builder-write>)
  - [func \(b \*Builder\) HotcoalString\(\) hotcoalString](<#func-builder-hotcoalstring>)
  - [func \(b \*Builder\) String\(\) string](<#func-builder-string>)


### type [hotcoalString](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L6>)

hotcoalString is an abstract data type, which is used for handcrafting SQL, protecting against SQL injection

```go
type hotcoalString string
```


#### func [Wrap](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L25>)

```go
func Wrap(s hotcoalString) hotcoalString
```

The Wrap function converts an untyped string constant to a hotcoalString. You can only use it with an untyped string constant, not with a string variable. For the latter, please use an Allowlist to validate the variable and guard against SQL injection.


#### func [W](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L30>)

```go
func W(s hotcoalString) hotcoalString
```

The W function is a shorthand for Wrap


#### func \(hotcoalString\) [String](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L17>)

```go
func (s hotcoalString) String() string
```

The String method converts a hotcoalString to a plain string. Please do all your SQL handcrafting using hotcoalStrings, and convert the result to a plain string only when you pass it to the SQL library.


### type [Slice](<https://github.com/motrboat/hotcoal/blob/main/hotcoal.go#L11>)

Slice is an alias for a slice of hotcoalStrings. Since hotcoalString is not exported, we export this alias, which allows you to create slices.

```go
type Slice = []hotcoalString
```


### type [allowlistT](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L7-L9>)

allowlistT holds an allowlist of items, which is used to validate string variables such as column names or table names, guarding against SQL injection

```go
type allowlistT struct {
    items map[hotcoalString]unitT
}
```


#### func [Allowlist](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L17>)

```go
func Allowlist(firstAllowlistItem hotcoalString, otherAllowlistItems ...hotcoalString) allowlistT
```

Allowlist creates an allowlistT, which is used to validate validate string variables such as column names or table names, guarding against SQL injection


#### func \(allowlistT\) [Validate](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L33>)

```go
func (a allowlistT) Validate(value string) (hotcoalString, error)
```

The Validate method validates a string variable against the allowlist and returns a hotcoalString. If the value is not in the allowlist, it returns an error.


#### func \(allowlistT\) [V](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L42>)

```go
func (a allowlistT) V(value string) (hotcoalString, error)
```

The V method is an shorthand for Validate


#### func \(allowlistT\) [MustValidate](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L48>)

```go
func (a allowlistT) MustValidate(value string) hotcoalString
```

The MustValidate method validates a string variable against the allowlist and returns a hotcoalString. If the value is not in the allowlist, it panics.


#### func \(allowlistT\) [MV](<https://github.com/motrboat/hotcoal/blob/main/allowlist.go#L58>)

```go
func (a allowlistT) MV(value string) hotcoalString
```

The MV method is an shorthand for MustValidate


### func [Join](<https://github.com/motrboat/hotcoal/blob/main/strings.go#L9>)

```go
func Join(elems []hotcoalString, sep hotcoalString) hotcoalString
```

Join concatenates the elements of its first argument to create a single hotcoalString. The separator hotcoalString sep is placed between elements in the resulting hotcoalString.

Under the hood, it uses strings.Join https://pkg.go.dev/strings#Join


### func \(hotcoalString\) [Replace](<https://github.com/motrboat/hotcoal/blob/main/strings.go#L68>)

```go
func (s hotcoalString) Replace(old, new hotcoalString, n int) hotcoalString
```

Replace returns a copy of the hotcoalString s with the first n non\-overlapping instances of old replaced by new. If old is empty, it matches at the beginning of the hotcoalString and after each UTF\-8 sequence, yielding up to k\+1 replacements for a k\-rune hotcoalString. If n \< 0, there is no limit on the number of replacements.

You can chain method calls.

Under the hood, it uses strings.Replace https://pkg.go.dev/strings#Replace

### func \(hotcoalString\) [ReplaceAll](<https://github.com/motrboat/hotcoal/blob/main/strings.go#L88>)

```go
func (s hotcoalString) ReplaceAll(old, new hotcoalString) hotcoalString
```

ReplaceAll returns a copy of the string s with all non\-overlapping instances of old replaced by new. If old is empty, it matches at the beginning of the string and after each UTF\-8 sequence, yielding up to k\+1 replacements for a k\-rune string.

You can chain method calls.

Under the hood, it uses strings.ReplaceAll https://pkg.go.dev/strings#ReplaceAll


### type [Builder](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L13-L15>)

A Builder is used to efficiently build a hotcoalString using the Write method. It minimizes memory copying. The zero value is ready to use. Do not copy a non\-zero Builder.

Under the hood, it uses strings.Builder https://pkg.go.dev/strings#Builder

```go
type Builder struct {
    stringBuilder strings.Builder
}
```


#### func \(\*Builder\) [Cap](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L20>)

```go
func (b *Builder) Cap() int
```

Cap returns the capacity of the builder's underlying byte slice. It is the total space allocated for the hotcoalString being built and includes any bytes already written.


#### func \(\*Builder\) [Grow](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L27>)

```go
func (b *Builder) Grow(n int)
```

Grow grows b's capacity, if necessary, to guarantee space for another n bytes. After Grow\(n\), at least n bytes can be written to b without another allocation. If n is negative, Grow panics.


#### func \(\*Builder\) [Len](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L32>)

```go
func (b *Builder) Len() int
```

Len returns the number of accumulated bytes; b.Len\(\) == len\(b.String\(\)\).


#### func \(\*Builder\) [Reset](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L37>)

```go
func (b *Builder) Reset()
```

Reset resets the Builder to be empty.


#### func \(\*Builder\) [Write](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L42>)

```go
func (b *Builder) Write(s hotcoalString) *Builder
```

Write appends the contents of s to b's buffer. It returns b, you can chain method calls.


#### func \(\*Builder\) [HotcoalString](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L54>)

```go
func (b *Builder) HotcoalString() hotcoalString
```

String returns the accumulated string as a hotcoalString.


#### func \(\*Builder\) [String](<https://github.com/motrboat/hotcoal/blob/main/builder.go#L59>)

```go
func (b *Builder) String() string
```

String returns the accumulated string as a plain string.

## Disclaimer

Hotcoal comes without any warranty.
