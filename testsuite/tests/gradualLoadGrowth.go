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
	Fastest = (iota*100 + 50) * time.Millisecond
	Fast
	Medium
	Slow
	Slowest

	sleepTime     = 200 * time.Millisecond
	maxConcurrent = 40
)

func generateDescription(duration time.Duration, resourceType string, isPostBurst bool) string {
	result := ""
	requestType := ""
	if isPostBurst {
		result += "send a burst of concurrent requests without benchmarking it. then, "
	}
	result += "gradually send more concurrent %srequests over time, with a %s delay between hits."
	if resourceType != "" {
		requestType = resourceType + " intensive "
	}
	return fmt.Sprintf(result, requestType, duration)
}

func init() {
	Tests.Register(Test{Id: "GradualLoadGrowthFastest", Fn: gradualLoadGrowthFastest, RequiredStack: "sleep", Description: generateDescription(Fastest, "", false)})
	Tests.Register(Test{Id: "GradualLoadGrowthFast", Fn: gradualLoadGrowthFast, RequiredStack: "sleep", Description: generateDescription(Fast, "", false)})
	Tests.Register(Test{Id: "GradualLoadGrowthMedium", Fn: gradualLoadGrowthMedium, RequiredStack: "sleep", Description: generateDescription(Medium, "", false)})
	Tests.Register(Test{Id: "GradualLoadGrowthSlow", Fn: gradualLoadGrowthSlow, RequiredStack: "sleep", Description: generateDescription(Slow, "", false)})
	Tests.Register(Test{Id: "GradualLoadGrowthSlowest", Fn: gradualLoadGrowthSlowest, RequiredStack: "sleep", Description: generateDescription(Slowest, "", false)})

	Tests.Register(Test{Id: "GradualCPULoadGrowthFastest", Fn: gradualLoadGrowthFastest, RequiredStack: "cpustress", Description: generateDescription(Fastest, "CPU", false)})
	Tests.Register(Test{Id: "GradualCPULoadGrowthFast", Fn: gradualLoadGrowthFast, RequiredStack: "cpustress", Description: generateDescription(Fast, "CPU", false)})
	Tests.Register(Test{Id: "GradualCPULoadGrowthMedium", Fn: gradualLoadGrowthMedium, RequiredStack: "cpustress", Description: generateDescription(Medium, "CPU", false)})
	Tests.Register(Test{Id: "GradualCPULoadGrowthSlow", Fn: gradualLoadGrowthSlow, RequiredStack: "cpustress", Description: generateDescription(Slow, "CPU", false)})
	Tests.Register(Test{Id: "GradualCPULoadGrowthSlowest", Fn: gradualLoadGrowthSlowest, RequiredStack: "cpustress", Description: generateDescription(Slowest, "CPU", false)})

	Tests.Register(Test{Id: "GradualIOLoadGrowthFastest", Fn: gradualLoadGrowthFastest, RequiredStack: "iostress", Description: generateDescription(Fastest, "IO", false)})
	Tests.Register(Test{Id: "GradualIOLoadGrowthFast", Fn: gradualLoadGrowthFast, RequiredStack: "iostress", Description: generateDescription(Fast, "IO", false)})
	Tests.Register(Test{Id: "GradualIOLoadGrowthMedium", Fn: gradualLoadGrowthMedium, RequiredStack: "iostress", Description: generateDescription(Medium, "IO", false)})
	Tests.Register(Test{Id: "GradualIOLoadGrowthSlow", Fn: gradualLoadGrowthSlow, RequiredStack: "iostress", Description: generateDescription(Slow, "IO", false)})
	Tests.Register(Test{Id: "GradualIOLoadGrowthSlowest", Fn: gradualLoadGrowthSlowest, RequiredStack: "iostress", Description: generateDescription(Slowest, "IO", false)})

	Tests.Register(Test{Id: "GradualMemLoadGrowthFastest", Fn: gradualLoadGrowthFastest, RequiredStack: "memstress", Description: generateDescription(Fastest, "memory", false)})
	Tests.Register(Test{Id: "GradualMemLoadGrowthFast", Fn: gradualLoadGrowthFast, RequiredStack: "memstress", Description: generateDescription(Fast, "memory", false)})
	Tests.Register(Test{Id: "GradualMemLoadGrowthMedium", Fn: gradualLoadGrowthMedium, RequiredStack: "memstress", Description: generateDescription(Medium, "memory", false)})
	Tests.Register(Test{Id: "GradualMemLoadGrowthSlow", Fn: gradualLoadGrowthSlow, RequiredStack: "memstress", Description: generateDescription(Slowest, "memory", false)})
	Tests.Register(Test{Id: "GradualMemLoadGrowthSlowest", Fn: gradualLoadGrowthSlowest, RequiredStack: "memstress", Description: generateDescription(Slow, "memory", false)})

	Tests.Register(Test{Id: "GradualLoadGrowthPostBurstFastest", Fn: gradualLoadGrowthPostBurstFastest, RequiredStack: "sleep", Description: generateDescription(Fastest, "", true)})
	Tests.Register(Test{Id: "GradualLoadGrowthPostBurstFast", Fn: gradualLoadGrowthPostBurstFast, RequiredStack: "sleep", Description: generateDescription(Fast, "", true)})
	Tests.Register(Test{Id: "GradualLoadGrowthPostBurstMedium", Fn: gradualLoadGrowthPostBurstMedium, RequiredStack: "sleep", Description: generateDescription(Medium, "", true)})
	Tests.Register(Test{Id: "GradualLoadGrowthPostBurstSlow", Fn: gradualLoadGrowthPostBurstSlow, RequiredStack: "sleep", Description: generateDescription(Slow, "", true)})
	Tests.Register(Test{Id: "GradualLoadGrowthPostBurstSlowest", Fn: gradualLoadGrowthPostBurstSlowest, RequiredStack: "sleep", Description: generateDescription(Slowest, "", true)})

	Tests.Register(Test{Id: "GradualCPULoadGrowthPostBurstFastest", Fn: gradualLoadGrowthPostBurstFastest, RequiredStack: "cpustress", Description: generateDescription(Fastest, "CPU", true)})
	Tests.Register(Test{Id: "GradualCPULoadGrowthPostBurstFast", Fn: gradualLoadGrowthPostBurstFast, RequiredStack: "cpustress", Description: generateDescription(Fast, "CPU", true)})
	Tests.Register(Test{Id: "GradualCPULoadGrowthPostBurstMedium", Fn: gradualLoadGrowthPostBurstMedium, RequiredStack: "cpustress", Description: generateDescription(Medium, "CPU", true)})
	Tests.Register(Test{Id: "GradualCPULoadGrowthPostBurstSlow", Fn: gradualLoadGrowthPostBurstSlow, RequiredStack: "cpustress", Description: generateDescription(Slow, "CPU", true)})
	Tests.Register(Test{Id: "GradualCPULoadGrowthPostBurstSlowest", Fn: gradualLoadGrowthPostBurstSlowest, RequiredStack: "cpustress", Description: generateDescription(Slowest, "CPU", true)})

	Tests.Register(Test{Id: "GradualIOLoadGrowthPostBurstFastest", Fn: gradualLoadGrowthPostBurstFastest, RequiredStack: "iostress", Description: generateDescription(Fastest, "IO", true)})
	Tests.Register(Test{Id: "GradualIOLoadGrowthPostBurstFast", Fn: gradualLoadGrowthPostBurstFast, RequiredStack: "iostress", Description: generateDescription(Fast, "IO", true)})
	Tests.Register(Test{Id: "GradualIOLoadGrowthPostBurstMedium", Fn: gradualLoadGrowthPostBurstMedium, RequiredStack: "iostress", Description: generateDescription(Medium, "IO", true)})
	Tests.Register(Test{Id: "GradualIOLoadGrowthPostBurstSlow", Fn: gradualLoadGrowthPostBurstSlow, RequiredStack: "iostress", Description: generateDescription(Slow, "IO", true)})
	Tests.Register(Test{Id: "GradualIOLoadGrowthPostBurstSlowest", Fn: gradualLoadGrowthPostBurstSlowest, RequiredStack: "iostress", Description: generateDescription(Slowest, "IO", true)})

	Tests.Register(Test{Id: "GradualMemLoadGrowthPostBurstFastest", Fn: gradualLoadGrowthPostBurstFastest, RequiredStack: "memstress", Description: generateDescription(Fastest, "memory", true)})
	Tests.Register(Test{Id: "GradualMemLoadGrowthPostBurstFast", Fn: gradualLoadGrowthPostBurstFast, RequiredStack: "memstress", Description: generateDescription(Fast, "memory", true)})
	Tests.Register(Test{Id: "GradualMemLoadGrowthPostBurstMedium", Fn: gradualLoadGrowthPostBurstMedium, RequiredStack: "memstress", Description: generateDescription(Medium, "memory", true)})
	Tests.Register(Test{Id: "GradualMemLoadGrowthPostBurstSlow", Fn: gradualLoadGrowthPostBurstSlow, RequiredStack: "memstress", Description: generateDescription(Slow, "memory", true)})
	Tests.Register(Test{Id: "GradualMemLoadGrowthPostBurstSlowest", Fn: gradualLoadGrowthPostBurstSlowest, RequiredStack: "memstress", Description: generateDescription(Slowest, "memory", true)})
}

func gradualLoadGrowthFastest(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Fastest),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, nil)
}

func gradualLoadGrowthFast(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Fast),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, nil)
}

func gradualLoadGrowthMedium(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Medium),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, nil)
}

func gradualLoadGrowthSlow(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Slow),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, nil)
}

func gradualLoadGrowthSlowest(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Slowest),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, nil)
}

func gradualLoadGrowthPostBurstFastest(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Fastest),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, sendOpeningBurst)
}

func gradualLoadGrowthPostBurstFast(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Fast),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, sendOpeningBurst)
}

func gradualLoadGrowthPostBurstMedium(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Medium),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, sendOpeningBurst)
}

func gradualLoadGrowthPostBurstSlow(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Slow),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, sendOpeningBurst)
}

func gradualLoadGrowthPostBurstSlowest(test *config.Test) {
	gradualLoadGrowth(test, config.Http{
		TestType:  httpbench.RequestsForTimeGraph.String(),
		HitsGraph: gradualHitGraph(maxConcurrent, Slowest),
		Hook:      test.Config.Provider.HttpInvocationTriggerStage(),
	}, sendOpeningBurst)
}

func gradualLoadGrowth(test *config.Test, httpConfig config.Http, cb func(hfConf *config.HttpFunction)) {
	headers := http.Header{}
	body := []byte{}
	queryParams := url.Values{}
	httpConfig.QueryParams = &queryParams
	httpConfig.Headers = &headers
	httpConfig.Body = &body
	for _, function := range test.Stack.ListFunctions() {
		hfConf, err := test.NewFunction(&httpConfig, function)

		if err != nil {
			continue
		}

		if cb != nil {
			cb(hfConf)
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
		httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, hfConf.Test.Config.Provider.HttpInvocationLatency)
	}()
	requestsResult := trace.RequestsForTimeGraph(*hfConf.HttpConfig.HitsGraph)
	wg.Wait()
	httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
}

func sendOpeningBurst(hfConf *config.HttpFunction) {
	newReq := hfConf.Test.Config.Provider.NewFunctionRequest(hfConf.Test.Stack, hfConf.Function, hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
	lastHit := (*hfConf.HttpConfig.HitsGraph)[len(*hfConf.HttpConfig.HitsGraph)-1]
	traceToDiscard := httpbench.New(newReq, hfConf.HttpConfig.Hook)
	// we send roughly the same number of concurrent requests as at the peak time of our hits graph test
	requestsToSend := lastHit.Concurrent * uint64(hfConf.HttpConfig.SleepTime/lastHit.Time+1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for range traceToDiscard.ResultCh {
		}
		wg.Done()
	}()
	traceToDiscard.ConcurrentRequestsSyncedOnce(requestsToSend, 0)
	wg.Wait()
}
