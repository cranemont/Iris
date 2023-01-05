package handler

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func (r *Response) Encode(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(r)
}
