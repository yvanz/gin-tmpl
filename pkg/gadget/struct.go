package gadget

import (
	"reflect"
	"regexp"

	"gorm.io/gorm"
)

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

type (
	ColumnField struct {
		Name  string
		Field reflect.StructField
	}
	MyStructFields struct {
		strField []string
		numField []string
		other    []string
		all      []ColumnField
	}
)

func (sf *MyStructFields) Add(fields ...ColumnField) {
	for _, columnField := range fields {
		sf.all = append(sf.all, columnField)
		fieldKind := columnField.Field.Type.Kind()

		switch {
		default:
			sf.other = append(sf.other, columnField.Name)
		case fieldKind == reflect.String:
			sf.strField = append(sf.strField, columnField.Name)
		case IsNumber(fieldKind):
			sf.numField = append(sf.numField, columnField.Name)
		}
	}
}

func (sf *MyStructFields) Merge(other MyStructFields) {
	sf.Add(other.all...)
}

func (sf MyStructFields) GetStringField() []string {
	return sf.strField
}

func NewColumnField(field reflect.StructField, name string) ColumnField {
	return ColumnField{
		Name:  name,
		Field: field,
	}
}

func IsNumber(kind reflect.Kind) bool {
	numberKinds := []reflect.Kind{
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64,
	}

	for i := range numberKinds {
		if kind == numberKinds[i] {
			return true
		}
	}

	return false
}

func FieldsFromModel(m interface{}, db *gorm.DB, recurse bool) (fields MyStructFields) {
	t := reflect.TypeOf(m)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	return FieldsFromStruct(t, db, recurse)
}

// FieldsFromStruct 拿到结构体中的所有属性的字段名，对于嵌套的Struct只递归一次
func FieldsFromStruct(t reflect.Type, db *gorm.DB, recurse bool) (fields MyStructFields) {
	for i := 0; i < t.NumField(); i++ {
		var columnName string
		if t.Field(i).Tag != "" {
			tag := t.Field(i).Tag.Get("gorm")
			reg := regexp.MustCompile(`column:\s*([^;\s]+)`)

			// 判断嵌套
			if t.Field(i).Type.Kind() == reflect.Struct {
				if recurse {
					subFields := FieldsFromStruct(t.Field(i).Type, db, false)
					fields.Merge(subFields)
				}

				continue
			}

			if reg.MatchString(tag) {
				columnName = reg.FindStringSubmatch(tag)[1]
				fields.Add(NewColumnField(t.Field(i), columnName))
			}
		}

		if columnName == "" {
			// 判断嵌套
			if t.Field(i).Type.Kind() == reflect.Struct {
				if recurse {
					subFields := FieldsFromStruct(t.Field(i).Type, db, false)
					fields.Merge(subFields)
				}
				continue
			}

			columnName = db.NamingStrategy.ColumnName("", t.Field(i).Name)
			fields.Add(NewColumnField(t.Field(i), columnName))
		}
	}

	return fields
}
