package mq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cranemont/judge-manager/handler"
	"github.com/cranemont/judge-manager/ingress/mq/producer"
)

// Judge Handler
const (
	JUDGE           = "Judge"
	SPECIAL_JUDGE   = "SpecialJudge"
	CUSTOM_TESTCASE = "CustomTestcase"
)

type RmqController interface {
	Call(handle string, data interface{}) []byte
}

// controller의 역할. OnJudge, OnRun, OnOutput등으로 여러 상황 구분
type rmqController struct {
	judgeHandler *handler.JudgeHandler
}

func NewRmqController(
	judgeHandler *handler.JudgeHandler,
) *rmqController {
	return &rmqController{
		judgeHandler: judgeHandler,
	}
}

// 요청을 받고 최종 response를 내보내는 책임
// logging, publish
// router, controller?
func (r *rmqController) Call(handle string, data interface{}) []byte {
	fmt.Println("From Controller: ", handle)
	var result []byte
	switch handle {
	case JUDGE:
		req := handler.JudgeRequest{}
		err := json.Unmarshal(data.([]byte), &req)
		// judgeRequest, ok := data.(handler.JudgeRequest)
		if err != nil {
			log.Printf("JUDGE: %s", err)
			fmt.Println(string(data.([]byte)))
			// panic(nil)
			break
		}

		res, err := r.judgeHandler.Handle(req)
		if err != nil {
			log.Printf("JUDGE: handler error: %s", err)
			break
		}
		result = r.judgeHandler.ResultToJson(res)
	case SPECIAL_JUDGE:
		// special-judge handler
	case CUSTOM_TESTCASE:
		// custom-testcase handler
	default:
		log.Printf("unregistered handler: %s", handle)
		fmt.Println(data)
		// panic(nil) // 컴파일시 잡아낼 수 없지만 어딘가 코드가 잘못되었기 때문에 panic
		// FIXME: string 말고 바로 호출할수있게 라우터에 해당 caller 만들기
		// 그리고 공통 response middleware를 만들고
		// result = handler.DefaultResult()
	}

	// publish result
	return result
}

func publish(result string) {
	uri := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + "rabbitmq" + ":5672/"
	done := make(chan bool)
	if err := producer.Publish(done, uri, result); err != nil {
		log.Println(fmt.Errorf("%w", err))
	}
}
