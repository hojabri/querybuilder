package querybuilder

import (
	"fmt"
	"reflect"
	"strings"
)

func unifyArgs(args ...interface{}) ([]interface{}, int) {
	count := 0
	var newArgs []interface{}
	for _, arg := range args {
		s := reflect.ValueOf(arg)
		switch reflect.TypeOf(arg).Kind() {
		case reflect.Slice:
			count = +s.Len()
			for i := 0; i < s.Len(); i++ {
				newArgs = append(newArgs, s.Index(i).Interface())
			}
		default:
			count++
			newArgs = append(newArgs, s.Interface())
		}
	}
	return newArgs, count
}

func In(column string, args ...interface{}) (string, []interface{}) {
	args, count := unifyArgs(args...)
	
	if count == 0 {
		return "", nil
	}
	values := strings.TrimSuffix(strings.Repeat("?,", count), ",")
	query := fmt.Sprintf("%s IN (%s)", column, values)
	return query, args
}
