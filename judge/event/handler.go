package judgeEvent

import (
	"fmt"
	"log"
	"time"

	"github.com/cranemont/judge-manager/common/exception"
	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/event"
	"github.com/cranemont/judge-manager/file"
	"github.com/cranemont/judge-manager/judge"
	"github.com/cranemont/judge-manager/sandbox"
)

// controller의 역할. OnJudge, OnRun, OnOutput등으로 여러 상황 구분
type handler struct {
	funcMap      map[string](func(task *judge.Task) error)
	judger       *judge.Judger
	eventEmitter event.Emitter
	// fileManager  fileManager.FileManager
	config *sandbox.LanguageConfig
}

func NewHandler(
	judger *judge.Judger,
	eventEmitter event.Emitter,
	// fileManager fileManager.FileManager,
	config *sandbox.LanguageConfig,
) *handler {
	handler := &handler{
		judger:       judger,
		eventEmitter: eventEmitter,
		// fileManager:  fileManager,
		config: config,
	}
	funcMap := map[string](func(task *judge.Task) error){
		"OnExec": (*handler).OnExec,
		"OnExit": (*handler).OnExit,
	}
	handler.funcMap = funcMap
	return handler
}

// controller의 역할!
func (h *handler) OnExec(task *judge.Task) error {
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
	// error 처리
	fmt.Println("triggerring event")
	if err := h.eventEmitter.Emit(constants.TASK_EXITED, task); err != nil {
		return fmt.Errorf("onexec: event emit failed: %w", err)
	}
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

func (h *handler) OnExit(task *judge.Task) error {
	// 파일 삭제, task 결과 업데이트 등 정리작업
	file.RemoveDir(task.GetDir())
	fmt.Println(time.Since(task.StartedAt))
	if err := h.eventEmitter.Emit(constants.PUBLISH_RESULT, task); err != nil {
		return fmt.Errorf("onexit: event emit failed: %w", err)
	}

	return nil
}

func (h *handler) Call(funcName string, args interface{}) {
	// TODO: Refactor
	if fn, ok := h.funcMap[funcName]; ok {
		if _, ok := args.(*judge.Task); ok {
			fmt.Println("handler function calling... ", funcName)
			if err := fn(args.(*judge.Task)); err != nil {
				log.Printf("error on %s: %s", funcName, err)
			}
		} else {
			log.Printf("%s: invalid task data", exception.ErrTypeAssertionFail)
		}
	} else {
		log.Printf("unregistered function: %s", funcName)
	}
}
