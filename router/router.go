package router

import (
	"fmt"
	"log"

	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/egress"
	"github.com/cranemont/judge-manager/handler"
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
	var result string
	switch handle {
	case JUDGE:
		judgeRequest, ok := data.(rmq.JudgeRequest)
		if !ok {
			log.Printf("JUDGE: %s", exception.ErrTypeAssertionFail)
			panic(nil)
		}

		res, err := r.judgeHandler.Handle(judgeRequest)
		if err != nil {
			log.Printf("JUDGE: handler error: %s", err)
		}
		result = r.judgeHandler.ResultToJson(res)
	case SPECIAL_JUDGE:
	case CUSTOM_TESTCASE:
	default:
		log.Printf("unregistered handler: %s", handle)
		panic(nil) // 컴파일시 잡아낼 수 없지만 어딘가 코드가 잘못되었기 때문에 panic
		// FIXME: string 말고 바로 호출할수있게 라우터에 해당 caller 만들기
		// 그리고 공통 response middleware를 만들고
		// result = handler.DefaultResult()
	}

	// publish result
	r.publisher.Publish(result)
	// goroutine으로 하는게 성능향상이 있나?
}
