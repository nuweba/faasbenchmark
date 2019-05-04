package config

import (
	"encoding/json"
	"github.com/nuweba/httpbench"
	"github.com/nuweba/httpbench/syncedtrace"
	"net/http"
	"net/url"
	"time"
)

type Http struct {
	SleepTime        time.Duration
	Hook             syncedtrace.TraceHookType
	QueryParams      *url.Values
	Headers          *http.Header
	Duration         time.Duration
	RequestDelay     time.Duration
	ConcurrencyLimit uint64
	Body             *[]byte
	TestType         string
	ConcurrentGraph  *httpbench.ConcurrentGraph
	HitsGraph        *httpbench.HitsGraph
}

func (h *Http) String() (string, error) {
	b, err := json.MarshalIndent(h, "", "\t")
	if err != nil {
		return "", err
	}
	//b, err := yaml.Marshal(h)
	//if err != nil {
	//	return "", err
	//}
	return string(b), nil
}
