package tests

import (
	"math"
	"net/http"
	"net/url"
	"github.com/nuweba/faasbenchmark/config"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
	"nwb.nu/httpbench"
	"strconv"
	"sync"
	"time"
)

func init() {
	Tests.Register(Test{Id: "RequestsFor1Min", Fn: RequestsFor1Minute, RequiredStack: "coldstart", Description: "1 minute test of a new request every 500ms with no sleep"})
	Tests.Register(Test{Id: "ColdStart", Fn: coldStart, RequiredStack: "coldstart", Description: "Test cold start"})
}

func sleepQueryParam(sleep time.Duration) url.Values {
	qParams := url.Values{}
	sleepTimeMillisecond := strconv.FormatInt(int64(math.Ceil(float64(sleep.Nanoseconds())/float64(time.Millisecond))), 10)
	qParams.Add("sleep", sleepTimeMillisecond)

	return qParams
}

func coldStart(test *config.Test) {
	sleep := 2000 * time.Millisecond
	qParams := sleepQueryParam(sleep)
	headers := http.Header{}
	body := []byte("")

	httpConfig := &config.Http{
		SleepTime:        sleep,
		Hook:             test.Config.Provider.HttpInvocationTriggerStage(),
		QueryParams:      &qParams,
		Headers:          &headers,
		Duration:         0,
		RequestDelay:     20 * time.Millisecond,
		ConcurrencyLimit: 3,
		Body:             &body,
	}
	wg := &sync.WaitGroup{}
	for _, function := range test.Stack.ListFunctions() {
		hfConf, err := test.NewFunction(httpConfig, function)
		if err != nil {
			continue
		}

		newReq := test.Config.Provider.NewFunctionRequest(hfConf.Function.Name(), hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
		trace := httpbench.New(newReq, hfConf.HttpConfig.Hook)

		wg.Add(1)
		go func() {
			defer wg.Done()
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, test.Config.Provider.HttpInvocationLatency)
		}()

		requestsResult := trace.ConcurrentRequestsSyncedOnce(hfConf.HttpConfig.ConcurrencyLimit, hfConf.HttpConfig.RequestDelay)
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
		wg.Wait()
	}
}

func RequestsFor1Minute(test *config.Test) {
	sleep := 0 * time.Millisecond
	qParams := sleepQueryParam(sleep)
	headers := http.Header{}
	body := []byte("")

	httpConfig := &config.Http{
		SleepTime:        sleep,
		Hook:             test.Config.Provider.HttpInvocationTriggerStage(),
		QueryParams:      &qParams,
		Headers:          &headers,
		Duration:         1 * time.Minute,
		RequestDelay:     500 * time.Millisecond,
		ConcurrencyLimit: 0,
		Body:             &body,
	}

	wg := &sync.WaitGroup{}
	for _, function := range test.Stack.ListFunctions() {
		hfConf, err := test.NewFunction(httpConfig, function)
		if err != nil {
			continue
		}

		newReq := test.Config.Provider.NewFunctionRequest(hfConf.Function.Name(), hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
		trace := httpbench.New(newReq, hfConf.HttpConfig.Hook)

		wg.Add(1)
		go func() {
			defer wg.Done()
			httpbenchReport.ReportRequestResults(hfConf, trace.ResultCh, test.Config.Provider.HttpInvocationLatency)
		}()

		requestsResult := trace.RequestPerDuration(hfConf.HttpConfig.RequestDelay, hfConf.HttpConfig.Duration)
		httpbenchReport.ReportFunctionResults(hfConf, requestsResult)
		wg.Wait()
	}
}

//func warmStart(testId string, provider provider.FaasProvider, testStack *Stack, gConfig *GlobalConfig) {
//	rc, err := NewResultConfig(testId, gConfig, testStack, 2000 * time.Millisecond, httptrace.TLSHandshakeDone)
//
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	qParams := sleepQueryParam(rc.SleepTime)
//	headers := http.Header{}
//	//duration := 1 * time.Second
//	reqDelay := 0 * time.Millisecond
//	concurrencyLimit := uint64(3)
//
//	for _, function := range testStack.ListFunctions() {
//
//		fr, err := NewFunctionResult(rc, function)
//		if err != nil {
//			fmt.Println(err)
//		}
//		newReq := func() (*http.Request, error) {
//			return provider.NewFunctionRequest(function.Name(), qParams, headers, []byte(""))
//		}
//
//		//warming up functions
//		trace := httptrace.New(fr.DiscardResult, newReq)
//		trace.ConcurrentRequestsSyncedOnce(rc.Hook, concurrencyLimit, reqDelay)
//
//		fr, err = NewFunctionResult(rc, function)
//		if err != nil {
//			fmt.Println(err)
//		}
//		time.Sleep(time.Second)
//		trace = httptrace.New(fr.HandleHttpsFunctionResult, newReq)
//		requestsResult := trace.ConcurrentRequestsSyncedOnce(rc.Hook, concurrencyLimit, reqDelay)
//		fr.HandleReqResults(requestsResult)
//
//	}
//}
