package querybuilder

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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
	
	u := Update()
	tests := []struct {
		name      string
		query     *UpdateQuery
		wantQuery string
		wantArgs  []any
		wantErr   error
	}{
		{
			name:      "test1",
			query:     u.Table("table1").MapValues(nil),
			wantQuery: "",
			wantArgs:  []any(nil),
			wantErr:   errors.New(ErrColumnValueMapIsEmpty),
		},
		{
			name:      "test2",
			query:     u,
			wantQuery: "",
			wantArgs:  []any(nil),
			wantErr:   errors.New(ErrTableIsEmpty),
		},
		{
			name:      "test3",
			query:     u.Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			wantQuery: "UPDATE table1 SET field1=?,field2=? WHERE (id=?)",
			wantArgs:  []any{10, "test", 5000},
			wantErr:   nil,
		},
		{
			name:      "test4 - non pointer struct",
			query:     u.Table("table1").StructValues(sampleStruct).Where("id=?", 5000),
			wantQuery: "UPDATE table1 SET name=?,email=?,ID=? WHERE (id=?)",
			wantArgs:  []any{"Omid", "o.hojabri@gmail.com", 74639876, 5000},
			wantErr:   nil,
		},
		{
			name:      "test5 -with pointer struct",
			query:     u.Table("table1").StructValues(&sampleStruct).Where("id=?", 5000),
			wantQuery: "UPDATE table1 SET name=?,email=?,ID=? WHERE (id=?)",
			wantArgs:  []any{"Omid", "o.hojabri@gmail.com", 74639876, 5000},
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

func TestUpdatePanicNotStruct(t *testing.T) {
	require.Panics(t, func() {
		u := Update()
		u.Table("table1").StructValues(123)
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
			updateQuery: UpdateByDriver(DriverPostgres).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=$1,field2=$2 WHERE (id=$3)",
		},
		{
			name:        "test PGX",
			updateQuery: UpdateByDriver(DriverPGX).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=$1,field2=$2 WHERE (id=$3)",
		},
		{
			name:        "test pq-timeouts",
			updateQuery: UpdateByDriver(DriverPqTimeout).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=$1,field2=$2 WHERE (id=$3)",
		},
		{
			name:        "test CloudSqlPostgres",
			updateQuery: UpdateByDriver(DriverCloudSqlPostgres).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=$1,field2=$2 WHERE (id=$3)",
		},
		{
			name:        "test MySQL",
			updateQuery: UpdateByDriver(DriverMySQL).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=?,field2=? WHERE (id=?)",
		},
		{
			name:        "test Sqlite3",
			updateQuery: UpdateByDriver(DriverSqlite3).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=?,field2=? WHERE (id=?)",
		},
		{
			name:        "test oci8",
			updateQuery: UpdateByDriver(DriverOCI8).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=:arg1,field2=:arg2 WHERE (id=:arg3)",
		},
		{
			name:        "test ora",
			updateQuery: UpdateByDriver(DriverORA).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=:arg1,field2=:arg2 WHERE (id=:arg3)",
		},
		{
			name:        "test goracle",
			updateQuery: UpdateByDriver(DriverGORACLE).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=:arg1,field2=:arg2 WHERE (id=:arg3)",
		},
		{
			name:        "test SqlServer",
			updateQuery: UpdateByDriver(DriverSqlServer).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
			want:        "UPDATE table1 SET field1=@p1,field2=@p2 WHERE (id=@p3)",
		},
		{
			name:        "test unknown",
			updateQuery: UpdateByDriver("abcdefg").Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}).Where("id=?", 5000),
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