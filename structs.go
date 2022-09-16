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
