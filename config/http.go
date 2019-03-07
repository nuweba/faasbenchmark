package config

import (
	"gopkg.in/yaml.v2"
	"net/http"
	"net/url"
	"nwb.nu/httpbench/syncedtrace"
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
}

func (h *Http) String() (string, error) {
	b, err := yaml.Marshal(h)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
