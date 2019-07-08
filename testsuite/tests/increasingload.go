package tests

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/config"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
	"github.com/nuweba/httpbench"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	Lvl3 = 10 * time.Millisecond
	Lvl2 = 100 * time.Millisecond
	Lvl1 = 1000 * time.Millisecond

	shortRuntime  = 50 * time.Millisecond
	mediumRuntime = 500 * time.Millisecond
	longRuntime   = 5000 * time.Millisecond

	maxConcurrent    = 40
)

func generateDescription(duration time.Duration, resourceType string) string {
	result := "Gradually invoke more concurrent %s over time with a %s delay between hits and benchmark their invocation overhead."
	return fmt.Sprintf(result, resourceType, duration)
}

func init() {
	Tests.Register(Test{Id: "IncreasingLoadLvl3", Fn: increasingLoadLvl3LongRuntime, RequiredStack: "sleep", Description: generateDescription(Lvl3, "long runtime functions")})
	Tests.Register(Test{Id: "IncreasingLoadLvl2", Fn: increasingLoadLvl2MedRuntime, RequiredStack: "sleep", Description: generateDescription(Lvl2, "medium runtime functions")})
	Tests.Register(Test{Id: "IncreasingLoadLvl1", Fn: increasingLoadLvl1ShortRuntime, RequiredStack: "sleep", Description: generateDescription(Lvl1, "short runtime functions")})

	Tests.Register(Test{Id: "IncreasingCPULoadLvl3", Fn: increasingLoadLvl3, RequiredStack: "cpustress", Description: generateDescription(Lvl3, "CPU intensive functions")})
	Tests.Register(Test{Id: "IncreasingCPULoadLvl2", Fn: increasingLoadLvl2, RequiredStack: "cpustress", Description: generateDescription(Lvl2, "CPU intensive functions")})
	Tests.Register(Test{Id: "IncreasingCPULoadLvl1", Fn: increasingLoadLvl1, RequiredStack: "cpustress", Description: generateDescription(Lvl1, "CPU intensive functions")})

	Tests.Register(Test{Id: "IncreasingIOLoadLvl3", Fn: increasingLoadLvl3, RequiredStack: "iostress", Description: generateDescription(Lvl3, "IO intensive functions")})
	Tests.Register(Test{Id: "IncreasingIOLoadLvl2", Fn: increasingLoadLvl2, RequiredStack: "iostress", Description: generateDescription(Lvl2, "IO intensive functions")})
	Tests.Register(Test{Id: "IncreasingIOLoadLvl1", Fn: increasingLoadLvl1, RequiredStack: "iostress", Description: generateDescription(Lvl1, "IO intensive functions")})

	Tests.Register(Test{Id: "IncreasingMemLoadLvl3", Fn: increasingLoadLvl3, RequiredStack: "memstress", Description: generateDescription(Lvl3, "memory intensive functions")})
	Tests.Register(Test{Id: "IncreasingMemLoadLvl2", Fn: increasingLoadLvl2, RequiredStack: "memstress", Description: generateDescription(Lvl2, "memory intensive functions")})
	Tests.Register(Test{Id: "IncreasingMemLoadLvl1", Fn: increasingLoadLvl1, RequiredStack: "memstress", Description: generateDescription(Lvl1, "memory intensive functions")})

	Tests.Register(Test{Id: "IncreasingNetLoadLvl3", Fn: increasingLoadLvl3, RequiredStack: "netstress", Description: generateDescription(Lvl3, "network intensive functions")})
	Tests.Register(Test{Id: "IncreasingNetLoadLvl2", Fn: increasingLoadLvl2, RequiredStack: "netstress", Description: generateDescription(Lvl2, "network intensive functions")})
	Tests.Register(Test{Id: "IncreasingNetLoadLvl1", Fn: increasingLoadLvl1, RequiredStack: "netstress", Description: generateDescription(Lvl1, "network intensive functions")})

	Tests.Register(Test{Id: "IncreasingProviderStorageLoadLvl3", Fn: increasingLoadLvl3, RequiredStack: "providerstorage", Description: generateDescription(Lvl3, "functions that upload to a provider storage service")})
	Tests.Register(Test{Id: "IncreasingProviderStorageLoadLvl2", Fn: increasingLoadLvl2, RequiredStack: "providerstorage", Description: generateDescription(Lvl2, "functions that upload to a provider storage service")})
	Tests.Register(Test{Id: "IncreasingProviderStorageLoadLvl1", Fn: increasingLoadLvl1, RequiredStack: "providerstorage", Description: generateDescription(Lvl1, "functions that upload to a provider storage service")})

	Tests.Register(Test{Id: "IncreasingLargeCodeLoadLvl3", Fn: increasingLoadLvl3, RequiredStack: "largecode", Description: generateDescription(Lvl3, "large codebase functions")})
	Tests.Register(Test{Id: "IncreasingLargeCodeLoadLvl2", Fn: increasingLoadLvl2, RequiredStack: "largecode", Description: generateDescription(Lvl2, "large codebase functions")})
	Tests.Register(Test{Id: "IncreasingLargeCodeLoadLvl1", Fn: increasingLoadLvl1, RequiredStack: "largecode", Description: generateDescription(Lvl1, "large codebase functions")})

	Tests.Register(Test{Id: "IncreasingLoggingLoadLvl3", Fn: increasingLoadLvl3, RequiredStack: "logging", Description: generateDescription(Lvl3, "logging functions")})
	Tests.Register(Test{Id: "IncreasingLoggingLoadLvl2", Fn: increasingLoadLvl2, RequiredStack: "logging", Description: generateDescription(Lvl2, "logging functions")})
	Tests.Register(Test{Id: "IncreasingLoggingLoadLvl1", Fn: increasingLoadLvl1, RequiredStack: "logging", Description: generateDescription(Lvl1, "logging functions")})

	Tests.Register(Test{Id: "IncreasingLoadOnVPCLvl3", Fn: increasingLoadLvl3, RequiredStack: "vpc", Description: generateDescription(Lvl3, "functions on a vpc")})
	Tests.Register(Test{Id: "IncreasingLoadOnVPCLvl2", Fn: increasingLoadLvl2, RequiredStack: "vpc", Description: generateDescription(Lvl2, "functions on a vpc")})
	Tests.Register(Test{Id: "IncreasingLoadOnVPCLvl1", Fn: increasingLoadLvl1, RequiredStack: "vpc", Description: generateDescription(Lvl1, "functions on a vpc")})

	Tests.Register(Test{Id: "IncreasingLargeRequestLoadLvl3", Fn: increasingLargeRequestLoadLvl3, RequiredStack: "largerequest", Description: generateDescription(Lvl3, "functions with a large (4mb) request body")})
	Tests.Register(Test{Id: "IncreasingLargeRequestLoadLvl2", Fn: increasingLargeRequestLoadLvl2, RequiredStack: "largerequest", Description: generateDescription(Lvl2, "functions with a large (4mb) request body")})
	Tests.Register(Test{Id: "IncreasingLargeRequestLoadLvl1", Fn: increasingLargeRequestLoadLvl1, RequiredStack: "largerequest", Description: generateDescription(Lvl1, "functions with a large (4mb) request body")})
}

func increasingLargeRequestLoadLvl3(test *config.Test) {
	increasingLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Lvl3),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func increasingLargeRequestLoadLvl2(test *config.Test) {
	increasingLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Lvl2),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func increasingLargeRequestLoadLvl1(test *config.Test) {
	increasingLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Lvl1),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func increasingLoadLvl3LongRuntime(test *config.Test) {
	increasingLoad(test, config.Http{
		SleepTime:   longRuntime,
		QueryParams: sleepQueryParam(longRuntime),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Lvl3),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func increasingLoadLvl2MedRuntime(test *config.Test) {
	increasingLoad(test, config.Http{
		SleepTime:   mediumRuntime,
		QueryParams: sleepQueryParam(mediumRuntime),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Lvl2),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func increasingLoadLvl1ShortRuntime(test *config.Test) {
	increasingLoad(test, config.Http{
		SleepTime:   shortRuntime,
		QueryParams: sleepQueryParam(shortRuntime),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Lvl1),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func increasingLoadLvl3(test *config.Test) {
	increasingLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Lvl3),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func increasingLoadLvl2(test *config.Test) {
	increasingLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Lvl2),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func increasingLoadLvl1(test *config.Test) {
	increasingLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Lvl1),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func increasingLoad(test *config.Test, httpConfig config.Http) {
	headers := http.Header{}
	body := []byte{}
	queryParams := url.Values{}
	if httpConfig.QueryParams == nil {
		httpConfig.QueryParams = &queryParams
	}
	httpConfig.QueryParams.Add("level", "1")
	httpConfig.Headers = &headers
	httpConfig.Body = &body
	for _, function := range test.Stack.ListFunctions() {
		hfConf, err := test.NewFunction(&httpConfig, function)

		if err != nil {
			continue
		}

		executeTest(hfConf)
	}
}

func executeTest(hfConf *config.HttpFunction) {
	newReq := hfConf.Test.Config.Provider.NewFunctionRequest(hfConf.Test.Stack, hfConf.Function, hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
	wg := &sync.WaitGroup{}
	trace := httpbench.New(newReq, hfConf.HttpConfig.Hook)
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, hfConf.Test.Config.Provider.HttpResult)
	}()
	requestsResult := trace.RequestsForTimeGraph(*hfConf.HttpConfig.HitsGraph)
	wg.Wait()
	httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
}
