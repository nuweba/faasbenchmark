package tests

import (
	"github.com/nuweba/httpbench"
	"math"
	"net/url"
	"strconv"
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
