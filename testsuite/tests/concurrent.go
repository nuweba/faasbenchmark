package tests

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/config"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/httpbench"
	"net/http"
	"sync"
	"time"
)

func init() {
	Tests.Register(Test{Id: "10FunctionsConcurrently1Each", Fn: C10FunctionsConcurrently1Each, RequiredStack: "identicalfunctions", Description: "Test concurrent load"})
	Tests.Register(Test{Id: "10FunctionsConcurrently1EachWindows", Fn: C10FunctionsConcurrently1Each, RequiredStack: "identicalfunctionswindows", Description: "Test concurrent load - azure functions on windows"})
	Tests.Register(Test{Id: "1Function10Concurrent", Fn: C1Function10Concurrent, RequiredStack: "singlefunction", Description: "Test concurrent load"})
	Tests.Register(Test{Id: "1Function10ConcurrentWindows", Fn: C1Function10Concurrent, RequiredStack: "singlefunctionwindows", Description: "Test concurrent load - azure functions on windows"})
}

func C10FunctionsConcurrently1Each(test *config.Test) {
	sleep := 2000 * time.Millisecond
	headers := http.Header{}
	body := []byte("")
	params := sleepQueryParam(sleep)

	httpConfig := &config.Http{
		SleepTime:        sleep,
		Hook:             test.Config.Provider.HttpInvocationTriggerStage(),
		QueryParams:      params,
		Headers:          &headers,
		Duration:         0,
		RequestDelay:     20 * time.Millisecond,
		ConcurrencyLimit: 1,
		Body:             &body,
		TestType:         httpbench.ConcurrentRequestsSyncedOnce.String(),
	}

	var m sync.Mutex
	c := sync.NewCond(&m)

	var start, end sync.WaitGroup

	for _, function := range test.Stack.ListFunctions() {
		end.Add(1)
		start.Add(1)
		go func(function stack.Function) {

			defer end.Done()

			var wg sync.WaitGroup

			hfConf, err := test.NewFunction(httpConfig, function)
			if err != nil {
				fmt.Println(err)
				return
			}

			newReq := test.Config.Provider.NewFunctionRequest(hfConf.Test.Stack, hfConf.Function, hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
			trace := httpbench.New(newReq, hfConf.HttpConfig.Hook)

			wg.Add(1)
			go func() {
				defer wg.Done()
				httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, hfConf.Test.Config.Provider.HttpResult, test.Config.Debug)
			}()

			c.L.Lock()
			start.Done()

			c.Wait()

			c.L.Unlock()

			requestsResult := trace.ConcurrentRequestsSyncedOnce(hfConf.HttpConfig.ConcurrencyLimit, hfConf.HttpConfig.RequestDelay)
			wg.Wait()
			httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
		}(function)
	}

	start.Wait()

	c.Broadcast()

	end.Wait()
}

func C1Function10Concurrent(test *config.Test) {
	sleep := 2000 * time.Millisecond
	qParams := sleepQueryParam(sleep)
	headers := http.Header{}
	body := []byte("")

	httpConfig := &config.Http{
		SleepTime:        sleep,
		Hook:             test.Config.Provider.HttpInvocationTriggerStage(),
		QueryParams:      qParams,
		Headers:          &headers,
		Duration:         0,
		RequestDelay:     20 * time.Millisecond,
		ConcurrencyLimit: 10,
		Body:             &body,
		TestType:         httpbench.ConcurrentRequestsSyncedOnce.String(),
	}
	wg := &sync.WaitGroup{}
	for _, function := range test.Stack.ListFunctions() {
		hfConf, err := test.NewFunction(httpConfig, function)
		if err != nil {
			continue
		}

		newReq := test.Config.Provider.NewFunctionRequest(hfConf.Test.Stack, hfConf.Function, hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
		trace := httpbench.New(newReq, hfConf.HttpConfig.Hook)

		wg.Add(1)
		go func() {
			defer wg.Done()
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, hfConf.Test.Config.Provider.HttpResult, test.Config.Debug)
		}()

		requestsResult := trace.ConcurrentRequestsSyncedOnce(hfConf.HttpConfig.ConcurrencyLimit, hfConf.HttpConfig.RequestDelay)
		wg.Wait()
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
	}
}
