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
	SleepTime        time.Duration              `json:"sleepTime"`
	Hook             syncedtrace.TraceHookType  `json:"hook"`
	QueryParams      *url.Values                `json:"queryParams"`
	Headers          *http.Header               `json:"headers"`
	Duration         time.Duration              `json:"duration"`
	RequestDelay     time.Duration              `json:"requestDelay"`
	ConcurrencyLimit uint64                     `json:"concurrencyLimit"`
	Body             *[]byte                    `json:"body"`
	TestType         string                     `json:"testType"`
	ConcurrentGraph  *httpbench.ConcurrentGraph `json:"concurrentGraph"`
	HitsGraph        *httpbench.HitsGraph       `json:"hitsGraph"`
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
