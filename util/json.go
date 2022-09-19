package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"server/model"
)

func RowsToMap(rows *sql.Rows) []map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	rows_map := make([]map[string]interface{}, 0)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		rows_map = append(rows_map, entry)
	}
	return rows_map
}

func StandardResponseWriter(w http.ResponseWriter, res model.StandardResponse) {
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))
}

func JSONRowsToString(rows_map []map[string]interface{}) []byte {
	out, err := json.Marshal(rows_map)
	if err != nil {
		panic(err)
	}
	return out
}
