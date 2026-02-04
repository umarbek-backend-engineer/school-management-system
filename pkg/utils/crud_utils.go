package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func GenerateInsertQuery(model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println(dbTag)
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		fmt.Println(dbTag)
		if dbTag != "" && dbTag != "id" {
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			placeholders += "?"
		}
	}
	return fmt.Sprintf("insert into teachers(%s) values (%s)", columns, placeholders)
}

func IsValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}
func IsValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func Getfilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for key, column := range params {
		val := r.URL.Query().Get(key)
		if val != "" {
			query += " AND " + column + " = ?"
			args = append(args, val)
		}
	}

	return query, args
}
