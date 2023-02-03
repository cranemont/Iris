package main

import (
	"context"
	"net/http"

	// _ "net/http/pprof"
	"os"

	"github.com/cranemont/judge-manager/src/common/constants"
	"github.com/cranemont/judge-manager/src/connector"
	"github.com/cranemont/judge-manager/src/handler"
	"github.com/cranemont/judge-manager/src/router"
	"github.com/cranemont/judge-manager/src/service/cache"
	"github.com/cranemont/judge-manager/src/service/file"
	"github.com/cranemont/judge-manager/src/service/logger"
	"github.com/cranemont/judge-manager/src/service/sandbox"
	"github.com/cranemont/judge-manager/src/service/testcase"
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
	logProvider := logger.NewLogger(logger.Console, constants.Env((os.Getenv("APP_ENV"))))

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

	sb := sandbox.NewSandbox(logProvider)
	compiler := sandbox.NewCompiler(sb, langConfig, fileManager)
	runner := sandbox.NewRunner(sb, langConfig, fileManager)

	judgeHandler := handler.NewJudgeHandler(
		compiler,
		runner,
		testcaseManager,
		langConfig,
		fileManager,
		logProvider,
	)

	routeProvider := router.NewRouter(judgeHandler, logProvider)

	logProvider.Log(logger.INFO, "Server Started")

	// rabbitmq.NewConnector(consumer, producer, routeProvider, logProvider).Connect()
	connector.Factory(
		connector.RABBIT_MQ,
		connector.Providers{Router: routeProvider, Logger: logProvider},
	).Connect()

	select {}
}
