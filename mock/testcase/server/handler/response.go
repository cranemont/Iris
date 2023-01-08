package handler

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Data    string `json:"data"`
	Error   string `json:"error"`
}

func (r *Response) Encode(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(r)
}
