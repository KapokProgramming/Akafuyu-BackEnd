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
	r.HandleFunc("/test", TestHandler)
	r.HandleFunc("/json", JSONHandler)

	r.NotFoundHandler = http.HandlerFunc(emptyJsonHandler)
	fmt.Println("Listening on :7700")
	log.Fatal(http.ListenAndServe(":7700", handlers.LoggingHandler(os.Stdout, r)))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `ola`)
}

func JSONHandler(w http.ResponseWriter, r *http.Request) {
	body := make([]map[string]interface{}, 0)
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", body)
	fmt.Fprintf(w, "body: %+v", body)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DB")))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	query := "SELECT NOW();"
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
