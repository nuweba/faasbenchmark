package graph

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/report"
)

type Request struct {
	upperLevel *Function
}

func (f *Function) Request() (report.Request, error) {
	r := &Request{upperLevel: f}

	return r, nil
}

func (r *Request) Result(result report.Result) error {
	_, err := r.upperLevel.upperLevel.upperLevel.graphWriter.Write([]byte(fmt.Sprintf("%s %s", result.Id() ,result.InvocationOverHead())))
	return err
}

func (r *Request) Summary(summary string) error {
	return nil
}

func (r *Request) Error(id uint64, error string) error {
	return nil
}

func (r *Request) RawResult(raw string) error {
	return nil
}
