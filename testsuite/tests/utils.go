package tests

import (
	"github.com/nuweba/faasbenchmark/config"
	httpbenchReport "github.com/nuweba/faasbenchmark/report/generate/httpbench"
	"github.com/nuweba/httpbench"
	"math"
	"net/url"
	"strconv"
	"sync"
	"time"
)

func sleepQueryParam(sleep time.Duration) *url.Values {
	qParams := url.Values{}
	sleepTimeMillisecond := strconv.FormatInt(int64(math.Ceil(float64(sleep.Nanoseconds())/float64(time.Millisecond))), 10)
	qParams.Add("sleep", sleepTimeMillisecond)

	return &qParams
}

func gradualHitGraph(maxConcurrent int, durationIntensity time.Duration) *httpbench.HitsGraph {
	var graph httpbench.HitsGraph
	for concurrent := 1; concurrent <= maxConcurrent; concurrent++ {
		graph = append(graph, httpbench.RequestsPerTime{Concurrent: uint64(concurrent), Time: durationIntensity})
	}
	return &graph
}

func sendPreWarmup(hfConf *config.HttpFunction, requestsToSend uint64) {
	newReq := hfConf.Test.Config.Provider.NewFunctionRequest(hfConf.Test.Stack, hfConf.Function, hfConf.HttpConfig.QueryParams, hfConf.HttpConfig.Headers, hfConf.HttpConfig.Body)
	trace := httpbench.New(newReq, hfConf.HttpConfig.Hook)
	// we send roughly the same number of concurrent requests as at the peak time of our hits graph test
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		httpbenchReport.DiscardRequestResults(trace.ResultCh)
		wg.Done()
	}()
	trace.ConcurrentRequestsSyncedOnce(requestsToSend, 0)
	wg.Wait()
}

const benchmarkDuration = 1 * time.Minute
