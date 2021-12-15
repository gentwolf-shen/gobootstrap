package util

import (
	"math"
	"reflect"
	"strings"
)

func ToMap(value interface{}) map[string]interface{} {
	if value == nil {
		return nil
	}

	values := reflect.ValueOf(value)
	if values.Kind() == reflect.Map {
		return value.(map[string]interface{})
	}

	types := reflect.TypeOf(value)
	if types.Kind() == reflect.Ptr {
		types = types.Elem()
		values = values.Elem()
	}

	size := types.NumField()
	var data = make(map[string]interface{})
	for i := 0; i < size; i++ {
		name := types.Field(i).Name
		name = strings.ToLower(name[0:1]) + name[1:]
		if v := values.Field(i).Interface(); v != nil {
			data[name] = v
		}
	}
	return data
}

func QueryDbTagField(value interface{}) string {
	types := reflect.TypeOf(value)
	if types.Kind() == reflect.Ptr {
		types = types.Elem()
	}

	if types.Kind() == reflect.Slice {
		types = types.Elem()
	}

	size := types.NumField()
	arr := make([]string, size)

	index := 0
	for i := 0; i < size; i++ {
		name, bl := types.Field(i).Tag.Lookup("db")
		if !bl {
			continue
		}
		arr[index] = name
		index++
	}
	return strings.Join(arr[0:index], ",")
}

func QueryDbTagMap(value interface{}, target string) map[string]interface{} {
	if value == nil {
		return nil
	}

	values := reflect.ValueOf(value)
	if values.Kind() == reflect.Map {
		return value.(map[string]interface{})
	}

	types := reflect.TypeOf(value)

	if types.Kind() == reflect.Ptr {
		types = types.Elem()
		values = values.Elem()
	}

	size := values.NumField()
	var data = make(map[string]interface{}, size)
	for i := 0; i < size; i++ {
		name, bl := types.Field(i).Tag.Lookup("db")
		if !bl {
			continue
		}

		if strings.Contains(name, ",") {
			if !strings.Contains(name, target) {
				continue
			} else {
				name = strings.Split(name, ",")[0]
			}
		}

		if v := values.Field(i).Interface(); v != nil {
			data[name] = v
		}
	}
	return data
}

func ToArray(value map[string]interface{}) ([]string, []interface{}) {
	size := len(value)
	keys := make([]string, size)
	values := make([]interface{}, size)

	index := 0
	for k, v := range value {
		keys[index] = k
		values[index] = v
		index++
	}
	return keys, values
}

func Ceil(size, count int64) int64 {
	return int64(math.Ceil(float64(count) / float64(size)))
}

func ToOffset(page, size int64) int64 {
	return (page - 1) * size
}
