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

// Judge Handler
const (
	JUDGE           = "Judge"
	SPECIAL_JUDGE   = "SpecialJudge"
	CUSTOM_TESTCASE = "CustomTestcase"
)

type Router interface {
	Route(handle string, data interface{})
}

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
func (r *router) Route(handle string, data interface{}) {
	fmt.Println("From Router: ", handle)
	var result string // []byte
	switch handle {
	case JUDGE:
		judgeRequest, ok := data.(rmq.JudgeRequest)
		if !ok {
			log.Printf("JUDGE: %s", exception.ErrTypeAssertionFail)
			return
		}
		task := judge.NewTask(judgeRequest)
		err := r.judgeHandler.Handle(task)
		if err != nil {
			log.Printf("JUDGE: handler error: %s", err)
			return
		}
		result, err = task.ResultToJson()
		if err != nil {
			log.Printf("JUDGE: %s", err)
			return
		}
	case SPECIAL_JUDGE:
	case CUSTOM_TESTCASE:
	default:
		log.Printf("unregistered handler: %s", handle)
	}

	// publish result
	r.publisher.Publish(result)
	// goroutine으로 하는게 성능향상이 있나?
}
