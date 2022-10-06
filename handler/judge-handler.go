package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cranemont/judge-manager/handler/judge"
	"github.com/cranemont/judge-manager/service/file"
	"github.com/cranemont/judge-manager/service/logger"
	"github.com/cranemont/judge-manager/service/sandbox"
)

var handler = "JudgeHandler"

type JudgeResposne struct {
	ServerStatusCode Code                  `json:"serverStatusCode"` // handler's status code
	Data             judge.JudgeTaskResult `json:"data"`
}

type JudgeRequest struct {
	SubmissionId int    `json:"submissionId"`
	Code         string `json:"code"`
	Language     string `json:"language"`
	ProblemId    int    `json:"problemId"`
	TimeLimit    int    `json:"timeLimit"`
	MemoryLimit  int    `json:"memoryLimit"`
}

type JudgeHandler struct {
	langConfig sandbox.LangConfig
	file       file.FileManager
	judger     *judge.Judger
	logger     *logger.Logger
}

func NewJudgeHandler(
	langConfig sandbox.LangConfig,
	file file.FileManager,
	judger *judge.Judger,
	logger *logger.Logger,
) *JudgeHandler {
	return &JudgeHandler{langConfig, file, judger, logger}
}

// handle top layer logical flow
func (j *JudgeHandler) Handle(req JudgeRequest) (result JudgeResposne, err error) {
	res := JudgeResposne{ServerStatusCode: INTERNAL_SERVER_ERROR, Data: judge.JudgeTaskResult{}}
	task := judge.NewTask(
		req.Code, req.Language, strconv.Itoa(req.ProblemId), req.TimeLimit, req.MemoryLimit,
	)
	dir := task.GetDir()

	defer func() {
		j.file.RemoveDir(dir)
		j.logger.Debug(fmt.Sprintf("task %s done: total time: %s", dir, time.Since(task.StartedAt)))
	}()

	if err := j.file.CreateDir(dir); err != nil {
		return res, fmt.Errorf("handler: %s: failed to create base directory: %w", handler, err)
	}

	srcPath, err := j.langConfig.MakeSrcPath(dir, task.GetLanguage())
	if err != nil {
		return res, fmt.Errorf("handler: %s: failed to create src path: %w", handler, err)
	}

	if err := j.file.CreateFile(srcPath, task.GetCode()); err != nil {
		return res, fmt.Errorf("handler: %s: failed to create src file: %w", handler, err)
	}

	err = j.judger.Judge(req.SubmissionId, task)
	if err != nil {
		if errors.Is(err, judge.ErrTestcaseGet) {
			res.ServerStatusCode = TESTCASE_GET_FAILED
		} else if errors.Is(err, judge.ErrCompile) {
			res.ServerStatusCode = COMPILE_ERROR
			res.Data = task.Result
			return res, nil
		} else {
			res.ServerStatusCode = INTERNAL_SERVER_ERROR
		} // run, grade 등 추가
		return res, fmt.Errorf("handler: judge failed: %w", err)
	} else {
		res.ServerStatusCode = SUCCESS
	}

	res.Data = task.Result
	return res, nil
}

func (h *JudgeHandler) ResultToJson(result JudgeResposne) []byte {
	res, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	return res
}
