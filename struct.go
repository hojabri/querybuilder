package querybuilder

import (
	"errors"
	"reflect"
	"sort"
	"strings"
)

type KeyValue struct {
	Key   string
	Value any
}

type IndexedColumnValues map[int]KeyValue

func mapToIndexColumnValue(columnValues map[string]any) IndexedColumnValues {
	indexedColumnValues := make(IndexedColumnValues, len(columnValues))
	columns := make([]string, len(columnValues))
	i := 0
	for column := range columnValues {
		columns[i] = column
		i++
	}
	sort.Strings(columns)
	
	for i = 0; i < len(columns); i++ {
		indexedColumnValues[i] = KeyValue{
			Key:   columns[i],
			Value: columnValues[columns[i]],
		}
	}
	
	return indexedColumnValues
}

func structToMap(s any) (IndexedColumnValues, error) {
	columnValues := make(IndexedColumnValues)
	v := reflect.ValueOf(s)
	// if its a pointer, resolve its value
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	
	if v.Kind() != reflect.Struct {
		return nil, errors.New("unexpected type")
	}
	e := v.Type()
	for i := 0; i < e.NumField(); i++ {
		name := e.Field(i).Name
		tag := strings.Split(e.Field(i).Tag.Get("db"), ",")[0] // use split to ignore tag "options"
		
		// ignore columns with -
		if tag == "-" {
			continue
		}
		value := v.FieldByIndex(e.Field(i).Index)
		column := tag
		if tag == "" {
			column = name
		}
		columnValues[i] = KeyValue{Key: column, Value: value.Interface()}
	}
	return columnValues, nil
}
