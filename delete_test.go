package querybuilder

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDeleteQuery_Build(t *testing.T) {
	tests := []struct {
		name      string
		query     *DeleteQuery
		wantQuery string
		wantArgs  []interface{}
		wantErr   error
	}{
		{
			name:      "test1",
			query:     Delete(""),
			wantQuery: "",
			wantArgs:  []interface{}(nil),
			wantErr:   errors.New(ErrTableIsEmpty),
		},
		{
			name:      "test2",
			query:     Delete("table1").Where("id=?", 5000),
			wantQuery: "DELETE FROM table1 WHERE (id=?)",
			wantArgs:  []interface{}{5000},
			wantErr:   nil,
		},
		{
			name:      "test3 - wrong number of arguments",
			query:     Delete("table1").Where("id=?", 5000, 10),
			wantQuery: "",
			wantArgs:  []interface{}(nil),
			wantErr:   errors.New(ErrWrongNumberOfArgs),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t1 := time.Now()
			gotQuery, gotArgs, err := tt.query.Build()
			t.Logf("duration: %s", time.Since(t1))
			require.Equal(t, tt.wantQuery, gotQuery, "Build() gotQuery = %v, wantQuery %v", gotQuery, tt.wantQuery)
			require.Equal(t, tt.wantArgs, gotArgs, "Build() gotArgs = %v, wantQuery %v", gotArgs, tt.wantArgs)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestDeleteQuery_Rebind(t *testing.T) {
	tests := []struct {
		name        string
		driver      DriverName
		deleteQuery *DeleteQuery
		want        string
	}{
		{
			name:        "test Postgres",
			driver:      DriverPostgres,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=$1)",
		},
		{
			name:        "test PGX",
			driver:      DriverPGX,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=$1)",
		},
		{
			name:        "test pq-timeouts",
			driver:      DriverPqTimeout,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=$1)",
		},
		{
			name:        "test CloudSqlPostgres",
			driver:      DriverCloudSqlPostgres,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=$1)",
		},
		{
			name:        "test MySQL",
			driver:      DriverMySQL,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=?)",
		},
		{
			name:        "test Sqlite3",
			driver:      DriverSqlite3,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=?)",
		},
		{
			name:        "test oci8",
			driver:      DriverOCI8,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=:arg1)",
		},
		{
			name:        "test ora",
			driver:      DriverORA,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=:arg1)",
		},
		{
			name:        "test goracle",
			driver:      DriverGORACLE,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=:arg1)",
		},
		{
			name:        "test SqlServer",
			driver:      DriverSqlServer,
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=@p1)",
		},
		{
			name:        "test unknown",
			driver:      "abcdefg",
			deleteQuery: Delete("table1").Where("id=?", 5000),
			want:        "DELETE FROM table1 WHERE (id=?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, _, err := tt.deleteQuery.Build()
			if err != nil {
				t.Errorf("can't build the query: %s", err)
			}
			Driver = tt.driver
			if got := Rebind(query); got != tt.want {
				require.Equal(t, tt.want, got, "Rebind() got = %v, want %v", got, tt.want)
			}
		})
	}
}
