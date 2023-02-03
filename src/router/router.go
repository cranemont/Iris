package router

import (
	"encoding/json"
	"fmt"

	"github.com/cranemont/judge-manager/src/handler"
	"github.com/cranemont/judge-manager/src/service/logger"
)

type To string

const (
	Judge        To = "Judge"
	SpecialJudge To = "SpecialJudge"
	Run          To = "Run"
	Interactive  To = "Interactive"
)

type Router interface {
	Route(path To, id string, data []byte) []byte
}

type router struct {
	judgeHandler *handler.JudgeHandler
	logger       logger.Logger
}

func NewRouter(
	judgeHandler *handler.JudgeHandler,
	logger logger.Logger,
) *router {
	return &router{judgeHandler, logger}
}

func (r *router) Route(path To, id string, data []byte) []byte {
	var handlerResult json.RawMessage
	var err error

	switch path {
	case Judge:
		handlerResult, err = r.judgeHandler.Handle(id, data)
	case SpecialJudge:
		// special-judge handler
	case Run:
		// custom-testcase handler
	default:
		err = fmt.Errorf("handler does not exist: %s", path)
	}

	if err != nil {
		if u, ok := err.(*handler.HandlerError); ok {
			r.logger.Log(u.Level(), err.Error())
		} else {
			r.logger.Log(logger.ERROR, fmt.Sprintf("router: %s", err.Error()))
		}
	}
	return NewResponse(id, handlerResult, err).Marshal()
}
