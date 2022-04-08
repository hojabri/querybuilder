package querybuilder

import (
	"errors"
	"log"
	"strings"
)

type InsertQuery struct {
	driver              DriverName
	table               string
	indexedColumnValues IndexedColumnValues
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
	newQuery.indexedColumnValues = mapToIndexColumnValue(columnValues)
	return &newQuery
}

// StructValues gets and struct and extract column/values,
func (s *InsertQuery) StructValues(structure any) *InsertQuery {
	newQuery := *s
	m, err := structToMap(structure)
	if err != nil {
		log.Panic(err)
	}
	newQuery.indexedColumnValues = m
	return &newQuery
}

func (s *InsertQuery) Build() (string, []any, error) {
	if s.table == "" {
		return "", nil, errors.New(ErrTableIsEmpty)
	}
	if len(s.indexedColumnValues) == 0 {
		return "", nil, errors.New(ErrColumnValueMapIsEmpty)
	}
	var query string
	
	args := make([]any, len(s.indexedColumnValues))
	
	// make column slice
	columns := make([]string, len(s.indexedColumnValues))
	
	for i := 0; i < len(s.indexedColumnValues); i++ {
		indexedColumnValue := s.indexedColumnValues[i]
		columns[i] = indexedColumnValue.Key
		args[i] = indexedColumnValue.Value
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
