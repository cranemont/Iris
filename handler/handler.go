package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/cranemont/judge-manager/egress"
	"github.com/cranemont/judge-manager/file"
	"github.com/cranemont/judge-manager/handler/judge"
	"github.com/cranemont/judge-manager/ingress/rmq"
	"github.com/cranemont/judge-manager/sandbox"
)

// controller의 역할. OnJudge, OnRun, OnOutput등으로 여러 상황 구분
type handler struct {
	funcMap   map[string](func(task *judge.Task) error)
	judger    *judge.Judger
	config    *sandbox.LanguageConfig
	publisher egress.Publisher
}

func NewHandler(
	judger *judge.Judger,
	config *sandbox.LanguageConfig,
	publisher egress.Publisher,
) *handler {
	handler := &handler{
		judger:    judger,
		config:    config,
		publisher: publisher,
	}
	funcMap := map[string](func(task *judge.Task) error){
		"Judge": (*handler).Judge,
	}
	handler.funcMap = funcMap
	return handler
}

func (h *handler) SpecialJudge() error {
	return nil
}

func (h *handler) CustomTestcaseRun() error {
	return nil
}

// controller의 역할!
// 여기서 JudgeResult객체 관리
func (h *handler) Judge(task *judge.Task) error {
	defer file.RemoveDir(task.GetDir())
	// Result의 status code는 여기서 관리,
	// task에 들어가있을게 아님. sentinel error로 잡아내기?

	task.StartedAt = time.Now()
	dir := task.GetDir()
	// 폴더 생성
	if err := file.CreateDir(dir); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	srcPath, err := h.config.MakeSrcPath(dir, task.GetLanguage())
	if err != nil {
		return err
	}
	if err := h.createSrcFile(srcPath, task.GetCode()); err != nil {
		return err
	}

	if err := h.judger.Judge(task); err != nil {
		return fmt.Errorf("onexec: %w", err)
	}
	// error 처리, defer에서 무조건 실행되도록 하기(폴더제거)
	fmt.Println("triggerring event")
	return nil
}

func (h *handler) createSrcFile(srcPath string, code string) error {
	if err := file.CreateFile(srcPath, code); err != nil {
		// ENUM으로 변경, result code 반환
		err := fmt.Errorf("failed to create src file: %s", err)
		return err
	}
	return nil
}

// 요청을 받고 최종 response를 내보내는 책임
func (h *handler) Handle(funcName string, dto rmq.SubmissionDto) {
	// TODO: Refactor
	// fn, ok := h.funcMap[funcName]
	// if !ok {
	// 	log.Printf("unregistered function: %s", funcName)
	// }

	// funcName별로 해당하는 Type의 Task로 만들어서 아래 함수 호출
	fmt.Println("handler function calling... ", funcName)
	var result string // []byte
	switch funcName {
	case "Judge":
		task := judge.NewTask(dto)
		err := h.Judge(task) // result 받기?
		if err != nil {
			log.Printf("error on %s: %s", funcName, err)
		}
		result = task.ResultToJson()
	case "SpecialJudge":
	case "CustomTestcaseRun":
	}

	// publish result
	h.publisher.Publish(result)
	// goroutine으로 하는게 성능향상이 있나?
}
