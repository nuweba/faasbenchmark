package azure

import (
	"fmt"
	"github.com/nuweba/httpbench/engine"
	"github.com/nuweba/httpbench/syncedtrace"
	"time"
)

func (azure *Azure) HttpInvocationLatency(sleepTime time.Duration, tr *engine.TraceResult, funcDuration time.Duration, reused bool) (string, error) {
	coldStart := float64(tr.Hooks[syncedtrace.GotFirstResponseByte].Duration) - float64(funcDuration)
	s := fmt.Sprintf("%f", coldStart/float64(time.Millisecond))

	return s, nil
}
