package stdio

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
	return nil
}

func (r *Request) Summary(summary string) error {
	_, err := r.upperLevel.upperLevel.upperLevel.stdoutWriter.Write([]byte(summary))
	return err
}

func (r *Request) Error(error string) error {
	_, err := r.upperLevel.upperLevel.upperLevel.stdoutWriter.Write([]byte(error))
	return err
}

func (r *Request) RawResult(raw string) error {
	return nil
}
