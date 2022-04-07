package querybuilder

import "testing"

func Test_rebind(t *testing.T) {
	type args struct {
		bindType int
		query    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				bindType: DOLLAR,
				query:    "SELECT * FROM table1 WHERE id=?",
			},
			want: "SELECT * FROM table1 WHERE id=$1",
		},
		{
			name: "test2",
			args: args{
				bindType: DOLLAR,
				query:    "SELECT * FROM table1 WHERE id=? and name=?",
			},
			want: "SELECT * FROM table1 WHERE id=$1 and name=$2",
		},
		{
			name: "test3",
			args: args{
				bindType: QUESTION,
				query:    "SELECT * FROM table1 WHERE id=?",
			},
			want: "SELECT * FROM table1 WHERE id=?",
		},
		{
			name: "test4",
			args: args{
				bindType: QUESTION,
				query:    "SELECT * FROM table1 WHERE id=? and name=?",
			},
			want: "SELECT * FROM table1 WHERE id=? and name=?",
		},
		{
			name: "test5",
			args: args{
				bindType: NAMED,
				query:    "SELECT * FROM table1 WHERE id=?",
			},
			want: "SELECT * FROM table1 WHERE id=:arg1",
		},
		{
			name: "test6",
			args: args{
				bindType: NAMED,
				query:    "SELECT * FROM table1 WHERE id=? and name=?",
			},
			want: "SELECT * FROM table1 WHERE id=:arg1 and name=:arg2",
		},
		{
			name: "test7",
			args: args{
				bindType: AT,
				query:    "SELECT * FROM table1 WHERE id=?",
			},
			want: "SELECT * FROM table1 WHERE id=@p1",
		},
		{
			name: "test8",
			args: args{
				bindType: AT,
				query:    "SELECT * FROM table1 WHERE id=? and name=?",
			},
			want: "SELECT * FROM table1 WHERE id=@p1 and name=@p2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rebind(tt.args.bindType, tt.args.query); got != tt.want {
				t.Errorf("rebind() = %v, want %v", got, tt.want)
			}
		})
	}
}
