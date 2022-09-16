package main

type PostData struct {
	Title   string `json:"title"`
	RawBody string `json:"raw_body"`
}

type StandardResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
