package router

import (
	"fmt"
	"log"

	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/egress"
	"github.com/cranemont/judge-manager/handler"
	"github.com/cranemont/judge-manager/handler/judge"
	"github.com/cranemont/judge-manager/ingress/rmq"
)

// Judge Mode
const (
	JUDGE = 0 + iota
	SPECIAL_JUDGE
	CUSTOM_TESTCASE
)

// controller의 역할. OnJudge, OnRun, OnOutput등으로 여러 상황 구분
type router struct {
	judgeHandler *handler.JudgeHandler
	publisher    egress.Publisher
}

func NewRouter(
	judgeHandler *handler.JudgeHandler,
	publisher egress.Publisher,
) *router {
	return &router{
		judgeHandler: judgeHandler,
		publisher:    publisher,
	}
}

// 요청을 받고 최종 response를 내보내는 책임
// logging, publish
// router, controller?
func (r *router) Route(funcName string, dto interface{}) {

	// funcName별로 해당하는 Type의 Task로 만들어서 아래 함수 호출
	fmt.Println("handler function calling... ", funcName)
	var result string // []byte
	switch funcName {
	case "Judge":
		// 여기서 직접 호출하는게 아니라 judge패키지에 handler가 있어서 거기서 result json만 받아옴
		judgeRequest, ok := dto.(rmq.JudgeRequest)
		if !ok {
			log.Printf("%s: %s", funcName, exception.ErrTypeAssertionFail)
			return
		}
		task := judge.NewTask(judgeRequest)
		// 여기서 dto을 type assertion
		err := r.judgeHandler.Handle(task) // result 받기?
		if err != nil {
			log.Printf("error on %s: %s", funcName, err)
		}
		result = task.ResultToJson()
	case "SpecialJudge":
	case "CustomTestcaseRun":
	default:
		log.Printf("unregistered function: %s", funcName)
	}

	// publish result
	r.publisher.Publish(result)
	// goroutine으로 하는게 성능향상이 있나?
}
