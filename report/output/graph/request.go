package graph

import (
	"io"
	"io/ioutil"
	"github.com/nuweba/faasbenchmark/report"
)


type Request struct {
	upperLevel *Function
}

func (f *Function) Request() (report.Request, error) {
	r := &Request{upperLevel: f}

	return r, nil
}

func (r *Request) ResultWriter() (io.Writer, error) {
	return r.upperLevel.upperLevel.upperLevel.graphWriter, nil
}

func (r *Request) SummaryWriter() (io.Writer, error) {
	return ioutil.Discard, nil
}

func (r *Request) ErrorWriter() (io.Writer, error) {
	return ioutil.Discard, nil
}

func (r *Request) RawResultWriter() (io.Writer, error) {
	return ioutil.Discard, nil
}
