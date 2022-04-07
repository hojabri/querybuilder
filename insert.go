package querybuilder

import (
	"errors"
	"log"
	"reflect"
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

// StructValues gets and struct and extract column/values,
func (s *InsertQuery) StructValues(structure any) *InsertQuery {
	newQuery := *s
	m, err := convertStructToMap(structure)
	if err != nil {
		log.Panic(err)
	}
	newQuery.columnValues = m
	return &newQuery
}

func convertStructToMap(s any) (map[string]any, error) {
	columnValues := make(map[string]any)
	v := reflect.ValueOf(s)
	// if its a pointer, resolve its value
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("unexpected type")
	}
	e := v.Type()
	for i := 0; i < e.NumField(); i++ {
		name := e.Field(i).Name
		tag := strings.Split(e.Field(i).Tag.Get("db"), ",")[0] // use split to ignore tag "options"

		// ignore columns with -
		if tag == "-" {
			continue
		}
		value := v.FieldByIndex(e.Field(i).Index)
		column := tag
		if tag == "" {
			column = name
		}
		columnValues[column] = value.Interface()
	}
	return columnValues, nil
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
	for column := range s.columnValues {
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

// Rebind transforms a query table QUESTION to the DB driver's bindvar type.
func (s *InsertQuery) Rebind(query string) string {
	return rebind(BindType(s.driver), query)
}
