package provider

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"github.com/nuweba/faasbenchmark/provider/aws"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/httpbench/engine"
	"github.com/nuweba/httpbench/syncedtrace"
	"time"
)

type RequestFilter = func(sleepTime time.Duration,tr *engine.TraceResult,funcDuration time.Duration, reused bool) (string, error)

type Filter interface {
	HttpInvocationLatency(sleepTime time.Duration,tr *engine.TraceResult,funcDuration time.Duration, reused bool) (string, error)
}

type FaasProvider interface {
	Filter
	Name() string
	HttpInvocationTriggerStage() syncedtrace.TraceHookType
	NewStack(stackPath string) (stack.Stack, error)
	NewFunctionRequest(funcName string, qParams *url.Values, headers *http.Header, body *[]byte) (func () (*http.Request, error))
}

type Providers int

const (
	AWS Providers = iota
	ProvidersCount
)

func (p Providers) String() string {
	return [...]string{
		"aws",
	}[p]
}

func (p Providers) Description() string {
	return [...]string{
		"aws lambda functions",
	}[p]
}

func NewProvider(providerName string) (FaasProvider, error) {
	var faasProvider FaasProvider
	var err error

	switch providerName {
	case AWS.String():
		faasProvider, err = aws.New()
	default:
		faasProvider, err = nil, errors.New(fmt.Sprintf("provider not supported: %s", providerName))
	}

	if err != nil {
		return nil, err
	}

	return faasProvider, nil
}
