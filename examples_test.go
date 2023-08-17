package hotcoal_test

import (
	"github.com/motrboat/hotcoal"
	"testing"
)

func TestValidateColumnName(t *testing.T) {
	sql, err := validateColumnName("middle_name", "Larry")
	if err != nil || sql != "SELECT COUNT(*) FROM users WHERE middle_name = ?;" {
		t.Fail()
	}

	_, err = validateColumnName("true; DROP TABLE users; --", "Larry")
	if err == nil {
		t.Fail()
	}
}

func validateColumnName(columnName, value string) (string, error) {
	validatedColumnName, err := hotcoal.
		Allowlist("first_name", "middle_name", "last_name").
		Validate(columnName)

	if err != nil {
		return "", err
	}

	query := hotcoal.Wrap("SELECT COUNT(*) FROM users WHERE ") +
		validatedColumnName +
		hotcoal.Wrap(" = ?;")

	return query.String(), nil
}

func TestHandcraftedSQL(t *testing.T) {
	sql, err := handcraftSQL(
		"users",
		Filter{"first_name", "John"},
		Filter{"last_name", "Doe"},
	)

	if err != nil || sql != "SELECT COUNT(*) FROM users WHERE first_name = ? OR last_name = ?;" {
		t.Fail()
	}

	_, err = handcraftSQL(
		"users",
		Filter{"first_name", "John"},
		Filter{"true; DROP TABLE users; --", "Doe"},
	)

	if err == nil {
		t.Fail()
	}
}

type Filter struct {
	ColumnName string
	Value      string
}

func handcraftSQL(tableName string, filters ...Filter) (string, error) {
	validatedTableName, err := hotcoal.Allowlist("users", "customers").V(tableName)
	if err != nil {
		return "", err
	}

	allowlist := hotcoal.Allowlist("first_name", "middle_name", "last_name", "nickname")

	sqlArr := hotcoal.Slice{}
	values := []string{}

	for _, filter := range filters {
		validatedColumnName, err := allowlist.Validate(filter.ColumnName)
		if err != nil {
			return "", err
		}

		sqlArr = append(sqlArr, validatedColumnName+hotcoal.Wrap(" = ?"))
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

	return query.String(), nil
}

func TestBuilder(t *testing.T) {
	sql, err := build(
		Filter{"middle_name", "Larry"},
		Filter{"nickname", "Larry"},
	)

	if err != nil || sql != "SELECT COUNT(*) FROM users WHERE middle_name = ? OR nickname = ?;" {
		t.Fail()
	}

	_, err = build(
		Filter{"middle_name", "Larry"},
		Filter{"true; DROP TABLE users; --", "Larry"},
	)

	if err == nil {
		t.Fail()
	}
}

func build(filters ...Filter) (string, error) {
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
			return "", err
		}

		builder.Write(validatedColumnName + hotcoal.Wrap(" = ?"))
		values = append(values, filter.Value)
	}

	builder.Write(";")

	return builder.String(), nil
}
