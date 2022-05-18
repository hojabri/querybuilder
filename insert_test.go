package querybuilder

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInsertQuery_Build(t *testing.T) {
	type NestedStructType struct {
		Level string `json:"level,omitempty" db:"level"`
	}
	type SampleStructType struct {
		Name  string      `json:"name,omitempty" db:"name"`
		Email string      `json:"email,omitempty" db:"email"`
		ID    interface{} `json:"id,omitempty"`
		Order float32     `json:"order" db:"-"`
		Image *[]byte     `json:"image" db:"image"`
		Grade int         `json:"grade" db:"grade"`
		NestedStructType
	}

	sampleImage := []byte("sample image bytes")

	tests := []struct {
		name      string
		query     *InsertQuery
		wantQuery string
		wantArgs  []interface{}
		wantErr   error
	}{
		{
			name:      "test1",
			query:     Insert("table1").MapValues(nil),
			wantQuery: "",
			wantArgs:  []interface{}(nil),
			wantErr:   errors.New(ErrColumnValueMapIsEmpty),
		},
		{
			name:      "test2",
			query:     Insert(""),
			wantQuery: "",
			wantArgs:  []interface{}(nil),
			wantErr:   errors.New(ErrTableIsEmpty),
		},
		{
			name:      "test3",
			query:     Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			wantQuery: "INSERT INTO table1(field1,field2) VALUES(?,?)",
			wantArgs:  []interface{}{10, "test"},
			wantErr:   nil,
		},
		{
			name: "test4 - non pointer struct",
			query: Insert("table1").StructValues(SampleStructType{
				Name:             "Omid",
				Email:            "o.hojabri@gmail.com",
				ID:               74639876,
				Image:            &sampleImage,
				Grade:            2,
				NestedStructType: NestedStructType{Level: "123"},
			}),
			wantQuery: "INSERT INTO table1(name,email,ID,image,grade,level) VALUES(?,?,?,?,?,?)",
			wantArgs:  []interface{}{"Omid", "o.hojabri@gmail.com", 74639876, sampleImage, 2, "123"},
			wantErr:   nil,
		},
		{
			name: "test5 -with pointer struct",
			query: Insert("table1").StructValues(&SampleStructType{
				Name:  "Omid",
				Email: "o.hojabri@gmail.com",
				ID:    74639876,
				Image: &sampleImage,
				Grade: 2,
			}),
			wantQuery: "INSERT INTO table1(name,email,ID,image,grade,level) VALUES(?,?,?,?,?,?)",
			wantArgs:  []interface{}{"Omid", "o.hojabri@gmail.com", 74639876, sampleImage, 2, ""},
			wantErr:   nil,
		},
		{
			name: "test6 - empty and non pointer field",
			query: Insert("table1").StructValues(SampleStructType{
				Name:  "Omid",
				Email: "o.hojabri@gmail.com",
				ID:    74639876,
				Image: &sampleImage,
			}),
			wantQuery: "INSERT INTO table1(name,email,ID,image,grade,level) VALUES(?,?,?,?,?,?)",
			wantArgs:  []interface{}{"Omid", "o.hojabri@gmail.com", 74639876, sampleImage, 0, ""},
			wantErr:   nil,
		},
		{
			name: "test7 - empty and pointer field",
			query: Insert("table1").StructValues(SampleStructType{
				Name:  "Omid",
				Email: "o.hojabri@gmail.com",
				ID:    74639876,
			}),
			wantQuery: "INSERT INTO table1(name,email,ID,grade,level) VALUES(?,?,?,?,?)",
			wantArgs:  []interface{}{"Omid", "o.hojabri@gmail.com", 74639876, 0, ""},
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

func TestInsertPanicNotStruct(t *testing.T) {
	require.Panics(t, func() {
		Insert("table1").StructValues(123)
	}, "should panic with non struct types")
}

func TestInsertQuery_Rebind(t *testing.T) {
	tests := []struct {
		name        string
		driver      DriverName
		insertQuery *InsertQuery
		want        string
	}{
		{
			name:        "test Postgres",
			driver:      DriverPostgres,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES($1,$2)",
		},
		{
			name:        "test PGX",
			driver:      DriverPGX,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES($1,$2)",
		},
		{
			name:        "test pq-timeouts",
			driver:      DriverPqTimeout,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES($1,$2)",
		},
		{
			name:        "test CloudSqlPostgres",
			driver:      DriverCloudSqlPostgres,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES($1,$2)",
		},
		{
			name:        "test MySQL",
			driver:      DriverMySQL,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(?,?)",
		},
		{
			name:        "test Sqlite3",
			driver:      DriverSqlite3,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(?,?)",
		},
		{
			name:        "test oci8",
			driver:      DriverOCI8,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(:arg1,:arg2)",
		},
		{
			name:        "test ora",
			driver:      DriverORA,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(:arg1,:arg2)",
		},
		{
			name:        "test goracle",
			driver:      DriverGORACLE,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(:arg1,:arg2)",
		},
		{
			name:        "test SqlServer",
			driver:      DriverSqlServer,
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(@p1,@p2)",
		},
		{
			name:        "test unknown",
			driver:      "abcdefg",
			insertQuery: Insert("table1").MapValues(map[string]interface{}{"field1": 10, "field2": "test"}),
			want:        "INSERT INTO table1(field1,field2) VALUES(?,?)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, _, err := tt.insertQuery.Build()
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
