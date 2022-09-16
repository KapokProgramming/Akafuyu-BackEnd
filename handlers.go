package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res StandardResponse
	res.Status = "success"
	res.Data = "test"
	StandardResponseWriter(w, res)
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := createConnectionToDatabase()
	var res StandardResponse
	res.Status = "error"
	switch r.Method {
	case "GET":
		page := r.URL.Query().Get("page")
		page_num, err := strconv.Atoi(page)
		if err != nil {
			page_num = 0
		}
		query := fmt.Sprintf("SELECT * FROM posts LIMIT %d,10;", page_num*10)
		rows, err := db.Query(query)
		if err != nil {
			panic(err)
		}
		res.Status = "success"
		res.Data = RowsToMap(rows)
		break
	case "POST":
		var id int
		var post_data PostData
		res.Status = "success"
		err := json.NewDecoder(r.Body).Decode(&post_data)
		if err != nil {
			res.Status = "error"
			panic(err)
		}
		tx, err := db.Begin()
		if err != nil {
			res.Status = "error"
			panic(err)
		}
		query, err := tx.Prepare(
			`INSERT INTO posts (
				title,
				raw_body,
			) VALUES ($1, $2)
			RETURNING id`)
		if err != nil {
			res.Status = "error"
			panic(err)
		}
		defer query.Close()
		err = query.QueryRow(
			post_data.Title,
			post_data.RawBody,
		).Scan(&id)
		if err != nil {
			res.Status = "error"
			panic(err)
		}
		err = tx.Commit()
		if err != nil {
			res.Status = "error"
			panic(err)
		}
		break
	}
	StandardResponseWriter(w, res)
}

func JSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res StandardResponse
	res.Status = "success"
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var body map[string]interface{}
	err = json.Unmarshal(reqbody, &body)
	if err != nil {
		panic(err)
	}
	res.Data = body
	fmt.Printf("%+v", body)
	StandardResponseWriter(w, res)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res StandardResponse
	res.Status = "success"
	query := "SELECT NOW();"
	db := createConnectionToDatabase()
	rows, err := db.Query(query)
	if err != nil {
		res.Status = "error"
	} else {
		res.Data = RowsToMap(rows)
	}
	StandardResponseWriter(w, res)
}

func emptyJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{}`)
}
