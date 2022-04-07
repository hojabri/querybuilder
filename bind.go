package querybuilder

import (
	"strconv"
	"strings"
)

// The Rebind function is inspired by sqlx: https://pkg.go.dev/github.com/jmoiron/sqlx#Rebind

// Bindvar types supported by Rebind
const (
	UNKNOWN = iota
	QUESTION
	DOLLAR
	NAMED
	AT
)

const (
	DriverPostgres         = "postgres"
	DriverPGX              = "pgx"
	DriverPqTimeout        = "pq-timeouts"
	DriverCloudSqlPostgres = "cloudsqlpostgres"
	DriverMySQL            = "mysql"
	DriverSqlite3          = "sqlite3"
	DriverOCI8             = "oci8"
	DriverORA              = "ora"
	DriverGORACLE          = "goracle"
	DriverSqlServer        = "sqlserver"
)

type DriverName string

// BindType returns the bindtype for a given a drivername/database.
func BindType(driverName DriverName) int {
	switch driverName {
	case "postgres", "pgx", "pq-timeouts", "cloudsqlpostgres":
		return DOLLAR
	case "mysql":
		return QUESTION
	case "sqlite3":
		return QUESTION
	case "oci8", "ora", "goracle":
		return NAMED
	case "sqlserver":
		return AT
	}
	return UNKNOWN
}

// rebind a query table the default bindtype (QUESTION) to the target bindtype.
func rebind(bindType int, query string) string {
	switch bindType {
	case QUESTION, UNKNOWN:
		return query
	}

	// Add space enough for 10 params before we have to allocate
	rqb := make([]byte, 0, len(query)+10)

	var i, j int

	for i = strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		rqb = append(rqb, query[:i]...)

		switch bindType {
		case DOLLAR:
			rqb = append(rqb, '$')
		case NAMED:
			rqb = append(rqb, ':', 'a', 'r', 'g')
		case AT:
			rqb = append(rqb, '@', 'p')
		}

		j++
		rqb = strconv.AppendInt(rqb, int64(j), 10)

		query = query[i+1:]
	}

	return string(append(rqb, query...))
}
