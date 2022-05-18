package querybuilder

import (
	"errors"
	"fmt"
	"strings"
)

type JoinType int

const (
	JoinInner = iota
	JoinLeft
	JoinRight
)

func joinTypeString(joinType JoinType) string {
	switch joinType {
	case JoinInner:
		return "JOIN"
	case JoinLeft:
		return "LEFT JOIN"
	case JoinRight:
		return "RIGHT JOIN"
	default:
		return "JOIN"
	}
}

type OrderDirection int

const (
	OrderAsc = iota
	OrderDesc
)

func orderDirectionString(direction OrderDirection) string {
	switch direction {
	case OrderDesc:
		return "DESC"
	default:
		return "ASC"
	}
}

type SelectQuery struct {
	columns    []columnClause
	table      string
	joins      []joinClause
	conditions []whereClause
	havings    []havingClause
	groupBy    []groupByClause
	orderBy    []orderByClause
	limit      interface{}
	offset     interface{}
}

func (s *SelectQuery) Columns(query string, args ...interface{}) *SelectQuery {
	args, _ = unifyArgs(args...)
	column := columnClause{
		query: query,
		args:  args,
	}
	newQuery := *s
	newQuery.columns = append(newQuery.columns, column)
	return &newQuery
}

func (s *SelectQuery) Joins(tableName string, on string, joinType JoinType, args ...interface{}) *SelectQuery {
	args, _ = unifyArgs(args...)
	join := joinClause{
		tableName: tableName,
		args:      args,
		on:        on,
		joinType:  joinType,
	}
	newQuery := *s
	newQuery.joins = append(newQuery.joins, join)
	return &newQuery
}

func (s *SelectQuery) Where(query string, args ...interface{}) *SelectQuery {
	args, _ = unifyArgs(args...)
	condition := whereClause{
		query: query,
		args:  args,
	}
	newQuery := *s

	newQuery.conditions = append(newQuery.conditions, condition)
	return &newQuery
}

func (s *SelectQuery) Having(query string, args ...interface{}) *SelectQuery {
	args, _ = unifyArgs(args...)
	having := havingClause{
		query: query,
		args:  args,
	}
	newQuery := *s

	newQuery.havings = append(newQuery.havings, having)
	return &newQuery
}

func (s *SelectQuery) Group(query string) *SelectQuery {
	clause := groupByClause{
		fields: query,
	}
	newQuery := *s
	newQuery.groupBy = append(newQuery.groupBy, clause)
	return &newQuery
}

func (s *SelectQuery) Order(column string, direction OrderDirection) *SelectQuery {
	clause := orderByClause{
		field:     column,
		direction: direction,
	}
	newQuery := *s
	newQuery.orderBy = append(newQuery.orderBy, clause)
	return &newQuery
}

func (s *SelectQuery) Limit(limit int64) *SelectQuery {
	newQuery := *s
	newQuery.limit = limit
	return &newQuery
}

func (s *SelectQuery) Offset(offset int64) *SelectQuery {
	newQuery := *s
	newQuery.offset = offset
	return &newQuery
}

func (s *SelectQuery) Build() (string, []interface{}, error) {
	if s.table == "" {
		return "", nil, errors.New(ErrTableIsEmpty)
	}
	var args []interface{}
	var columns string
	//
	// check for columns
	if len(s.columns) > 0 {
		var columnsSlice []string
		for _, column := range s.columns {
			columnsSlice = append(columnsSlice, column.query)
			args = append(args, column.args...)
		}
		columns = strings.Join(columnsSlice, ",")
	} else {
		columns = "*"
	}
	// add columns
	query := "SELECT " + columns
	//
	// add table name
	query = query + " FROM " + s.table
	//
	// add joins
	if len(s.joins) > 0 {
		for _, join := range s.joins {
			args = append(args, join.args...)
			query = query + " " + joinTypeString(join.joinType) + " " + join.tableName + " ON " + join.on
		}
	}
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
	//
	// add group by
	if len(s.groupBy) > 0 {
		var groupBySlice []string
		for _, groupBy := range s.groupBy {
			groupBySlice = append(groupBySlice, groupBy.fields)
		}
		query = query + " GROUP BY " + strings.Join(groupBySlice, ",")
	}
	//
	// add having
	if len(s.havings) > 0 {
		var havingSlice []string
		for _, having := range s.havings {
			havingSlice = append(havingSlice, fmt.Sprintf("(%s)", having.query))
			args = append(args, having.args...)
		}
		query = query + " HAVING " + strings.Join(havingSlice, " AND ")
	}
	//
	// add order by
	if len(s.orderBy) > 0 {
		var orderBySlice []string
		for _, orderBy := range s.orderBy {
			orderBySlice = append(orderBySlice, orderBy.field+" "+orderDirectionString(orderBy.direction))
		}
		query = query + " ORDER BY " + strings.Join(orderBySlice, ",")
	}
	//
	// add limit
	if s.limit != nil {
		if limit, ok := s.limit.(int64); ok {
			query = query + fmt.Sprintf(" LIMIT %d", limit)
		} else {
			return "", nil, errors.New(ErrLimitNotInteger)
		}
	}
	//
	// add offset
	if s.offset != nil {
		if offset, ok := s.offset.(int64); ok {
			query = query + fmt.Sprintf(" OFFSET %d", offset)
		} else {
			return "", nil, errors.New(ErrOffsetNotInteger)
		}
	}

	// compare the number of args and ? in tableName
	if len(args) != strings.Count(query, "?") {
		return "", nil, errors.New(ErrWrongNumberOfArgs)
	}
	//
	// return built tableName and args
	return query, args, nil
}
