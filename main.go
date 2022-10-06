package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/cranemont/judge-manager/file"
	"github.com/cranemont/judge-manager/handler"
	"github.com/cranemont/judge-manager/handler/judge"
	"github.com/cranemont/judge-manager/ingress/mq"
	"github.com/cranemont/judge-manager/ingress/mq/rabbitmq"
	"github.com/cranemont/judge-manager/service/cache"
	"github.com/cranemont/judge-manager/service/logger"
	"github.com/cranemont/judge-manager/service/sandbox"
	"github.com/cranemont/judge-manager/service/testcase"
)

func profile() {
	// http.HandleFunc("/debug/pprof/", Index) // Profile Endpoint for Heap, Block, ThreadCreate, Goroutine, Mutex
	// http.HandleFunc("/debug/pprof/cmdline", Cmdline)
	// http.HandleFunc("/debug/pprof/profile", Profile) // Profile Endpoint for CPU
	// http.HandleFunc("/debug/pprof/symbol", Symbol)
	// http.HandleFunc("/debug/pprof/trace", Trace)
	go func() {
		http.ListenAndServe("localhost:6060", nil)
		// use <go tool pprof -http :8080 http://localhost:6060/debug/pprof/profile\?seconds\=30> to profile
	}()
}

func init() {

}

func main() {

	zapLogger := logger.NewLogger(logger.Console, logger.Env(os.Getenv("APP_ENV")))

	ctx := context.Background()
	cache := cache.NewCache(ctx)
	testcaseManager := testcase.NewManager(cache)

	fileManager := file.NewFileManager()
	langConfig := sandbox.NewLangConfig(fileManager)

	sb := sandbox.NewSandbox(zapLogger)
	compiler := sandbox.NewCompiler(sb, langConfig, fileManager)
	runner := sandbox.NewRunner(sb, langConfig, fileManager)

	judger := judge.NewJudger(
		compiler,
		runner,
		testcaseManager,
		zapLogger,
	)

	judgeHandler := handler.NewJudgeHandler(langConfig, fileManager, judger, zapLogger)
	// specialJudger
	// customTestcaseRunner 만들어서 같이 넣어주기

	rmqController := mq.NewRmqController(judgeHandler, zapLogger)

	uri := "amqp://" +
		os.Getenv("RABBITMQ_DEFAULT_USER") + ":" +
		os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" +
		os.Getenv("RABBITMQ_HOST") + ":" +
		os.Getenv("RABBITMQ_PORT") + "/"
	consumer, err := rabbitmq.NewConsumer(uri, "ctag", "go-consumer")
	if err != nil {
		panic(err)
	}

	producer, err := rabbitmq.NewProducer(uri, "go-producer", zapLogger)
	if err != nil {
		panic(err)
	}

	zapLogger.Info("Server started")
	mq.NewIngress(consumer, producer, rmqController, zapLogger).Activate()
	select {}

	// for debug
	// judgeRouter := router.NewRouter(
	// 	judgeHandler,
	// )

	// c := "#include <stdio.h>\n\nint main (void) {\nint a=0; scanf(\"%d\", &a);\nprintf(\"%d\t\\n\", a);\n\nreturn 0;\n}\n"
	// cpp := "#include<iostream>\n using namespace std;\n	int main() {\n	cout << \"1 1  \t\\n\";\n	return 0;\n}"
	// py := "a = input() \nprint(a)" // (\"1 1  \t\\n\")"
	// javaTimeout := "public class Main {\n	public static void main(String[] args) {\n		while(true) {\n		int i=0; i++;\n}\n}\n}"
	// java := "public class Main {\n	public static void main(String[] args) {\n		 System.out.println(\"1 1  \t\\n\");\n }\n}"

	// for {
	// 	var input string
	// 	fmt.Scanln(&input)
	//  problemId, err := strconv.Atoi(input)
	//  if err != nil {
	//    panic(err)
	//  }
	// 	submissionDto := handler.JudgeRequest{
	// 		// Code: "#include <stdio.h>\n\nint main (void) {\n  printf(\"Hello world!\\n\");\n  char buf[100];\n  scanf(\"%s\", buf);\n  printf(\"%s\\n\", buf);\n  return 0;\n}\n",
	// 		// Code: "#include <stdio.h>\n\nint main (void) {\nprintf(\"Hello world!\");\nreturn 0;\n}\n",
	// 		// Code: "#include <stdio.h>\n\nint main (void) {\nwhile(1) {printf(\"Hello world!\");}\nreturn 0;\n}\n",
	// 		Code:        c,
	// 		Language:    sandbox.C,
	// 		ProblemId:   problemId,
	// 		TimeLimit:   1000,
	// 		MemoryLimit: 256 * 1024 * 1024,
	// 	}
	// 	go judgeRouter.Route(router.JUDGE, submissionDto)
	// }
	// 여기서 rabbitMQ consumer가 돌고
	// 메시지 수신시 채점자 호출
}
