package models

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})
var timePtrType = reflect.TypeOf((*time.Time)(nil))
var flexTimeType = reflect.TypeOf(FlexTime{})

// CVSchema generates a JSON schema example from the CV struct.
func CVSchema() string {
	b, _ := json.MarshalIndent(buildExample(reflect.TypeOf(CV{})), "", "  ")
	return string(b)
}

// JobSchema generates a JSON schema example from the Job struct.
func JobSchema() string {
	b, _ := json.MarshalIndent(buildExample(reflect.TypeOf(Job{})), "", "  ")
	return string(b)
}

// CVVariantSchema generates a JSON schema example from the CVVariant struct.
func CVVariantSchema() string {
	b, _ := json.MarshalIndent(buildExample(reflect.TypeOf(CVVariant{})), "", "  ")
	return string(b)
}

func buildExample(t reflect.Type) interface{} {
	if t.Kind() == reflect.Ptr {
		if t == timePtrType {
			return "YYYY-MM-DDT00:00:00Z or null"
		}
		t = t.Elem()
	}

	if t == timeType {
		return "YYYY-MM-DDT00:00:00Z"
	}

	if t == flexTimeType {
		return "YYYY-MM-DDT00:00:00Z or null"
	}

	switch t.Kind() {
	case reflect.String:
		return "string"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return 0

	case reflect.Bool:
		return false

	case reflect.Slice:
		elem := t.Elem()
		if elem.Kind() == reflect.String {
			return []string{"string"}
		}
		return []interface{}{buildExample(elem)}

	case reflect.Struct:
		m := make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() || f.Anonymous {
				continue // skip unexported and embedded (e.g. gorm.Model)
			}
			jsonTag := f.Tag.Get("json")
			if jsonTag == "-" || strings.HasPrefix(jsonTag, "-,") {
				continue
			}
			name := strings.SplitN(jsonTag, ",", 2)[0]
			if name == "" {
				name = f.Name
			}
			m[name] = buildExample(f.Type)
		}
		return m
	}

	return nil
}
