package main

import (
	"context"
	"net/http"

	// _ "net/http/pprof"
	"os"

	"github.com/cranemont/iris/src/common/constants"
	"github.com/cranemont/iris/src/connector"
	"github.com/cranemont/iris/src/connector/rabbitmq"
	"github.com/cranemont/iris/src/handler"
	"github.com/cranemont/iris/src/router"
	"github.com/cranemont/iris/src/service/cache"
	"github.com/cranemont/iris/src/service/file"
	"github.com/cranemont/iris/src/service/logger"
	"github.com/cranemont/iris/src/service/sandbox"
	"github.com/cranemont/iris/src/service/testcase"
	"github.com/cranemont/iris/src/utils"
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

	uri := "amqp://" +
		utils.Getenv("RABBITMQ_DEFAULT_USER", "skku") + ":" +
		utils.Getenv("RABBITMQ_DEFAULT_PASS", "1234") + "@" +
		utils.Getenv("RABBITMQ_HOST", "localhost") + ":" +
		utils.Getenv("RABBITMQ_PORT", "5672") + "/"

	connector.Factory(
		connector.RABBIT_MQ,
		connector.Providers{Router: routeProvider, Logger: logProvider},
		rabbitmq.ConsumerConfig{
			AmqpURI:        uri,
			ConnectionName: utils.Getenv("RABBITMQ_CONSUMER_CONNECTION_NAME", "iris-consumer"),
			QueueName:      utils.Getenv("RABBITMQ_CONSUMER_QUEUE_NAME", "client.q.judge.submission"),
			Ctag:           utils.Getenv("RABBITMQ_CONSUMER_TAG", "consumer-tag"),
		},
		rabbitmq.ProducerConfig{
			AmqpURI:        uri,
			ConnectionName: utils.Getenv("RABBITMQ_PRODUCER_CONNECTION_NAME", "iris-producer"),
			ExchangeName:   utils.Getenv("RABBITMQ_PRODUCER_EXCHANGE_NAME", "iris.e.direct.judge"),
			RoutingKey:     utils.Getenv("RABBITMQ_PRODUCER_ROUTING_KEY", "judge.result"),
		},
	).Connect(context.Background())

	select {}
}
