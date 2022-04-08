package querybuilder

import (
	"errors"
	"fmt"
	"strings"
)

type DeleteQuery struct {
	driver     DriverName
	table      string
	conditions []whereClause
}

func (s *DeleteQuery) Table(name string) *DeleteQuery {
	newQuery := *s
	newQuery.table = name
	return &newQuery
}

func (s *DeleteQuery) Where(query string, args ...any) *DeleteQuery {
	args, _ = unifyArgs(args...)
	condition := whereClause{
		query: query,
		args:  args,
	}
	newQuery := *s
	
	newQuery.conditions = append(newQuery.conditions, condition)
	return &newQuery
}

func (s *DeleteQuery) Build() (string, []any, error) {
	if s.table == "" {
		return "", nil, errors.New(ErrTableIsEmpty)
	}
	
	var query string
	var args []any
	
	query = "DELETE " + s.table
	
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
	
	return query, args, nil
}

// Rebind transforms a query table QUESTION to the DB driver's bindvar type.
func (s *DeleteQuery) Rebind(query string) string {
	return rebind(BindType(s.driver), query)
}
