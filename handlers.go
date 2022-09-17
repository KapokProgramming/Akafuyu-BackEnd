package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res StandardResponse
	res.Status = "success"
	res.Data = "running"
	StandardResponseWriter(w, res)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res StandardResponse
	res.Status = "fail"
	switch r.Method {
	case "POST":
		var register_data RegisterData
		err := json.NewDecoder(r.Body).Decode(&register_data)
		if err != nil {
			res.Status = "error"
			panic(err)
		}
		db := createConnectionToDatabase()
		query := "INSERT INTO users (username, password, email) VALUES (?, ?, ?);"
		db.Exec(query, register_data.Username, register_data.Password, register_data.Email)
		res.Status = "success"
	}
	StandardResponseWriter(w, res)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res StandardResponse
	res.Status = "fail"
	switch r.Method {
	case "POST":
		var login_data LoginData
		err := json.NewDecoder(r.Body).Decode(&login_data)
		if err != nil {
			panic(err)
		}
		db := createConnectionToDatabase()
		query := "SELECT user_id FROM users WHERE username=? AND password=?;"
		var user_id int
		err = db.QueryRow(query, login_data.Username, login_data.Password).Scan(&user_id)
		switch {
		case err == sql.ErrNoRows:
			res.Status = "error"
			res.Data = "Wrong Login Data"
		case err != nil:
			panic(err)
		default:
			res.Status = "success"
			res.Data, err = CreateJWT(user_id)
			if err != nil {
				panic(err)
			}
		}
	}
	StandardResponseWriter(w, res)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	db := createConnectionToDatabase()
	var res StandardResponse
	res.Status = "error"
	switch r.Method {
	case "GET":
		var post PostData
		query := fmt.Sprintf("SELECT title, raw_body FROM posts where posts.post_id=?;")
		err := db.QueryRow(query, vars["id"]).Scan(&post.Title, &post.RawBody)
		switch {
		case err == sql.ErrNoRows:
			res.Status = "fail"
			res.Data = fmt.Sprintf("Invalid Post ID: %s", vars["id"])
		case err != nil:
			panic(err)
		default:
			res.Status = "success"
			res.Data = post
		}
	}
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
		page_size := r.URL.Query().Get("n")
		page_size_num, err := strconv.Atoi(page_size)
		if err != nil {
			page_size_num = 9
		}
		query := fmt.Sprintf("SELECT * FROM posts LIMIT ?,?;")
		rows, err := db.Query(query, page_num*page_size_num, page_size_num)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		res.Status = "success"
		res.Data = RowsToMap(rows)
	case "POST":
		var post_data PostData
		err := json.NewDecoder(r.Body).Decode(&post_data)
		if err != nil {
			panic(err)
		}
		query := "INSERT INTO posts (title, raw_body) VALUES (?, ?);"
		db.Exec(query, post_data.Title, post_data.RawBody)
		res.Status = "success"
	}
	StandardResponseWriter(w, res)
}

func JSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		return
	}
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
	if r.Method != "GET" {
		return
	}
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

func EmptyJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{}`)
}
