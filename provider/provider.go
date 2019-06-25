package provider

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/provider/aws"
	"github.com/nuweba/faasbenchmark/provider/azure"
	"github.com/nuweba/faasbenchmark/provider/google"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/httpbench/engine"
	"github.com/nuweba/httpbench/syncedtrace"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RequestFilter = func(sleepTime time.Duration, tr *engine.TraceResult, funcDuration time.Duration, reused bool) (string, error)

type Filter interface {
	HttpInvocationLatency(sleepTime time.Duration, tr *engine.TraceResult, funcDuration time.Duration, reused bool) (string, error)
}

type FaasProvider interface {
	Filter
	Name() string
	HttpInvocationTriggerStage() syncedtrace.TraceHookType
	NewStack(stackPath string) (stack.Stack, error)
	NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) func() (*http.Request, error)
}

type Providers int

const (
	AWS Providers = iota
	Google
	Azure
	ProvidersCount
)

func (p Providers) String() string {
	return [...]string{
		"aws",
		"google",
		"azure",
	}[p]
}

func (p Providers) Description() string {
	return [...]string{
		"aws lambda functions",
		"google cloud functions",
		"azure functions",
	}[p]
}

func NewProvider(providerName string) (FaasProvider, error) {
	var faasProvider FaasProvider
	var err error

	switch strings.ToLower(providerName) {
	case strings.ToLower(AWS.String()):
		faasProvider, err = aws.New()
	case strings.ToLower(Google.String()):
		faasProvider, err = google.New()
	case strings.ToLower(Azure.String()):
		faasProvider, err = azure.New()
	default:
		faasProvider, err = nil, errors.New(fmt.Sprintf("provider not supported: %s", providerName))
	}

	if err != nil {
		return nil, err
	}

	return faasProvider, nil
}

func ProviderList() []string {
	var providers []string
	for providerId := Providers(0); providerId < ProvidersCount; providerId++ {
		providers = append(providers, providerId.String())
	}

	return providers
}
