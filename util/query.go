// string
package util

import (
	"database/sql"
)

type Query struct {
	Limit  int
	Offset int
}

func FilterStruct(arr []string, a interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	dat := a.(map[string]interface{})

	for _, s := range arr {
		if _, ok := dat[s]; ok {
			result[s] = dat[s]
		}
	}
	return result
}

func ParseRows(rows *sql.Rows) (result []map[string]interface{}, err error) {
	columns, err := rows.Columns()
	if err != nil {
		return result, err
	}
	size := len(columns)
	pts := make([]interface{}, size)
	container := make([]interface{}, size)
	for i := range pts {
		pts[i] = &container[i]
	}
	for rows.Next() {
		err = rows.Scan(pts...)
		if err != nil {
			return result, err
		}
		r := make(map[string]interface{}, size)
		for i, column := range columns {
			r[column] = container[i]
		}
		result = append(result, r)
	}
	return result, err
}

func Unit8ToString(v interface{}) string {
	r := ""
	switch v.(type) {
	case []uint8:
		arr := v.([]uint8)
		r = string(arr)
	case nil:
		r = ""
	}
	return r
}
