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
	Tests.Register(Test{Id: "LargeResponse", Fn: largeResponse, RequiredStack: "largeresponse", Description: "Benchmark the response time of a function returning a large (4mb) response body. Invoke the function once at a time for a minute."})
}

func largeResponse(test *config.Test) {
	headers := http.Header{}
	body := []byte{}
	queryParams := url.Values{}
	httpConfig := config.Http{
		QueryParams:      &queryParams,
		Headers:          &headers,
		Body:             &body,
		TestType:         httpbench.ConcurrentRequestsSynced.String(),
		Hook:             test.Config.Provider.HttpInvocationTriggerStage(),
		ConcurrencyLimit: 1,
		RequestDelay:     time.Millisecond,
		Duration:         time.Minute,
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
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, test.Config.Provider.HttpResult, test.Config.Debug)
		}()
		requestsResult := trace.ConcurrentRequestsSynced(httpConfig.ConcurrencyLimit, httpConfig.RequestDelay, httpConfig.Duration)
		wg.Wait()
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
	}
}
