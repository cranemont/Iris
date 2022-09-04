package handler

import (
	"fmt"
	"time"

	"github.com/cranemont/judge-manager/file"
	"github.com/cranemont/judge-manager/handler/judge"
	"github.com/cranemont/judge-manager/sandbox"
)

var handler = "JudgeHandler"

type JudgeHandler struct {
	judger *judge.Judger
	config *sandbox.LanguageConfig
}

func NewJudgeHandler(
	judger *judge.Judger,
	config *sandbox.LanguageConfig,
) *JudgeHandler {
	return &JudgeHandler{
		judger: judger,
		config: config,
	}
}

// handle top layer logical flow
func (h *JudgeHandler) Handle(task *judge.Task) error {
	defer func() {
		file.RemoveDir(task.GetDir())
		fmt.Println(time.Since(task.StartedAt))
	}()
	// Result의 status code는 여기서 관리,
	// task에 들어가있을게 아님. sentinel error로 잡아내기?

	task.StartedAt = time.Now()
	dir := task.GetDir()

	if err := file.CreateDir(dir); err != nil {
		return fmt.Errorf("%s: failed to create directory: %w", handler, err)
	}

	srcPath, err := h.config.MakeSrcPath(dir, task.GetLanguage())
	if err != nil {
		return fmt.Errorf("%s: failed to create src path: %w", handler, err)
	}
	if err := file.CreateFile(srcPath, task.GetCode()); err != nil {
		return fmt.Errorf("%s: failed to create src file: %w", handler, err)
	}

	if err := h.judger.Judge(task); err != nil {
		return fmt.Errorf("%s: judge failed: %w", handler, err)
	}
	// error 처리, defer에서 무조건 실행되도록 하기(폴더제거)
	fmt.Println("JudgeHandler: Handle Done!")
	return nil
}
