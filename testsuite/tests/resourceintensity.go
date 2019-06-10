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
	descTemplate = "Invoke one %s intensive function (intensity level %d) at a time for %.0f minute(s)."
)

func init() {
	Tests.Register(Test{Id: "IntensiveCPUDurationLvl1", Fn: resourceIntensityLvl1, RequiredStack: "cpustress", Description: fmt.Sprintf(descTemplate, "CPU", 1, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveCPUDurationLvl2", Fn: resourceIntensityLvl2, RequiredStack: "cpustress", Description: fmt.Sprintf(descTemplate, "CPU", 2, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveCPUDurationLvl3", Fn: resourceIntensityLvl3, RequiredStack: "cpustress", Description: fmt.Sprintf(descTemplate, "CPU", 3, benchmarkDuration.Minutes())})

	Tests.Register(Test{Id: "IntensiveIODurationLvl1", Fn: resourceIntensityLvl1, RequiredStack: "iostress", Description: fmt.Sprintf(descTemplate, "IO", 1, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveIODurationLvl2", Fn: resourceIntensityLvl2, RequiredStack: "iostress", Description: fmt.Sprintf(descTemplate, "IO", 2, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveIODurationLvl3", Fn: resourceIntensityLvl3, RequiredStack: "iostress", Description: fmt.Sprintf(descTemplate, "IO", 3, benchmarkDuration.Minutes())})

	Tests.Register(Test{Id: "IntensiveMemDurationLvl1", Fn: resourceIntensityLvl1, RequiredStack: "memstress", Description: fmt.Sprintf(descTemplate, "memory", 1, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveMemDurationLvl2", Fn: resourceIntensityLvl2, RequiredStack: "memstress", Description: fmt.Sprintf(descTemplate, "memory", 2, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveMemDurationLvl3", Fn: resourceIntensityLvl3, RequiredStack: "memstress", Description: fmt.Sprintf(descTemplate, "memory", 3, benchmarkDuration.Minutes())})

	Tests.Register(Test{Id: "IntensiveLoggingDurationLvl1", Fn: resourceIntensityLvl1, RequiredStack: "logging", Description: fmt.Sprintf(descTemplate, "logging", 1, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveLoggingDurationLvl2", Fn: resourceIntensityLvl2, RequiredStack: "logging", Description: fmt.Sprintf(descTemplate, "logging", 2, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveLoggingDurationLvl3", Fn: resourceIntensityLvl3, RequiredStack: "logging", Description: fmt.Sprintf(descTemplate, "logging", 3, benchmarkDuration.Minutes())})

	Tests.Register(Test{Id: "IntensiveNetDurationLvl1", Fn: resourceIntensityLvl1, RequiredStack: "netstress", Description: fmt.Sprintf(descTemplate, "netstress", 1, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveNetDurationLvl2", Fn: resourceIntensityLvl2, RequiredStack: "netstress", Description: fmt.Sprintf(descTemplate, "netstress", 2, benchmarkDuration.Minutes())})
	Tests.Register(Test{Id: "IntensiveNetDurationLvl3", Fn: resourceIntensityLvl3, RequiredStack: "netstress", Description: fmt.Sprintf(descTemplate, "netstress", 3, benchmarkDuration.Minutes())})
}

func resourceIntensityLvl1(test *config.Test) {
	params := url.Values(map[string][]string{"level": {"1"}})
	resourceIntensity(test, config.Http{
		QueryParams: &params,
		TestType:    httpbench.ConcurrentRequestsSynced.String(),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func resourceIntensityLvl2(test *config.Test) {
	params := url.Values(map[string][]string{"level": {"2"}})
	resourceIntensity(test, config.Http{
		QueryParams: &params,
		TestType:    httpbench.ConcurrentRequestsSynced.String(),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func resourceIntensityLvl3(test *config.Test) {
	params := url.Values(map[string][]string{"level": {"3"}})
	resourceIntensity(test, config.Http{
		QueryParams: &params,
		TestType:    httpbench.ConcurrentRequestsSynced.String(),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	})
}

func resourceIntensity(test *config.Test, httpConfig config.Http) {
	headers := http.Header{}
	body := []byte{}
	httpConfig.Headers = &headers
	httpConfig.Body = &body

	for _, function := range test.Stack.ListFunctions() {
		hfConf, err := test.NewFunction(&httpConfig, function)

		if err != nil {
			continue
		}

		newReq := hfConf.Test.Config.Provider.NewFunctionRequest(hfConf.Test.Stack, hfConf.Function, hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
		wg := &sync.WaitGroup{}
		trace := httpbench.New(newReq, hfConf.HttpConfig.Hook)
		wg.Add(1)
		go func() {
			defer wg.Done()
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, test.Config.Provider.HttpResult)
		}()
		requestsResult := trace.ConcurrentRequestsSynced(1, time.Millisecond, benchmarkDuration)
		wg.Wait()
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
	}
}
