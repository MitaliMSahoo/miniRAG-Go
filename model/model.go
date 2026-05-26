package model

type Document struct {
	Text string `json:"text"`
}

type AddRequest struct {
	Document []Document `json:"documents"`
}

type AddResponse struct {
	Added   int    `json:"added"`
	Message string `json:"message"`
}

type QueryRequest struct {
	Content string `json:"content"`
}

type QueryResponse struct {
	Answer string `json:"answer"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
