package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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
		hash_pw, err := GetHashedPassword(register_data.Password)
		if err != nil {
			panic(err)
		}
		db := createConnectionToDatabase()
		var count int
		err = db.QueryRow("SELECT count(*) FROM users WHERE username=?", register_data.Username).Scan(&count)
		if err != nil {
			panic(err)
		}
		if count == 0 {
			query := "INSERT INTO users (username, password, email) VALUES (?, ?, ?);"
			db.Exec(query, register_data.Username, hash_pw, register_data.Email)
			res.Status = "success"
		} else {
			res.Status = "error"
			res.Data = "Username already exists"
		}
	}
	StandardResponseWriter(w, res)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var res StandardResponse
	var user_id int
	var err error
	var user UserData
	if len(vars["id"]) > 0 {
		user_id, err = strconv.Atoi(vars["id"])
		if err != nil {
			res.Status = "error"
			res.Data = err.Error()
		}
	} else {
		reqToken := r.Header.Get("Authorization")
		fmt.Printf("reqToken: %v", reqToken)
		if len(reqToken) > 0 {
			splitToken := strings.Split(reqToken, "Bearer ")
			reqToken = splitToken[1]
			user_id, err = ValidateJWT(reqToken)
			if err != nil {
				res.Status = "fail"
				res.Data = "Invalid token"
				StandardResponseWriter(w, res)
				return
			}
		} else {
			res.Status = "fail"
			res.Data = "Not Logged in"
			StandardResponseWriter(w, res)
			return
		}
	}
	query := "SELECT * FROM users WHERE user_id=?;"
	db.QueryRow(query, user_id).Scan(&user.UserID, &user.Username, &user.DisplayName, &user.Password, &user.Email, &user.Bio, &user.Timestamp)
	res.Status = "success"
	res.Data = user
	StandardResponseWriter(w, res)
}

func TokenTestHandler(w http.ResponseWriter, r *http.Request) {
	var res StandardResponse
	w.Header().Set("Content-Type", "application/json")
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	user_id, err := ValidateJWT(reqToken)
	if err != nil {
		res.Status = "error"
		res.Data = err.Error()
		StandardResponseWriter(w, res)
		return
	}
	db := createConnectionToDatabase()
	query := "SELECT * FROM users WHERE user_id=?;"
	var user UserData
	// err = db.QueryRow(query, user_id).Scan(&user)
	err = db.QueryRow(query, user_id).Scan(&user.UserID, &user.Username, &user.DisplayName, &user.Password, &user.Email, &user.Bio, &user.Timestamp)
	if err != nil {
		res.Status = "error"
		res.Data = err.Error()
		StandardResponseWriter(w, res)
		return
	}
	res.Status = "success"
	res.Data = user
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
		query := "SELECT user_id,password FROM users WHERE username=?;"
		var user_id int
		var hash_pw string
		err = db.QueryRow(query, login_data.Username).Scan(&user_id, &hash_pw)
		switch {
		case err == sql.ErrNoRows:
			res.Status = "error"
			res.Data = "Wrong Login Data"
		case err != nil:
			panic(err)
		default:
			err := ValidatePassword(login_data.Password, hash_pw)
			if err != nil {
				res.Status = "fail"
				res.Data = err
				break
			}
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
		var post PostWithAuthorData
		query := "SELECT posts.post_title, posts.post_body, IFNULL(users.display_name, users.username), COUNT(DISTINCT users_star.us_id) FROM posts INNER JOIN users ON posts.author=users.user_id AND post_id=? LEFT JOIN users_star ON posts.post_id=users_star.post_id GROUP BY posts.post_id;"
		err := db.QueryRow(query, vars["id"]).Scan(&post.PostTitle, &post.PostBody, &post.Author, &post.StarCount)
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
			page_size_num = 12
		}
		var user_id int
		reqToken := r.Header.Get("Authorization")
		fmt.Printf("reqToken: %v", reqToken)
		if len(reqToken) > 0 {
			splitToken := strings.Split(reqToken, "Bearer ")
			reqToken = splitToken[1]
			user_id, err = ValidateJWT(reqToken)
			if err != nil {
				user_id = -1
			}
		}
		var query string
		var rows *sql.Rows
		if user_id == -1 {
			query = "SELECT * FROM posts LIMIT ?,?;"
			rows, err = db.Query(query, page_num*page_size_num, page_size_num)
		} else {
			query = "SELECT posts.* FROM `posts` INNER JOIN users_follow ON users_follow.follower_id=? AND posts.author=users_follow.followed_id AND posts.isFollowerOnlyPost=1 UNION SELECT * FROM `posts` WHERE posts.isFollowerOnlyPost=0 ORDER BY post_id LIMIT ?,?;"
			rows, err = db.Query(query, user_id, page_num*page_size_num, page_size_num)
		}
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
		query := "INSERT INTO posts (post_title, post_body, author, isFollowerOnlyPost) VALUES (?, ?, ?, ?);"
		db.Exec(query, post_data.PostTitle, post_data.PostBody, post_data.Author, post_data.FollowerOnly)
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
