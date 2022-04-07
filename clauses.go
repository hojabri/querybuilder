package querybuilder

type columnClause struct {
	query string
	args  []any
}

type whereClause struct {
	query string
	args  []any
}

type havingClause struct {
	query string
	args  []any
}

type joinClause struct {
	tableName string
	on        string
	joinType  JoinType
	args      []any
}

type groupByClause struct {
	fields string
}

type orderByClause struct {
	field     string
	direction OrderDirection
}
