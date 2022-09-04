package handler

// 여기서 statecode를?
type Result struct {
}

type Handler interface {
	Handle(data interface{}) error
}
