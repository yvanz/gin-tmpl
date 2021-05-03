package gadget

import "reflect"

func GetTableColumn(obj interface{}) []string {
	t := reflect.TypeOf(obj).Elem()
	var column = make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		columnName := t.Field(i).Tag.Get("json")
		if columnName != "-" {
			column = append(column, columnName)
		}
	}
	return column
}

func GetTableColumnByTag(obj interface{}, tag string) []string {
	t := reflect.TypeOf(obj).Elem()
	var column = make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		columnName := t.Field(i).Tag.Get(tag)
		column = append(column, columnName)
	}
	return column
}

func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

func StructToMapByJSONTag(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		key := obj1.Field(i).Tag.Get("json")
		if key != "-" {
			data[key] = obj2.Field(i).Interface()
		}
	}
	return data
}
