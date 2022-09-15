package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/employees", EmployeesHandler)

	r.NotFoundHandler = http.HandlerFunc(emptyJsonHandler)
	fmt.Println("Listening on :7700")
	log.Fatal(http.ListenAndServe(":7700", handlers.LoggingHandler(os.Stdout, r)))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `ola`)
}

func EmployeesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", "root:@/employees")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	query := "SELECT * FROM employees"
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(JSONifyRows(rows)))
}

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

func emptyJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{}`)
}
