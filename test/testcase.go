package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/cranemont/judge-manager/cache"
	"github.com/cranemont/judge-manager/testcase"
)

func testcaseTest() {
	// reg([]byte("1 2   \t\r\n  1 2   \n"))
	// te([]byte("1 2   \t\r\n  1 2   \n"))
	ctx := context.Background()
	cache := cache.NewCache(ctx)
	testcaseManager := testcase.NewManager(cache)

	res, err := testcaseManager.GetTestcase("1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Data)
	fmt.Println(runtime.GOMAXPROCS(0))
}