package judge

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/cranemont/judge-manager/constants"
	"github.com/cranemont/judge-manager/judge/config"
)

type Compiler interface {
	Compile(task *Task) (int, error) // 얘는 task 몰라도 됨
	createSrcFile(srcPath string, code string) error
}

type compiler struct {
	sandbox Sandbox
	option  *config.CompileOption
}

func NewCompiler(sandbox Sandbox, option *config.CompileOption) *compiler {
	option.Init()
	return &compiler{sandbox, option}
}

func (c *compiler) Compile(task *Task) (int, error) {
	fmt.Println("Compile! from Compiler")

	options := c.option.Get(task.language) // 이게 된다고? private 아닌가? GetLanguage 가 필요없어?
	srcPath := constants.BASE_DIR + "/" + task.GetDir() + "/" + options.SrcName
	exePath := constants.BASE_DIR + "/" + task.GetDir() + "/" + options.ExeName

	// task.code로 srcName에 파일 생성, 얘는 다른곳에서 생성해줘야됨. 컴파일이 아님
	if err := c.createSrcFile(srcPath, task.code); err != nil {
		// ENUM으로 변경, result code 반환
		fmt.Println("error from createSrcFile")
		return -1, err
	}

	// option에서 바로 매칭시켜서 sadnbox인자 넘겨주기

	args := strings.Replace(options.Args, "{srcPath}", srcPath, 1)
	args = strings.Replace(args, "{exePath}", exePath, 1)
	argSlice := strings.Split(args, " ")
	// sandbox 받지말고 그냥 여기서 arg처리한다음에 libjudger 실행하기

	c.sandbox.Execute(
		&SandboxArgs{
			ExePath:     options.CompilerPath,
			MaxCpuTime:  options.MaxCpuTime,
			MaxRealTime: options.MaxRealTime,
			MaxMemory:   options.MaxMemory,
			Args:        argSlice,
		})
	time.Sleep(time.Second * 2)
	// 채널로 결과반환?
	return 0, nil
}

func (c *compiler) createSrcFile(srcPath string, code string) error {
	err := ioutil.WriteFile(srcPath, []byte(code), constants.BASE_FILE_MODE)
	if err != nil {
		fmt.Println("파일 생성 실패", err)
		return err
	}
	return nil
}
