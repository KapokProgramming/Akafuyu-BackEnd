package main

type PostData struct {
	Title   string `json:"title"`
	RawBody string `json:"raw_body"`
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
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Bio         string `json:"bio"`
	Timestamp   string `json:"timestamp"`
}
