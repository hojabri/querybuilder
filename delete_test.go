package querybuilder

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDeleteQuery_Build(t *testing.T) {
	d := Delete()
	tests := []struct {
		name      string
		query     *DeleteQuery
		wantQuery string
		wantArgs  []any
		wantErr   error
	}{
		{
			name:      "test1",
			query:     d,
			wantQuery: "",
			wantArgs:  []any(nil),
			wantErr:   errors.New(ErrTableIsEmpty),
		},
		{
			name:      "test2",
			query:     d.Table("table1").Where("id=?", 5000),
			wantQuery: "DELETE table1 WHERE (id=?)",
			wantArgs:  []any{5000},
			wantErr:   nil,
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
		deleteQuery *DeleteQuery
		want        string
	}{
		{
			name:        "test Postgres",
			deleteQuery: DeleteByDriver(DriverPostgres).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=$1)",
		},
		{
			name:        "test PGX",
			deleteQuery: DeleteByDriver(DriverPGX).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=$1)",
		},
		{
			name:        "test pq-timeouts",
			deleteQuery: DeleteByDriver(DriverPqTimeout).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=$1)",
		},
		{
			name:        "test CloudSqlPostgres",
			deleteQuery: DeleteByDriver(DriverCloudSqlPostgres).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=$1)",
		},
		{
			name:        "test MySQL",
			deleteQuery: DeleteByDriver(DriverMySQL).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=?)",
		},
		{
			name:        "test Sqlite3",
			deleteQuery: DeleteByDriver(DriverSqlite3).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=?)",
		},
		{
			name:        "test oci8",
			deleteQuery: DeleteByDriver(DriverOCI8).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=:arg1)",
		},
		{
			name:        "test ora",
			deleteQuery: DeleteByDriver(DriverORA).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=:arg1)",
		},
		{
			name:        "test goracle",
			deleteQuery: DeleteByDriver(DriverGORACLE).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=:arg1)",
		},
		{
			name:        "test SqlServer",
			deleteQuery: DeleteByDriver(DriverSqlServer).Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=@p1)",
		},
		{
			name:        "test unknown",
			deleteQuery: DeleteByDriver("abcdefg").Table("table1").Where("id=?", 5000),
			want:        "DELETE table1 WHERE (id=?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, _, err := tt.deleteQuery.Build()
			if err != nil {
				t.Errorf("can't build the query: %s", err)
			}
			if got := tt.deleteQuery.Rebind(query); got != tt.want {
				require.Equal(t, tt.want, got, "Rebind() got = %v, want %v", got, tt.want)
			}
		})
	}
}
