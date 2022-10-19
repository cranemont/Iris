package router

import (
	"encoding/json"
	"fmt"

	"github.com/cranemont/judge-manager/handler"
	"github.com/cranemont/judge-manager/service/logger"
)

const (
	Judge          = "Judge"
	SpecialJudge   = "SpecialJudge"
	CustomTestcase = "CustomTestcase"
)

type Router interface {
	Route(mode string, id string, data []byte) []byte
}

type router struct {
	judgeHandler *handler.JudgeHandler
	logger       *logger.Logger
}

func NewRouter(
	judgeHandler *handler.JudgeHandler,
	logger *logger.Logger,
) *router {
	return &router{judgeHandler, logger}
}

func (r *router) Route(mode string, id string, data []byte) []byte {
	var handlerResult json.RawMessage
	var err error

	switch mode {
	case Judge:
		handlerResult, err = r.judgeHandler.Handle(id, data)
	case SpecialJudge:
		// special-judge handler
	case CustomTestcase:
		// custom-testcase handler
	default:
		err = fmt.Errorf("invalid mode: %s", mode)
	}

	if err != nil {
		if u, ok := err.(*handler.HandlerError); ok {
			r.logger.Log(u.Level(), err.Error())
		} else {
			r.logger.Error(fmt.Sprintf("router: %s", err.Error()))
		}
	}
	return NewResponse(id, handlerResult, err).Marshal()
}
