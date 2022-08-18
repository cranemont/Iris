package event

// register를 공통으로?
type EventHandler interface {
	RegisterFn()
	Call(funcName string, args interface{})
}
