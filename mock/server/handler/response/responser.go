package response

import (
	"net/http"
)

type Responser interface {
	Ok(w http.ResponseWriter, data Data, code int)
	Error(w http.ResponseWriter, error string, code int)
}

type responser struct {
}

func NewResponser() *responser {
	return &responser{}
}

func (r *responser) Ok(w http.ResponseWriter, data Data, code int) {
	w.WriteHeader(code)
	data.Encode(w)
}

func (r *responser) Error(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	d := Data{Message: error}
	d.Encode(w)
}
