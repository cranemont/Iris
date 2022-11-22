package main

import (
	"context"
	"net/http"

	// _ "net/http/pprof"
	"os"

	"github.com/cranemont/judge-manager/connector/rabbitmq"
	"github.com/cranemont/judge-manager/handler"
	"github.com/cranemont/judge-manager/router"
	"github.com/cranemont/judge-manager/service/cache"
	"github.com/cranemont/judge-manager/service/file"
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

func main() {
	// profile()
	zapLogger := logger.NewLogger(logger.Console, logger.Env(os.Getenv("APP_ENV")))

	ctx := context.Background()
	cache := cache.NewCache(ctx)

	// source := testcase.NewPreset()
	source := testcase.NewServer(
		os.Getenv("TESTCASE_SERVER_URL"),
		os.Getenv("TESTCASE_SERVER_AUTH_TOKEN"),
	)
	testcaseManager := testcase.NewManager(source, cache)

	fileManager := file.NewFileManager()
	langConfig := sandbox.NewLangConfig(fileManager)

	sb := sandbox.NewSandbox(zapLogger)
	compiler := sandbox.NewCompiler(sb, langConfig, fileManager)
	runner := sandbox.NewRunner(sb, langConfig, fileManager)

	judgeHandler := handler.NewJudgeHandler(
		compiler,
		runner,
		testcaseManager,
		langConfig,
		fileManager,
		zapLogger,
	)

	routing := router.NewRouter(judgeHandler, zapLogger)

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
	rabbitmq.NewConnector(consumer, producer, routing, zapLogger).Connect()
	select {}
}
