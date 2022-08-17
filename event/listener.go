package event

type EventListener interface {
	On(eventCh <-chan interface{}, handlerFn string)
}
