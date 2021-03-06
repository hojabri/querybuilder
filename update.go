package querybuilder

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type UpdateQuery struct {
	table               string
	indexedColumnValues IndexedColumnValues
	conditions          []whereClause
}

func (s *UpdateQuery) Where(query string, args ...interface{}) *UpdateQuery {
	args, _ = unifyArgs(args...)
	condition := whereClause{
		query: query,
		args:  args,
	}
	newQuery := *s

	newQuery.conditions = append(newQuery.conditions, condition)
	return &newQuery
}

// MapValues gets columns and values,
// Enter Column/Values as a key/value map
func (s *UpdateQuery) MapValues(columnValues map[string]interface{}) *UpdateQuery {
	newQuery := *s
	newQuery.indexedColumnValues = mapToIndexColumnValue(columnValues)
	return &newQuery
}

// StructValues gets and struct and extract column/values,
func (s *UpdateQuery) StructValues(structure interface{}) *UpdateQuery {
	newQuery := *s
	m, err := structToMap(structure)
	if err != nil {
		log.Panic(err)
	}
	newQuery.indexedColumnValues = m
	return &newQuery
}

func (s *UpdateQuery) Build() (string, []interface{}, error) {
	if s.table == "" {
		return "", nil, errors.New(ErrTableIsEmpty)
	}
	if len(s.indexedColumnValues) == 0 {
		return "", nil, errors.New(ErrColumnValueMapIsEmpty)
	}
	var query string
	args := make([]interface{}, len(s.indexedColumnValues))

	// make column slice
	columns := make([]string, len(s.indexedColumnValues))
	var setQuery []string

	for i := 0; i < len(s.indexedColumnValues); i++ {
		indexedColumnValue := s.indexedColumnValues[i]
		columns[i] = indexedColumnValue.Key
		args[i] = indexedColumnValue.Value
		setQuery = append(setQuery, columns[i]+"=?")
	}

	query = "UPDATE " + s.table + " SET " + strings.Join(setQuery, ",")

	//
	// check for where part
	if len(s.conditions) > 0 {
		var conditionsSlice []string
		for _, condition := range s.conditions {
			conditionsSlice = append(conditionsSlice, fmt.Sprintf("(%s)", condition.query))
			args = append(args, condition.args...)
		}
		query = query + " WHERE " + strings.Join(conditionsSlice, " AND ")
	}

	// compare the number of args and ? in tableName
	if len(args) != strings.Count(query, "?") {
		return "", nil, errors.New(ErrWrongNumberOfArgs)
	}

	return query, args, nil
}
