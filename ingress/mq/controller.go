package mq

import (
	"encoding/json"
	"fmt"

	"github.com/cranemont/judge-manager/handler"
	"github.com/cranemont/judge-manager/logger"
)

const (
	Judge          = "Judge"
	SpecialJudge   = "SpecialJudge"
	CustomTestcase = "CustomTestcase"
)

type RmqController interface {
	Call(handle string, data interface{}) []byte
}

type rmqController struct {
	judgeHandler *handler.JudgeHandler
	logging      *logger.Logger
}

func NewRmqController(
	judgeHandler *handler.JudgeHandler,
	logging *logger.Logger,
) *rmqController {
	return &rmqController{judgeHandler, logging}
}

// 요청을 받고 최종 response를 내보내는 책임
func (r *rmqController) Call(handle string, data interface{}) []byte {
	result := handler.DefaultResult()
	switch handle {
	case Judge:
		req := handler.JudgeRequest{}
		err := json.Unmarshal(data.([]byte), &req)
		if err != nil {
			r.logging.Error(fmt.Sprintf("judge: invalid request data: %s", err))
			break
		}

		res, err := r.judgeHandler.Handle(req)
		if err != nil {
			r.logging.Error(fmt.Sprintf("judge: failed to handle request: %s", err))
			break
		}
		result = r.judgeHandler.ResultToJson(res)
	case SpecialJudge:
		// special-judge handler
	case CustomTestcase:
		// custom-testcase handler
	default:
		r.logging.Error("unregistered handler")
	}

	return result
}
