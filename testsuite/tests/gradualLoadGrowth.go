package tests

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/config"
	"github.com/nuweba/faasbenchmark/provider"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
	"github.com/nuweba/httpbench"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	Var3 = 10 * time.Millisecond
	Var2 = 100 * time.Millisecond
	Var1 = 1000 * time.Millisecond

	sleep50ms  = 50 * time.Millisecond
	sleep500ms = 500 * time.Millisecond
	sleep5s    = 5000 * time.Millisecond

	maxConcurrent    = 40
	estimatedRuntime = 500 * time.Millisecond // used for functions we can't know the runtime of in advance
)

func generateDescription(duration time.Duration, resourceType string, isWarm bool) string {
	result := "gradually invoke more concurrent %s functions over time, with a %s delay between hits."
	if isWarm {
		result += " warmed up."
	}
	return fmt.Sprintf(result, resourceType, duration)
}

func init() {
	Tests.Register(Test{Id: "LinAscLoadVar3Sleep50ms", Fn: linAscLoadVar3Sleep50ms, RequiredStack: "sleep", Description: generateDescription(Var3, sleep50ms.String(), false)})
	Tests.Register(Test{Id: "LinAscLoadVar2Sleep50ms", Fn: linAscLoadVar2Sleep50ms, RequiredStack: "sleep", Description: generateDescription(Var2, sleep50ms.String(), false)})
	Tests.Register(Test{Id: "LinAscLoadVar1Sleep50ms", Fn: linAscLoadVar1Sleep50ms, RequiredStack: "sleep", Description: generateDescription(Var1, sleep50ms.String(), false)})

	Tests.Register(Test{Id: "LinAscLoadVar3Sleep500ms", Fn: linAscLoadVar3Sleep500ms, RequiredStack: "sleep", Description: generateDescription(Var3, sleep500ms.String(), false)})
	Tests.Register(Test{Id: "LinAscLoadVar2Sleep500ms", Fn: linAscLoadVar2Sleep500ms, RequiredStack: "sleep", Description: generateDescription(Var2, sleep500ms.String(), false)})
	Tests.Register(Test{Id: "LinAscLoadVar1Sleep500ms", Fn: linAscLoadVar1Sleep500ms, RequiredStack: "sleep", Description: generateDescription(Var1, sleep500ms.String(), false)})

	Tests.Register(Test{Id: "LinAscLoadVar3Sleep5s", Fn: linAscLoadVar3Sleep5s, RequiredStack: "sleep", Description: generateDescription(Var3, sleep5s.String(), false)})
	Tests.Register(Test{Id: "LinAscLoadVar2Sleep5s", Fn: linAscLoadVar2Sleep5s, RequiredStack: "sleep", Description: generateDescription(Var2, sleep5s.String(), false)})
	Tests.Register(Test{Id: "LinAscLoadVar1Sleep5s", Fn: linAscLoadVar1Sleep5s, RequiredStack: "sleep", Description: generateDescription(Var1, sleep5s.String(), false)})

	Tests.Register(Test{Id: "LinAscCPULoadVar3", Fn: linAscLoadVar3, RequiredStack: "cpustress", Description: generateDescription(Var3, "CPU intensive", false)})
	Tests.Register(Test{Id: "LinAscCPULoadVar2", Fn: linAscLoadVar2, RequiredStack: "cpustress", Description: generateDescription(Var2, "CPU intensive", false)})
	Tests.Register(Test{Id: "LinAscCPULoadVar1", Fn: linAscLoadVar1, RequiredStack: "cpustress", Description: generateDescription(Var1, "CPU intensive", false)})

	Tests.Register(Test{Id: "LinAscIOLoadVar3", Fn: linAscLoadVar3, RequiredStack: "iostress", Description: generateDescription(Var3, "IO intensive", false)})
	Tests.Register(Test{Id: "LinAscIOLoadVar2", Fn: linAscLoadVar2, RequiredStack: "iostress", Description: generateDescription(Var2, "IO intensive", false)})
	Tests.Register(Test{Id: "LinAscIOLoadVar1", Fn: linAscLoadVar1, RequiredStack: "iostress", Description: generateDescription(Var1, "IO intensive", false)})

	Tests.Register(Test{Id: "LinAscMemLoadVar3", Fn: linAscLoadVar3, RequiredStack: "memstress", Description: generateDescription(Var3, "memory intensive", false)})
	Tests.Register(Test{Id: "LinAscMemLoadVar2", Fn: linAscLoadVar2, RequiredStack: "memstress", Description: generateDescription(Var2, "memory intensive", false)})
	Tests.Register(Test{Id: "LinAscMemLoadVar1", Fn: linAscLoadVar1, RequiredStack: "memstress", Description: generateDescription(Var1, "memory intensive", false)})

	Tests.Register(Test{Id: "LinAscLoadWarmVar3Sleep50ms", Fn: linAscLoadWarmVar3Sleep50ms, RequiredStack: "sleep", Description: generateDescription(Var3, sleep50ms.String(), true)})
	Tests.Register(Test{Id: "LinAscLoadWarmVar2Sleep50ms", Fn: linAscLoadWarmVar2Sleep50ms, RequiredStack: "sleep", Description: generateDescription(Var2, sleep50ms.String(), true)})
	Tests.Register(Test{Id: "LinAscLoadWarmVar1Sleep50ms", Fn: linAscLoadWarmVar1Sleep50ms, RequiredStack: "sleep", Description: generateDescription(Var1, sleep50ms.String(), true)})

	Tests.Register(Test{Id: "LinAscLoadWarmVar3Sleep500ms", Fn: linAscLoadWarmVar3Sleep500ms, RequiredStack: "sleep", Description: generateDescription(Var3, sleep500ms.String(), true)})
	Tests.Register(Test{Id: "LinAscLoadWarmVar2Sleep500ms", Fn: linAscLoadWarmVar2Sleep500ms, RequiredStack: "sleep", Description: generateDescription(Var2, sleep500ms.String(), true)})
	Tests.Register(Test{Id: "LinAscLoadWarmVar1Sleep500ms", Fn: linAscLoadWarmVar1Sleep500ms, RequiredStack: "sleep", Description: generateDescription(Var1, sleep500ms.String(), true)})

	Tests.Register(Test{Id: "LinAscLoadWarmVar3Sleep5s", Fn: linAscLoadWarmVar3Sleep5s, RequiredStack: "sleep", Description: generateDescription(Var3, sleep5s.String(), true)})
	Tests.Register(Test{Id: "LinAscLoadWarmVar2Sleep5s", Fn: linAscLoadWarmVar2Sleep5s, RequiredStack: "sleep", Description: generateDescription(Var2, sleep5s.String(), true)})
	Tests.Register(Test{Id: "LinAscLoadWarmVar1Sleep5s", Fn: linAscLoadWarmVar1Sleep5s, RequiredStack: "sleep", Description: generateDescription(Var1, sleep5s.String(), true)})

	Tests.Register(Test{Id: "LinAscCPULoadWarmVar3", Fn: linAscLoadWarmVar3, RequiredStack: "cpustress", Description: generateDescription(Var3, "CPU intensive", true)})
	Tests.Register(Test{Id: "LinAscCPULoadWarmVar2", Fn: linAscLoadWarmVar2, RequiredStack: "cpustress", Description: generateDescription(Var2, "CPU intensive", true)})
	Tests.Register(Test{Id: "LinAscCPULoadWarmVar1", Fn: linAscLoadWarmVar1, RequiredStack: "cpustress", Description: generateDescription(Var1, "CPU intensive", true)})

	Tests.Register(Test{Id: "LinAscIOLoadWarmVar3", Fn: linAscLoadWarmVar3, RequiredStack: "iostress", Description: generateDescription(Var3, "IO intensive", true)})
	Tests.Register(Test{Id: "LinAscIOLoadWarmVar2", Fn: linAscLoadWarmVar2, RequiredStack: "iostress", Description: generateDescription(Var2, "IO intensive", true)})
	Tests.Register(Test{Id: "LinAscIOLoadWarmVar1", Fn: linAscLoadWarmVar1, RequiredStack: "iostress", Description: generateDescription(Var1, "IO intensive", true)})

	Tests.Register(Test{Id: "LinAscMemLoadWarmVar3", Fn: linAscLoadWarmVar3, RequiredStack: "memstress", Description: generateDescription(Var3, "memory intensive", true)})
	Tests.Register(Test{Id: "LinAscMemLoadWarmVar2", Fn: linAscLoadWarmVar2, RequiredStack: "memstress", Description: generateDescription(Var2, "memory intensive", true)})
	Tests.Register(Test{Id: "LinAscMemLoadWarmVar1", Fn: linAscLoadWarmVar1, RequiredStack: "memstress", Description: generateDescription(Var1, "memory intensive", true)})
}

func linAscLoadVar3Sleep5s(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep5s,
		QueryParams: sleepQueryParam(sleep5s),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var3),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar2Sleep5s(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep5s,
		QueryParams: sleepQueryParam(sleep5s),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var2),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar1Sleep5s(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep5s,
		QueryParams: sleepQueryParam(sleep5s),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var1),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar3Sleep500ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep500ms,
		QueryParams: sleepQueryParam(sleep500ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var3),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar2Sleep500ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep500ms,
		QueryParams: sleepQueryParam(sleep500ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var2),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar1Sleep500ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep500ms,
		QueryParams: sleepQueryParam(sleep500ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var1),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar3Sleep50ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep50ms,
		QueryParams: sleepQueryParam(sleep50ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var3),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar2Sleep50ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep50ms,
		QueryParams: sleepQueryParam(sleep50ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var2),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar1Sleep50ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep50ms,
		QueryParams: sleepQueryParam(sleep50ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var1),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar3(test *config.Test) {
	linAscLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Var3),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar2(test *config.Test) {
	linAscLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Var2),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadVar1(test *config.Test) {
	linAscLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Var1),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, false, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar3(test *config.Test) {
	linAscLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Var3),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar2(test *config.Test) {
	linAscLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Var2),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar1(test *config.Test) {
	linAscLoad(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Var1),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar3Sleep500ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep500ms,
		QueryParams: sleepQueryParam(sleep500ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var3),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar2Sleep500ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep500ms,
		QueryParams: sleepQueryParam(sleep500ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var2),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar1Sleep500ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep500ms,
		QueryParams: sleepQueryParam(sleep500ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var1),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar3Sleep5s(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep5s,
		QueryParams: sleepQueryParam(sleep5s),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var3),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar2Sleep5s(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep5s,
		QueryParams: sleepQueryParam(sleep5s),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var2),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar1Sleep5s(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep5s,
		QueryParams: sleepQueryParam(sleep5s),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var1),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar3Sleep50ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep50ms,
		QueryParams: sleepQueryParam(sleep50ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var3),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar2Sleep50ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep50ms,
		QueryParams: sleepQueryParam(sleep50ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var2),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoadWarmVar1Sleep50ms(test *config.Test) {
	linAscLoad(test, config.Http{
		SleepTime:   sleep50ms,
		QueryParams: sleepQueryParam(sleep50ms),
		TestType:    httpbench.RequestsForTimeGraph.String(),
		HitsGraph:   gradualHitGraph(maxConcurrent, Var1),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}, true, test.Config.Provider.HttpInvocationLatency)
}

func linAscLoad(test *config.Test, httpConfig config.Http, warmup bool, filter provider.RequestFilter) {
	headers := http.Header{}
	body := []byte{}
	queryParams := url.Values{}
	if httpConfig.QueryParams == nil {
		httpConfig.QueryParams = &queryParams
	}
	httpConfig.Headers = &headers
	httpConfig.Body = &body
	for _, function := range test.Stack.ListFunctions() {
		hfConf, err := test.NewFunction(&httpConfig, function)

		if err != nil {
			continue
		}

		if warmup {
			lastHit := (*hfConf.HttpConfig.HitsGraph)[len(*hfConf.HttpConfig.HitsGraph)-1]
			expectedRuntime := estimatedRuntime
			if httpConfig.SleepTime != 0 {
				expectedRuntime = httpConfig.SleepTime
			}
			requestsToSend := lastHit.Concurrent * (uint64(expectedRuntime/lastHit.Time) + 1)
			sendWarmup(hfConf, requestsToSend)
		}

		executeTest(hfConf, filter)
	}
}

func executeTest(hfConf *config.HttpFunction, filter provider.RequestFilter) {
	newReq := hfConf.Test.Config.Provider.NewFunctionRequest(hfConf.Test.Stack, hfConf.Function, hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
	wg := &sync.WaitGroup{}
	trace := httpbench.New(newReq, hfConf.HttpConfig.Hook)
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, filter)
	}()
	requestsResult := trace.RequestsForTimeGraph(*hfConf.HttpConfig.HitsGraph)
	wg.Wait()
	httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
}
