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

https://pkg.go.dev/github.com/motrboat/hotcoal

## Disclaimer

Hotcoal comes without any warranty.
