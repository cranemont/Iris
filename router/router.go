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
	Route(handle string, data interface{}) []byte
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

func (r *router) Route(mode string, data interface{}) []byte {
	result := handler.DefaultResult()
	switch mode {
	case Judge:
		req := handler.JudgeRequest{}
		err := json.Unmarshal(data.([]byte), &req)
		if err != nil {
			r.logger.Error(fmt.Sprintf("judge: invalid request data: %s, %s", string(data.([]byte)), err))
			break
		}

		res, err := r.judgeHandler.Handle(req)
		if err != nil {
			r.logger.Error(fmt.Sprintf("judge: failed to handle request: %s", err))
			break
		}
		result = r.judgeHandler.ResultToJson(res)
	case SpecialJudge:
		// special-judge handler
	case CustomTestcase:
		// custom-testcase handler
	default:
		r.logger.Error("unregistered handler")
	}

	return result
}
