package querybuilder

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpdateQuery_Build(t *testing.T) {
	sampleStruct := struct {
		Name  string      `json:"name,omitempty" db:"name"`
		Email string      `json:"email,omitempty" db:"email"`
		ID    interface{} `json:"id,omitempty"`
		Order float32     `json:"order" db:"-"`
	}{
		Name:  "Omid",
		Email: "o.hojabri@gmail.com",
		ID:    74639876,
	}

	tests := []struct {
		name      string
		query     *UpdateQuery
		wantQuery string
		wantArgs  []interface{}
		wantErr   error
	}{
		{
			name:      "test1",
			query:     Update("table1").MapValues(nil),
			wantQuery: "",
			wantArgs:  []interface{}(nil),
			wantErr:   errors.New(ErrColumnValueMapIsEmpty),
		},
		{
			name:      "test2",
			query:     Update(""),
			wantQuery: "",
			wantArgs:  []interface{}(nil),
			wantErr:   errors.New(ErrTableIsEmpty),
		},
		{
			name:      "test3",
			query:     Update("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			wantQuery: "UPDATE table1 SET field1=?,field2=? WHERE (id=?)",
			wantArgs:  []interface{}{10, "test", 5000},
			wantErr:   nil,
		},
		{
			name:      "test4 - non pointer struct",
			query:     Update("table1").StructValues(sampleStruct).Where("id=?", 5000),
			wantQuery: "UPDATE table1 SET name=?,email=?,ID=? WHERE (id=?)",
			wantArgs:  []interface{}{"Omid", "o.hojabri@gmail.com", 74639876, 5000},
			wantErr:   nil,
		},
		{
			name:      "test5 -with pointer struct",
			query:     Update("table1").StructValues(&sampleStruct).Where("id=?", 5000),
			wantQuery: "UPDATE table1 SET name=?,email=?,ID=? WHERE (id=?)",
			wantArgs:  []interface{}{"Omid", "o.hojabri@gmail.com", 74639876, 5000},
			wantErr:   nil,
		},
		{
			name:      "test6 -wrong number of arguments",
			query:     Update("table1").StructValues(&sampleStruct).Where("id=?"),
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

func TestUpdatePanicNotStruct(t *testing.T) {
	require.Panics(t, func() {
		Update("table1").StructValues(123)
	}, "should panic with non struct types")
}

func TestUpdateQuery_Rebind(t *testing.T) {
	tests := []struct {
		name        string
		updateQuery *UpdateQuery
		want        string
	}{
		{
			name:        "test Postgres",
			updateQuery: UpdateByDriver("table1", DriverPostgres).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=$1,field2=$2 WHERE (id=$3)",
		},
		{
			name:        "test PGX",
			updateQuery: UpdateByDriver("table1", DriverPGX).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=$1,field2=$2 WHERE (id=$3)",
		},
		{
			name:        "test pq-timeouts",
			updateQuery: UpdateByDriver("table1", DriverPqTimeout).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=$1,field2=$2 WHERE (id=$3)",
		},
		{
			name:        "test CloudSqlPostgres",
			updateQuery: UpdateByDriver("table1", DriverCloudSqlPostgres).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=$1,field2=$2 WHERE (id=$3)",
		},
		{
			name:        "test MySQL",
			updateQuery: UpdateByDriver("table1", DriverMySQL).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=?,field2=? WHERE (id=?)",
		},
		{
			name:        "test Sqlite3",
			updateQuery: UpdateByDriver("table1", DriverSqlite3).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=?,field2=? WHERE (id=?)",
		},
		{
			name:        "test oci8",
			updateQuery: UpdateByDriver("table1", DriverOCI8).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=:arg1,field2=:arg2 WHERE (id=:arg3)",
		},
		{
			name:        "test ora",
			updateQuery: UpdateByDriver("table1", DriverORA).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=:arg1,field2=:arg2 WHERE (id=:arg3)",
		},
		{
			name:        "test goracle",
			updateQuery: UpdateByDriver("table1", DriverGORACLE).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=:arg1,field2=:arg2 WHERE (id=:arg3)",
		},
		{
			name:        "test SqlServer",
			updateQuery: UpdateByDriver("table1", DriverSqlServer).MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=@p1,field2=@p2 WHERE (id=@p3)",
		},
		{
			name:        "test unknown",
			updateQuery: UpdateByDriver("table1", "abcdefg").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=?,field2=? WHERE (id=?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, _, err := tt.updateQuery.Build()
			if err != nil {
				t.Errorf("can't build the query: %s", err)
			}
			if got := tt.updateQuery.Rebind(query); got != tt.want {
				require.Equal(t, tt.want, got, "Rebind() got = %v, want %v", got, tt.want)
			}
		})
	}
}
