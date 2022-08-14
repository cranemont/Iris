package judge

type Sandbox interface {
	Execute()
}

type Compiler interface {
	Compile(args CompileRequestDto)
}

type Runner interface {
	Run(args RunRequestDto)
}
