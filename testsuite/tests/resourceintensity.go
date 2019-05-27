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
	descTemplate = "Invoke one %s intensive function (intensity level %d) at a time for %d minute(s) and benchmark the runtime duration."
)

func init() {
	Tests.Register(Test{Id: "CPUIntensiveLvl1", Fn: resourceIntensityLvl1, RequiredStack: "cpustress", Description: fmt.Sprintf(descTemplate, "CPU", 1, 1)})
	Tests.Register(Test{Id: "CPUIntensiveLvl2", Fn: resourceIntensityLvl2, RequiredStack: "cpustress", Description: fmt.Sprintf(descTemplate, "CPU", 2, 1)})
	Tests.Register(Test{Id: "CPUIntensiveLvl3", Fn: resourceIntensityLvl3, RequiredStack: "cpustress", Description: fmt.Sprintf(descTemplate, "CPU", 3, 1)})

	Tests.Register(Test{Id: "IOIntensiveLvl1", Fn: resourceIntensityLvl1, RequiredStack: "iostress", Description: fmt.Sprintf(descTemplate, "IO", 1, 1)})
	Tests.Register(Test{Id: "IOIntensiveLvl2", Fn: resourceIntensityLvl2, RequiredStack: "iostress", Description: fmt.Sprintf(descTemplate, "IO", 2, 1)})
	Tests.Register(Test{Id: "IOIntensiveLvl3", Fn: resourceIntensityLvl3, RequiredStack: "iostress", Description: fmt.Sprintf(descTemplate, "IO", 3, 1)})

	Tests.Register(Test{Id: "MemIntensiveLvl1", Fn: resourceIntensityLvl1, RequiredStack: "memstress", Description: fmt.Sprintf(descTemplate, "memory", 1, 1)})
	Tests.Register(Test{Id: "MemIntensiveLvl2", Fn: resourceIntensityLvl2, RequiredStack: "memstress", Description: fmt.Sprintf(descTemplate, "memory", 2, 1)})
	Tests.Register(Test{Id: "MemIntensiveLvl3", Fn: resourceIntensityLvl3, RequiredStack: "memstress", Description: fmt.Sprintf(descTemplate, "memory", 3, 1)})

	Tests.Register(Test{Id: "LogIntensiveLvl1", Fn: resourceIntensityLvl1, RequiredStack: "logging", Description: fmt.Sprintf(descTemplate, "logging", 1, 1)})
	Tests.Register(Test{Id: "LogIntensiveLvl2", Fn: resourceIntensityLvl2, RequiredStack: "logging", Description: fmt.Sprintf(descTemplate, "logging", 2, 1)})
	Tests.Register(Test{Id: "LogIntensiveLvl3", Fn: resourceIntensityLvl3, RequiredStack: "logging", Description: fmt.Sprintf(descTemplate, "logging", 3, 1)})

	Tests.Register(Test{Id: "NetIntensiveLvl1", Fn: resourceIntensityLvl1, RequiredStack: "network", Description: fmt.Sprintf(descTemplate, "network", 1, 1)})
	Tests.Register(Test{Id: "NetIntensiveLvl2", Fn: resourceIntensityLvl2, RequiredStack: "network", Description: fmt.Sprintf(descTemplate, "network", 2, 1)})
	Tests.Register(Test{Id: "NetIntensiveLvl3", Fn: resourceIntensityLvl3, RequiredStack: "network", Description: fmt.Sprintf(descTemplate, "network", 3, 1)})

	Tests.Register(Test{Id: "InternalService", Fn: resourceIntensityLvl1, RequiredStack: "internalservice", Description: "benchmark the runtime duration of a function that accesses an internal provider service. Invokes the function once at a time for 1 minute(s)."})
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
	httpConfig.ConcurrencyLimit = 1
	httpConfig.RequestDelay = time.Millisecond
	httpConfig.Duration = time.Minute

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
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, duration)
		}()
		requestsResult := trace.ConcurrentRequestsSynced(httpConfig.ConcurrencyLimit, httpConfig.RequestDelay, httpConfig.Duration)
		wg.Wait()
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
	}
}
