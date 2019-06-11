package graph

import (
	"encoding/json"
	"github.com/nuweba/faasbenchmark/report"
)

type Result struct {
	Id                 uint64
	InvocationOverHead float64
	Duration           float64
	ContentTransfer       float64
	Reused             bool
}

type Request struct {
	upperLevel *Function
}

func (f *Function) Request() (report.Request, error) {
	r := &Request{upperLevel: f}

	return r, nil
}

func (r *Request) Result(result report.Result) error {
	//todo: nasty
	res := &Result{
		Id: result.Id(),
		InvocationOverHead: result.InvocationOverHead(),
		Duration: result.Duration(),
		ContentTransfer: result.ContentTransfer(),
		Reused: result.Reused(),
	}
	b, err := json.Marshal(res)
	_, err = r.upperLevel.upperLevel.upperLevel.graphWriter.Write([]byte(b))
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
