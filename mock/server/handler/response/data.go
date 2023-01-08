package response

import (
	"encoding/json"
	"net/http"
)

type Data struct {
	Message string `json:"message"`
	Result  string `json:"result"`
}

func (r *Data) Encode(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(r)
}
