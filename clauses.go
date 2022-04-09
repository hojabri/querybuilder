package querybuilder

type columnClause struct {
	query string
	args  []interface{}
}

type whereClause struct {
	query string
	args  []interface{}
}

type havingClause struct {
	query string
	args  []interface{}
}

type joinClause struct {
	tableName string
	on        string
	joinType  JoinType
	args      []interface{}
}

type groupByClause struct {
	fields string
}

type orderByClause struct {
	field     string
	direction OrderDirection
}
