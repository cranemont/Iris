package event

// register를 공통으로?
type EventHandler interface {
	Call(funcName string, args interface{})
}
