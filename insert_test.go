package querybuilder

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInsertQuery_Build(t *testing.T) {
	i := Insert()

	tests := []struct {
		name      string
		query     *InsertQuery
		wantQuery string
		wantArgs  []any
		wantErr   error
	}{
		{
			name:      "test1",
			query:     i.Table("table1").MapValues(map[string]any{"field1": 10, "field2": "test"}),
			wantQuery: "INSERT INTO table1(field1,field2) VALUES(?,?)",
			wantArgs:  []any{10, "test"},
			wantErr:   nil,
		},
		{
			name:      "test2",
			query:     i.Table("table1").MapValues(nil),
			wantQuery: "",
			wantArgs:  []any(nil),
			wantErr:   errors.New(ErrColumnValueMapIsEmpty),
		},
		{
			name:      "test3",
			query:     i,
			wantQuery: "",
			wantArgs:  []any(nil),
			wantErr:   errors.New(ErrTableIsEmpty),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotArgs, err := tt.query.Build()
			require.Equal(t, tt.wantQuery, gotQuery, "Build() gotQuery = %v, wantQuery %v", gotQuery, tt.wantQuery)
			require.Equal(t, tt.wantArgs, gotArgs, "Build() gotArgs = %v, wantQuery %v", gotArgs, tt.wantArgs)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
