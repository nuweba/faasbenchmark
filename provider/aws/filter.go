package aws

import (
	"github.com/nuweba/faasbenchmark/report"
	"github.com/nuweba/httpbench/engine"
	"github.com/nuweba/httpbench/syncedtrace"
	"time"
)

type Result struct {
	id                 uint64
	invocationOverHead time.Duration
	duration           time.Duration
	responseTime       time.Duration
	reused             bool
}

func (r *Result) Id() uint64 {
	return r.id
}

func (r *Result) InvocationOverHead() float64 {
	return float64(r.invocationOverHead)/float64(time.Millisecond)
}

func (r *Result) Duration() float64 {
	return float64(r.duration)/float64(time.Millisecond)
}

func (r *Result) ContentTransfer() float64 {
	return float64(r.responseTime)/float64(time.Millisecond)
}

func (r *Result) Reused() bool {
	return r.reused
}


func (aws *Aws) HttpResult(sleepTime time.Duration, tr *engine.TraceResult, funcDuration time.Duration, reused bool) (report.Result, error) {
	invocationOverHead := tr.Hooks[syncedtrace.GotFirstResponseByte].Duration - funcDuration
	responseTime := tr.Total - tr.Hooks[syncedtrace.GotFirstResponseByte].Duration
	duration := funcDuration - sleepTime

	r := &Result{
		invocationOverHead: invocationOverHead,
		duration:           duration,
		responseTime:       responseTime,
		reused:             reused,
		id:                 tr.Id,
	}

	return r, nil
}
