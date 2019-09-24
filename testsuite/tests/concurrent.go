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

func init() {
	Tests.Register(Test{Id: "ConcurrentIncreasingLoad", Fn: C10FunctionsConcurrently1Each, RequiredStack: "identicalfunctions", Description: "Test concurrent load"})
	Tests.Register(Test{Id: "ConcurrentIncreasingLoadWindows", Fn: C10FunctionsConcurrently1Each, RequiredStack: "identicalfunctionswindows", Description: "Test concurrent load - azure functions on windows"})
}

func C10FunctionsConcurrently1Each(test *config.Test) {
	graphConcurrencyLimit := 34
	headers := http.Header{}
	body := []byte("")
	params := sleepQueryParam(mediumRuntime)

	httpConfig := &config.Http{
		SleepTime:   mediumRuntime,
		QueryParams: params,
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(graphConcurrencyLimit, Lvl2),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
		Body:        &body,
		Headers:     &headers,
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

