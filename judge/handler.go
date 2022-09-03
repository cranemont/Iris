package judge

import (
	"fmt"
	"log"
	"time"

	"github.com/cranemont/judge-manager/file"
	"github.com/cranemont/judge-manager/sandbox"
)

// controller의 역할. OnJudge, OnRun, OnOutput등으로 여러 상황 구분
type handler struct {
	funcMap map[string](func(task *Task) error)
	judger  *Judger
	config  *sandbox.LanguageConfig
}

func NewHandler(
	judger *Judger,
	config *sandbox.LanguageConfig,
) *handler {
	handler := &handler{
		judger: judger,
		config: config,
	}
	funcMap := map[string](func(task *Task) error){
		"Judge": (*handler).Judge,
	}
	handler.funcMap = funcMap
	return handler
}

// controller의 역할!
// 여기서 JudgeResult객체 관리
func (h *handler) Judge(task *Task) error {
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
func (h *handler) Handle(funcName string, task *Task) {
	// TODO: Refactor
	if fn, ok := h.funcMap[funcName]; ok {
		fmt.Println("handler function calling... ", funcName)
		if err := fn(task); err != nil {
			log.Printf("error on %s: %s", funcName, err)
		}
	} else {
		log.Printf("unregistered function: %s", funcName)
	}
}
