package handler

type Handler interface {
	Handle(data interface{}) error
}
