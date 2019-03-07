package aws

import (
	"fmt"
	"nwb.nu/httpbench/engine"
	"nwb.nu/httpbench/syncedtrace"
	"time"
)

func (aws *Aws) HttpInvocationLatency(sleepTime time.Duration,tr *engine.TraceResult,funcDuration time.Duration, reused bool) (string, error) {
	coldStart := float64(tr.Hooks[syncedtrace.GotFirstResponseByte].Duration) - float64(funcDuration)
	s := fmt.Sprintf("%f", coldStart/float64(time.Millisecond))

	return s, nil
}
