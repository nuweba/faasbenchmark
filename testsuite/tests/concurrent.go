package tests

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/config"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/httpbench"
	"net/http"
	"sync"
)

/*
  We use less than the default maxConcurrent to ensure the total number of requests running concurrently is not more
  than ~800 - the number of requests sent in a regular increasingLoad test
*/
const (
	Lvl2ConcurrencyLimit = 34
	Lvl3ConcurrencyLimit = 18
)

func init() {
	Tests.Register(Test{Id: "ConcurrentIncreasingLoadLvl1", Fn: ConcurrentIncreasingLoadLvl1, RequiredStack: "identicalfunctions", Description: "Test concurrent load"})
	Tests.Register(Test{Id: "ConcurrentIncreasingLoadWindowsLvl1", Fn: ConcurrentIncreasingLoadLvl1, RequiredStack: "identicalfunctionswindows", Description: "Test concurrent load - azure functions on windows"})
	Tests.Register(Test{Id: "ConcurrentIncreasingLoadLvl2", Fn: ConcurrentIncreasingLoadLvl2, RequiredStack: "identicalfunctions", Description: "Test concurrent load"})
	Tests.Register(Test{Id: "ConcurrentIncreasingLoadWindowsLvl2", Fn: ConcurrentIncreasingLoadLvl2, RequiredStack: "identicalfunctionswindows", Description: "Test concurrent load - azure functions on windows"})
	Tests.Register(Test{Id: "ConcurrentIncreasingLoadLvl3", Fn: ConcurrentIncreasingLoadLvl3, RequiredStack: "identicalfunctions", Description: "Test concurrent load"})
	Tests.Register(Test{Id: "ConcurrentIncreasingLoadWindowsLvl3", Fn: ConcurrentIncreasingLoadLvl3, RequiredStack: "identicalfunctionswindows", Description: "Test concurrent load - azure functions on windows"})
}

func ConcurrentIncreasingLoadLvl1(test *config.Test){
	runtime := shortRuntime
	headers := http.Header{}
	body := []byte("")
	params := sleepQueryParam(runtime)
	httpConfig := &config.Http{
		SleepTime:   runtime,
		QueryParams: params,
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Lvl1),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
		Body:        &body,
		Headers:     &headers,
	}
	testConcurrently(test, httpConfig)
}

func ConcurrentIncreasingLoadLvl2(test *config.Test){
	runtime := mediumRuntime
	headers := http.Header{}
	body := []byte("")
	params := sleepQueryParam(runtime)
	httpConfig := &config.Http{
		SleepTime:   runtime,
		QueryParams: params,
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(Lvl2ConcurrencyLimit, Lvl2),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
		Body:        &body,
		Headers:     &headers,
	}
	testConcurrently(test, httpConfig)
}

func ConcurrentIncreasingLoadLvl3(test *config.Test){
	runtime := longRuntime
	headers := http.Header{}
	body := []byte("")
	params := sleepQueryParam(runtime)
	httpConfig := &config.Http{
		SleepTime:   runtime,
		QueryParams: params,
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(Lvl3ConcurrencyLimit, Lvl3),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
		Body:        &body,
		Headers:     &headers,
	}
	testConcurrently(test, httpConfig)
}

func testConcurrently(test *config.Test, httpConfig *config.Http) {
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
				httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, hfConf.Test.Config.Provider.HttpResult)
			}()

			c.L.Lock()
			start.Done()

			c.Wait()

			c.L.Unlock()

			requestsResult := trace.RequestsForTimeGraph(*hfConf.HttpConfig.HitsGraph)
			wg.Wait()
			httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
		}(function)
	}

	start.Wait()

	c.Broadcast()

	end.Wait()
}

