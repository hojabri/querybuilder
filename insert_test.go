package querybuilder

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInsertQuery_Build(t *testing.T) {
	i := Insert()

	sampleStruct := struct {
		Name  string      `json:"name,omitempty" db:"name"`
		Email string      `json:"email,omitempty" db:"email"`
		ID    interface{} `json:"id,omitempty" db:"id"`
		Order float32     `json:"order" db:"-"`
	}{
		Name:  "Omid",
		Email: "o.hojabri@gmail.com",
		ID:    74639876,
	}

	tests := []struct {
		name      string
		query     *InsertQuery
		wantQuery string
		wantArgs  []any
		wantErr   error
	}{
		{
			name:      "test1",
			query:     i.Table("table1").MapValues(nil),
			wantQuery: "",
			wantArgs:  []any(nil),
			wantErr:   errors.New(ErrColumnValueMapIsEmpty),
		},

		{
			name:      "test2",
			query:     i,
			wantQuery: "",
			wantArgs:  []any(nil),
			wantErr:   errors.New(ErrTableIsEmpty),
		},
		{
			name:      "test3",
			query:     i.Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			wantQuery: "INSERT INTO table1(field1,field2) VALUES(?,?)",
			wantArgs:  []any{10, "test"},
			wantErr:   nil,
		},
		{
			name:      "test4 - non pointer struct",
			query:     i.Table("table1").StructValues(sampleStruct),
			wantQuery: "INSERT INTO table1(email,id,name) VALUES(?,?,?)",
			wantArgs:  []any{"o.hojabri@gmail.com", 74639876, "Omid"},
			wantErr:   nil,
		},
		{
			name:      "test5 -with pointer struct",
			query:     i.Table("table1").StructValues(&sampleStruct),
			wantQuery: "INSERT INTO table1(email,id,name) VALUES(?,?,?)",
			wantArgs:  []any{"o.hojabri@gmail.com", 74639876, "Omid"},
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

func TestInsertQuery_Rebind(t *testing.T) {
	tests := []struct {
		name        string
		insertQuery *InsertQuery
		want        string
	}{
		{
			name:        "test Postgres",
			insertQuery: InsertByDriver(DriverPostgres).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES($1,$2)",
		},
		{
			name:        "test PGX",
			insertQuery: InsertByDriver(DriverPGX).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES($1,$2)",
		},
		{
			name:        "test pq-timeouts",
			insertQuery: InsertByDriver(DriverPqTimeout).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES($1,$2)",
		},
		{
			name:        "test CloudSqlPostgres",
			insertQuery: InsertByDriver(DriverCloudSqlPostgres).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES($1,$2)",
		},
		{
			name:        "test MySQL",
			insertQuery: InsertByDriver(DriverMySQL).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(?,?)",
		},
		{
			name:        "test Sqlite3",
			insertQuery: InsertByDriver(DriverSqlite3).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(?,?)",
		},
		{
			name:        "test oci8",
			insertQuery: InsertByDriver(DriverOCI8).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(:arg1,:arg2)",
		},
		{
			name:        "test ora",
			insertQuery: InsertByDriver(DriverORA).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(:arg1,:arg2)",
		},
		{
			name:        "test goracle",
			insertQuery: InsertByDriver(DriverGORACLE).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(:arg1,:arg2)",
		},
		{
			name:        "test SqlServer",
			insertQuery: InsertByDriver(DriverSqlServer).Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(@p1,@p2)",
		},
		{
			name:        "test unknown",
			insertQuery: InsertByDriver("abcdefg").Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(?,?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, _, err := tt.insertQuery.Build()
			if err != nil {
				t.Errorf("can't build the query: %s", err)
			}
			if got := tt.insertQuery.Rebind(query); got != tt.want {
				require.Equal(t, tt.want, got, "Rebind() got = %v, want %v", got, tt.want)
			}
		})
	}
}
