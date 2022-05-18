package querybuilder

import (
	"errors"
	"fmt"
	"strings"
)

type DeleteQuery struct {
	table      string
	conditions []whereClause
}

func (s *DeleteQuery) Where(query string, args ...interface{}) *DeleteQuery {
	args, _ = unifyArgs(args...)
	condition := whereClause{
		query: query,
		args:  args,
	}
	newQuery := *s

	newQuery.conditions = append(newQuery.conditions, condition)
	return &newQuery
}

func (s *DeleteQuery) Build() (string, []interface{}, error) {
	if s.table == "" {
		return "", nil, errors.New(ErrTableIsEmpty)
	}

	var query string
	var args []interface{}

	query = "DELETE FROM " + s.table

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
