package main

import "database/sql"

type PostData struct {
	Author    int    `json:"author"`
	PostTitle string `json:"post_title"`
	PostBody  string `json:"post_body"`
}

type PostWithAuthorData struct {
	PostTitle string `json:"post_title"`
	PostBody  string `json:"post_body"`
	Author    string `json:"author"`
	StarCount string `json:"star_count"`
}

type RegisterData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type StandardResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type UserData struct {
	UserID      int            `json:"user_id"`
	Username    string         `json:"username"`
	DisplayName sql.NullString `json:"display_name"`
	Password    string         `json:"password"`
	Email       string         `json:"email"`
	Bio         sql.NullString `json:"bio"`
	Timestamp   string         `json:"timestamp"`
}
