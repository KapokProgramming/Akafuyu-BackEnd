package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func JSONifyRows(rows *sql.Rows) []byte {
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	masterData := make([]map[string]interface{}, 0)
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
		masterData = append(masterData, entry)
	}
	out, err := json.Marshal(masterData)
	if err != nil {
		panic(err)
	}
	return out
}

func StandardResponseWriter(w http.ResponseWriter, res StandardResponse) {
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))
}
