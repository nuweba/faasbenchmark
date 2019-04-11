package graph

import (
	"github.com/nuweba/faasbenchmark/report"
)

type Request struct {
	upperLevel *Function
}

func (f *Function) Request() (report.Request, error) {
	r := &Request{upperLevel: f}

	return r, nil
}

func (r *Request) Result(result string) error {
	_, err := r.upperLevel.upperLevel.upperLevel.graphWriter.Write([]byte(result))
	return err
}

func (r *Request) Summary(summary string) error {
	return nil
}

func (r *Request) Error(error string) error {
	return nil
}

func (r *Request) RawResult(raw string) error {
	return nil
}
