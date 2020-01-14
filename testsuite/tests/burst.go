package tests

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/config"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
	"github.com/nuweba/httpbench"
	"net/http"
	"sync"
	"time"
)

const (
	burstSizeLvl1 = 100
	burstSizeLvl2 = 500
	burstSizeLvl3 = 1000
)

func burstDescription(size uint64) string {
	return fmt.Sprintf("send a burst of %d invocations to a single function", size)
}

func init() {
	Tests.Register(Test{Id: "BurstLvl1", Fn: burstLvl1, RequiredStack: "singlefunction", Description: burstDescription(burstSizeLvl1)})
	Tests.Register(Test{Id: "BurstLvl1", Fn: burstLvl2, RequiredStack: "singlefunction", Description: burstDescription(burstSizeLvl2)})
	Tests.Register(Test{Id: "BurstLvl1", Fn: burstLvl3, RequiredStack: "singlefunction", Description: burstDescription(burstSizeLvl3)})
}

func burstLvl1(test *config.Test) {
	burst(test, burstSizeLvl1)
}

func burstLvl2(test *config.Test) {
	burst(test, burstSizeLvl2)
}

func burstLvl3(test *config.Test) {
	burst(test, burstSizeLvl3)
}

func burst(test *config.Test, size uint64) {
	sleep := 500 * time.Millisecond
	headers := http.Header{}
	body := []byte("")

	httpConfig := &config.Http{
		SleepTime:        sleep,
		Hook:             test.Config.Provider.HttpInvocationTriggerStage(),
		QueryParams:      sleepQueryParam(sleep),
		Headers:          &headers,
		Duration:         0,
		RequestDelay:     0,
		ConcurrencyLimit: size,
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
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, test.Config.Provider.HttpResult)
		}()

		requestsResult := trace.ConcurrentRequestsSyncedOnce(hfConf.HttpConfig.ConcurrencyLimit, hfConf.HttpConfig.RequestDelay)
		wg.Wait()
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
	}
}
