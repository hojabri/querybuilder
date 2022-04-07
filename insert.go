package querybuilder

import (
	"errors"
	"sort"
	"strings"
)

type InsertQuery struct {
	driver       DriverName
	table        string
	columnValues map[string]any
}

func (s *InsertQuery) Table(name string) *InsertQuery {
	newQuery := *s
	newQuery.table = name
	return &newQuery
}

// MapValues gets columns and values,
// Enter Column/Values as a key/value map
func (s *InsertQuery) MapValues(columnValues map[string]any) *InsertQuery {
	newQuery := *s
	newQuery.columnValues = columnValues
	return &newQuery
}

func (s *InsertQuery) Build() (string, []any, error) {
	if s.table == "" {
		return "", nil, errors.New(ErrTableIsEmpty)
	}
	if len(s.columnValues) == 0 {
		return "", nil, errors.New(ErrColumnValueMapIsEmpty)
	}
	var query string
	var args []any

	// make column slice
	var columns []string
	for column, _ := range s.columnValues {
		columns = append(columns, column)
	}

	// sort columns
	sort.Strings(columns)

	// adding values to args
	for _, column := range columns {
		args = append(args, s.columnValues[column])

	}
	//
	// add table name
	query = "INSERT INTO " + s.table + "(" + strings.Join(columns, ",") + ") VALUES(" + strings.TrimSuffix(strings.Repeat("?,", len(columns)), ",") + ")"

	return query, args, nil
}
