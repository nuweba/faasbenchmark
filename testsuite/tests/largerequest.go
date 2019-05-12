package tests

import (
	"github.com/nuweba/faasbenchmark/config"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
	"github.com/nuweba/httpbench"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const parameterRepeats = 1042 * 1024 * 4

func init() {
	Tests.Register(Test{Id: "largeRequest", Fn: largeRequest, RequiredStack: "sleep", Description: "benchmark the invocation latency of a function invoked with a large request"})
}

func largeRequest(test *config.Test) {
	param := strings.Repeat("a", parameterRepeats) // This makes the request large
	headers := http.Header{}
	body := []byte{}
	queryParams := url.Values(map[string][]string{"param": {param}, "sleep": {"0"}})
	httpConfig := config.Http{
		QueryParams: &queryParams,
		Headers:     &headers,
		Body:        &body,
		TestType:    httpbench.ConcurrentRequestsSynced.String(),
		Hook:        test.Config.Provider.HttpInvocationTriggerStage(),
	}
	httpConfig.QueryParams = &queryParams
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
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, test.Config.Provider.HttpInvocationLatency)
		}()
		requestsResult := trace.ConcurrentRequestsSynced(1, time.Millisecond, benchmarkDuration)
		wg.Wait()
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
	}
}
