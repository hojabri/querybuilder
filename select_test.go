package querybuilder

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectQuery_Build(t *testing.T) {
	ids := []int{1, 2, 3}
	tests := []struct {
		name           string
		query          *SelectQuery
		wantBuiltQuery string
		wantArgs       []interface{}
		wantErr        error
	}{
		{
			name:           "test1",
			query:          Select("table1"),
			wantBuiltQuery: "SELECT * FROM table1",
			wantArgs:       nil,
			wantErr:        nil,
		},
		{
			name:           "test2",
			query:          Select(""),
			wantBuiltQuery: "",
			wantArgs:       nil,
			wantErr:        errors.New(ErrTableIsEmpty),
		},
		{
			name:           "test3",
			query:          Select("table1").Columns("field1"),
			wantBuiltQuery: "SELECT field1 FROM table1",
			wantArgs:       nil,
			wantErr:        nil,
		},
		{
			name:           "test4",
			query:          Select("table1").Columns("field1").Columns("field2"),
			wantBuiltQuery: "SELECT field1,field2 FROM table1",
			wantArgs:       nil,
			wantErr:        nil,
		},
		{
			name:           "test5",
			query:          Select("table1").Columns("field1").Columns("field2,field3"),
			wantBuiltQuery: "SELECT field1,field2,field3 FROM table1",
			wantArgs:       nil,
			wantErr:        nil,
		},
		{
			name:           "test6",
			query:          Select("table1").Columns("field1").Where("id > ?", 120),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ?)",
			wantArgs:       []interface{}{120},
			wantErr:        nil,
		},
		{
			name:           "test7",
			query:          Select("table1").Columns("field1").Where("id > ?", 120).Order("timestamp", OrderDesc),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ?) ORDER BY timestamp DESC",
			wantArgs:       []interface{}{120},
			wantErr:        nil,
		},
		{
			name:           "test8",
			query:          Select("table1").Columns("field1").Where("id > ?", 120).Order("timestamp", OrderDesc).Order("id", OrderAsc).Group("field1,field2"),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ?) GROUP BY field1,field2 ORDER BY timestamp DESC,id ASC",
			wantArgs:       []interface{}{120},
			wantErr:        nil,
		},
		{
			name:           "test9",
			query:          Select("table1").Columns("field1").Where("id > ?", 120).Order("timestamp", OrderDesc).Group("field1,field2").Limit(1000).Offset(0),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ?) GROUP BY field1,field2 ORDER BY timestamp DESC LIMIT 1000 OFFSET 0",
			wantArgs:       []interface{}{120},
			wantErr:        nil,
		},
		{
			name:           "test10",
			query:          Select("table1").Columns("field1").Where("id > ? OR name = ?", ids, 120, "Omid"),
			wantBuiltQuery: "",
			wantArgs:       nil,
			wantErr:        errors.New(ErrWrongNumberOfArgs),
		},
		{
			name:           "test11",
			query:          Select("table1").Columns("field1").Where(In("id", []interface{}{120, 140, 160})).Where("name=?", "Omid"),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id IN (?,?,?)) AND (name=?)",
			wantArgs:       []interface{}{120, 140, 160, "Omid"},
			wantErr:        nil,
		},
		{
			name:           "test12",
			query:          Select("table1").Columns("field1").Where(In("id", ids)),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id IN (?,?,?))",
			wantArgs:       []interface{}{ids[0], ids[1], ids[2]},
			wantErr:        nil,
		},
		{
			name:           "test13",
			query:          Select("table1").Columns("field1").Where("id > ? OR name = ?", 120, "Omid"),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ? OR name = ?)",
			wantArgs:       []interface{}{120, "Omid"},
			wantErr:        nil,
		},
		{
			name:           "test14",
			query:          Select("table1").Columns("field1").Where("id > ? OR name = ?", 120, "Omid").Where(In("id", ids, 4, 5, 6)),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ? OR name = ?) AND (id IN (?,?,?,?,?,?))",
			wantArgs:       []interface{}{120, "Omid", 1, 2, 3, 4, 5, 6},
			wantErr:        nil,
		},
		{
			name:           "test15",
			query:          Select("table1").Columns("field1").Joins("table2", "table1.id=table2.t_id", JoinInner).Where(In("id", ids, 4, 5, 6)),
			wantBuiltQuery: "SELECT field1 FROM table1 JOIN table2 ON table1.id=table2.t_id WHERE (id IN (?,?,?,?,?,?))",
			wantArgs:       []interface{}{1, 2, 3, 4, 5, 6},
			wantErr:        nil,
		},
		{
			name:           "test16",
			query:          Select("table1").Columns("field1").Joins("table2", "table1.id=table2.t_id", JoinLeft).Where(In("id", ids, 4, 5, 6)),
			wantBuiltQuery: "SELECT field1 FROM table1 LEFT JOIN table2 ON table1.id=table2.t_id WHERE (id IN (?,?,?,?,?,?))",
			wantArgs:       []interface{}{1, 2, 3, 4, 5, 6},
			wantErr:        nil,
		},
		{
			name:           "test17",
			query:          Select("table1").Columns("field1").Joins("table2", "table1.id=table2.t_id", JoinRight).Where(In("id", ids, 4, 5, 6)),
			wantBuiltQuery: "SELECT field1 FROM table1 RIGHT JOIN table2 ON table1.id=table2.t_id WHERE (id IN (?,?,?,?,?,?))",
			wantArgs:       []interface{}{1, 2, 3, 4, 5, 6},
			wantErr:        nil,
		},
		{
			name:           "test18",
			query:          Select("table1").Columns("field1").Joins("table2", "table1.id=table2.t_id", 10).Where(In("id", ids, 4, 5, 6)),
			wantBuiltQuery: "SELECT field1 FROM table1 JOIN table2 ON table1.id=table2.t_id WHERE (id IN (?,?,?,?,?,?))",
			wantArgs:       []interface{}{1, 2, 3, 4, 5, 6},
			wantErr:        nil,
		},
		{
			name: "test19",
			query: Select("table1").
				Columns("SUM(field1) AS total").
				Where("id > ?", 120).
				Order("timestamp", OrderDesc).
				Group("field1,field2").
				Having("SUM(field1)>10").
				Limit(1000).Offset(0),
			wantBuiltQuery: "SELECT SUM(field1) AS total FROM table1 WHERE (id > ?) GROUP BY field1,field2 HAVING (SUM(field1)>10) ORDER BY timestamp DESC LIMIT 1000 OFFSET 0",
			wantArgs:       []interface{}{120},
			wantErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBuiltQuery, gotBuiltArgs, err := tt.query.Build()
			require.Equal(t, tt.wantBuiltQuery, gotBuiltQuery, "Build() gotBuiltQuery = %v, wantBuiltQuery %v", gotBuiltQuery, tt.wantBuiltQuery)
			require.Equal(t, tt.wantArgs, gotBuiltArgs, "Build() gotBuiltArgs = %v, wantBuiltQuery %v", gotBuiltArgs, tt.wantArgs)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestIn(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		query, args := In("id", 1, 2, 3)
		require.Equal(t, "id IN (?,?,?)", query)
		require.Equal(t, []interface{}{1, 2, 3}, args)
	})

	t.Run("test2", func(t *testing.T) {
		ids := []int{1, 2, 3}
		query, args := In("id", ids)
		require.Equal(t, "id IN (?,?,?)", query)
		require.Equal(t, []interface{}{1, 2, 3}, args)
	})

	t.Run("test3", func(t *testing.T) {
		ids := []int{1, 2, 3}
		query, args := In("id", ids, 4, 5, 6)
		require.Equal(t, "id IN (?,?,?,?,?,?)", query)
		require.Equal(t, []interface{}{1, 2, 3, 4, 5, 6}, args)
	})

	t.Run("test4", func(t *testing.T) {
		var ids []int
		query, args := In("id", ids)
		require.Equal(t, "", query)
		require.Equal(t, []interface{}(nil), args)
	})
}

func TestSelectQuery_Rebind(t *testing.T) {
	tests := []struct {
		name        string
		driver      DriverName
		selectQuery *SelectQuery
		want        string
	}{
		{
			name:        "test Postgres",
			driver:      DriverPostgres,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=$1)",
		},
		{
			name:        "test PGX",
			driver:      DriverPGX,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=$1)",
		},
		{
			name:        "test pq-timeouts",
			driver:      DriverPqTimeout,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=$1)",
		},
		{
			name:        "test CloudSqlPostgres",
			driver:      DriverCloudSqlPostgres,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=$1)",
		},
		{
			name:        "test MySQL",
			driver:      DriverMySQL,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=?)",
		},
		{
			name:        "test Sqlite3",
			driver:      DriverSqlite3,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=?)",
		},
		{
			name:        "test oci8",
			driver:      DriverOCI8,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=:arg1)",
		},
		{
			name:        "test ora",
			driver:      DriverORA,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=:arg1)",
		},
		{
			name:        "test goracle",
			driver:      DriverGORACLE,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=:arg1)",
		},
		{
			name:        "test SqlServer",
			driver:      DriverSqlServer,
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=@p1)",
		},
		{
			name:        "test unknown",
			driver:      "abcdefg",
			selectQuery: Select("table1").Where("id=?", 100),
			want:        "SELECT * FROM table1 WHERE (id=?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, _, err := tt.selectQuery.Build()
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
