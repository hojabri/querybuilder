package querybuilder

import (
	"errors"
	"testing"
	
	"github.com/stretchr/testify/require"
)

func TestSelectQuery_Build(t *testing.T) {
	s := &SelectQuery{}
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
			query:          s.Table("table1"),
			wantBuiltQuery: "SELECT * FROM table1",
			wantArgs:       nil,
			wantErr:        nil,
		},
		{
			name:           "test2",
			query:          s,
			wantBuiltQuery: "",
			wantArgs:       nil,
			wantErr:        errors.New(ErrTableIsEmpty),
		},
		{
			name:           "test3",
			query:          s.Table("table1").Select("field1"),
			wantBuiltQuery: "SELECT field1 FROM table1",
			wantArgs:       nil,
			wantErr:        nil,
		},
		{
			name:           "test4",
			query:          s.Table("table1").Select("field1").Select("field2"),
			wantBuiltQuery: "SELECT field1,field2 FROM table1",
			wantArgs:       nil,
			wantErr:        nil,
		},
		{
			name:           "test5",
			query:          s.Table("table1").Select("field1").Select("field2,field3"),
			wantBuiltQuery: "SELECT field1,field2,field3 FROM table1",
			wantArgs:       nil,
			wantErr:        nil,
		},
		{
			name:           "test6",
			query:          s.Table("table1").Select("field1").Where("id > ?", 120),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ?)",
			wantArgs:       []interface{}{120},
			wantErr:        nil,
		},
		{
			name:           "test7",
			query:          s.Table("table1").Select("field1").Where("id > ?", 120).Order("timestamp"),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ?) ORDER BY timestamp",
			wantArgs:       []interface{}{120},
			wantErr:        nil,
		},
		{
			name:           "test8",
			query:          s.Table("table1").Select("field1").Where("id > ?", 120).Order("timestamp").Group("field1,field2"),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ?) GROUP BY field1,field2 ORDER BY timestamp",
			wantArgs:       []interface{}{120},
			wantErr:        nil,
		},
		{
			name:           "test9",
			query:          s.Table("table1").Select("field1").Where("id > ?", 120).Order("timestamp").Group("field1,field2").Limit(1000).Offset(0),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ?) GROUP BY field1,field2 ORDER BY timestamp LIMIT 1000 OFFSET 0",
			wantArgs:       []interface{}{120},
			wantErr:        nil,
		},
		{
			name:           "test10",
			query:          s.Table("table1").Select("field1").Where("id > ? OR name = ?", ids, 120, "Omid"),
			wantBuiltQuery: "",
			wantArgs:       nil,
			wantErr:        errors.New(ErrWrongNumberOfArgs),
		},
		{
			name:           "test11",
			query:          s.Table("table1").Select("field1").Where(In("id", []interface{}{120, 140, 160})).Where("name=?", "Omid"),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id IN (?,?,?)) AND (name=?)",
			wantArgs:       []interface{}{120, 140, 160, "Omid"},
			wantErr:        nil,
		},
		{
			name:           "test12",
			query:          s.Table("table1").Select("field1").Where(In("id", ids)),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id IN (?,?,?))",
			wantArgs:       []interface{}{ids[0], ids[1], ids[2]},
			wantErr:        nil,
		},
		{
			name:           "test13",
			query:          s.Table("table1").Select("field1").Where("id > ? OR name = ?", 120, "Omid"),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ? OR name = ?)",
			wantArgs:       []interface{}{120, "Omid"},
			wantErr:        nil,
		},
		{
			name:           "test14",
			query:          s.Table("table1").Select("field1").Where("id > ? OR name = ?", 120, "Omid").Where(In("id", ids, 4, 5, 6)),
			wantBuiltQuery: "SELECT field1 FROM table1 WHERE (id > ? OR name = ?) AND (id IN (?,?,?,?,?,?))",
			wantArgs:       []interface{}{120, "Omid", 1, 2, 3, 4, 5, 6},
			wantErr:        nil,
		},
		{
			name:           "test15",
			query:          s.Table("table1").Select("field1").Joins("table2", "table1.id=table2.t_id", JoinInner).Where(In("id", ids, 4, 5, 6)),
			wantBuiltQuery: "SELECT field1 FROM table1 JOIN table2 ON table1.id=table2.t_id WHERE (id IN (?,?,?,?,?,?))",
			wantArgs:       []interface{}{1, 2, 3, 4, 5, 6},
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
