package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cranemont/judge-manager/file"
	"github.com/cranemont/judge-manager/handler/judge"
	"github.com/cranemont/judge-manager/ingress/rmq"
	"github.com/cranemont/judge-manager/sandbox"
)

var handler = "JudgeHandler"

type JudgeResult struct {
	StatusCode Code                  `json:"statusCode"` // handler's status code
	Data       judge.JudgeTaskResult `json:"data"`
}

type JudgeHandler struct {
	langConfig sandbox.LangConfig
	file       file.FileManager
	judger     *judge.Judger
}

func NewJudgeHandler(
	langConfig sandbox.LangConfig,
	file file.FileManager,
	judger *judge.Judger,
) *JudgeHandler {
	return &JudgeHandler{langConfig, file, judger}
}

// handle top layer logical flow
func (h *JudgeHandler) Handle(request rmq.JudgeRequest) (result JudgeResult, err error) {
	res := JudgeResult{StatusCode: INTERNAL_SERVER_ERROR, Data: judge.JudgeTaskResult{}}
	task := judge.NewTask(request)
	task.StartedAt = time.Now()
	dir := task.GetDir()

	defer func() {
		h.file.RemoveDir(task.GetDir())
		fmt.Println(time.Since(task.StartedAt)) // for debug
	}()

	if err := h.file.CreateDir(dir); err != nil {
		return res, fmt.Errorf("%s: failed to create base directory: %w", handler, err)
	}

	srcPath, err := h.langConfig.MakeSrcPath(dir, task.GetLanguage())
	if err != nil {
		return res, fmt.Errorf("%s: failed to create src path: %w", handler, err)
	}

	if err := h.file.CreateFile(srcPath, task.GetCode()); err != nil {
		return res, fmt.Errorf("%s: failed to create src file: %w", handler, err)
	}

	err = h.judger.Judge(task)
	if err != nil {
		if errors.Is(err, judge.ErrTestcaseGet) {
			res.StatusCode = TESTCASE_GET_FAILED
		} else if errors.Is(err, judge.ErrCompile) {
			res.StatusCode = COMPILE_ERROR
			// 이때는 아래 코드 실행해야됨(오류로 던지는게 아님)
			res.Data = task.Result
			return res, nil
		} else {
			res.StatusCode = INTERNAL_SERVER_ERROR
		} // run, grade 등 추가
		return res, fmt.Errorf("%s: judge failed: %w", handler, err)
	} else {
		res.StatusCode = SUCCESS
	}

	res.Data = task.Result
	fmt.Println("JudgeHandler: Handle Done!")
	return res, nil
}

func (h *JudgeHandler) ResultToJson(result JudgeResult) string {
	res, err := json.Marshal(result)
	if err != nil {
		// 적절한 err 처리
		panic(err)
	}
	return string(res)
}
