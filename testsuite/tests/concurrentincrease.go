package tests

import (
	"github.com/nuweba/faasbenchmark/config"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
	"github.com/nuweba/httpbench"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func init() {
	Tests.Register(Test{Id: "ConcurrentIncreaseFastest", Fn: concurrentIncreaseFastest, RequiredStack: "coldstart", Description: "Test gradual growth of concurrent requests (up to 40) over time. Wait 50ms between hits."})
	Tests.Register(Test{Id: "ConcurrentIncreaseFast", Fn: concurrentIncreaseFast, RequiredStack: "coldstart", Description: "Test gradual growth of concurrent requests (up to 40) over time. Wait 150ms between hits."})
	Tests.Register(Test{Id: "ConcurrentIncreaseMedium", Fn: concurrentIncreaseMedium, RequiredStack: "coldstart", Description: "Test gradual growth of concurrent requests (up to 40) over time. Wait 250ms between hits."})
	Tests.Register(Test{Id: "ConcurrentIncreaseSlow", Fn: concurrentIncreaseSlow, RequiredStack: "coldstart", Description: "Test gradual growth of concurrent requests (up to 40) over time. Wait 350ms between hits."})
	Tests.Register(Test{Id: "ConcurrentIncreaseSlowest", Fn: concurrentIncreaseSlowest, RequiredStack: "coldstart", Description: "Test gradual growth of concurrent requests (up to 40) over time. Wait 450ms between hits."})
}

func concurrentIncreaseFastest(test *config.Test) {
	concurrentIncrease(test, config.Http{
		QueryParams: sleepQueryParam(200 * time.Millisecond),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(40, 50*time.Millisecond),
	})
}

func concurrentIncreaseFast(test *config.Test) {
	concurrentIncrease(test, config.Http{
		QueryParams: sleepQueryParam(200 * time.Millisecond),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(40, 150*time.Millisecond),
	})
}

func concurrentIncreaseMedium(test *config.Test) {
	concurrentIncrease(test, config.Http{
		QueryParams: sleepQueryParam(200 * time.Millisecond),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(40, 250*time.Millisecond),
	})
}

func concurrentIncreaseSlow(test *config.Test) {
	concurrentIncrease(test, config.Http{
		QueryParams: sleepQueryParam(200 * time.Millisecond),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(40, 350*time.Millisecond),
	})
}

func concurrentIncreaseSlowest(test *config.Test) {
	concurrentIncrease(test, config.Http{
		QueryParams: sleepQueryParam(200 * time.Millisecond),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(40, 450*time.Millisecond),
	})
}

func concurrentIncrease(test *config.Test, httpConfig config.Http) {
	if httpConfig.QueryParams == nil {
		httpConfig.QueryParams = new(url.Values)
	}
	if httpConfig.Headers == nil {
		httpConfig.Headers = new(http.Header)
	}
	if httpConfig.Body == nil {
		httpConfig.Body = new([]byte)
	}
	wg := &sync.WaitGroup{}
	for _, function := range test.Stack.ListFunctions() {
		hfConf, err := test.NewFunction(&httpConfig, function)
		if err != nil {
			continue
		}

		newReq := test.Config.Provider.NewFunctionRequest(hfConf.Test.Stack, hfConf.Function, hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
		trace := httpbench.New(newReq, hfConf.HttpConfig.Hook)

		wg.Add(1)
		go func() {
			defer wg.Done()
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, test.Config.Provider.HttpInvocationLatency)
		}()

		requestsResult := trace.RequestsForTimeGraph(*httpConfig.HitsGraph)

		wg.Wait()
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
	}
}
