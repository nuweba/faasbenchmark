package graph

import (
	"github.com/nuweba/faasbenchmark/report"
	"io"
	"io/ioutil"
)

type Top struct {
	graphWriter io.Writer
}

func New(graphWriter io.Writer) (report.Top, error) {

	t := &Top{graphWriter: graphWriter}

	return t, nil
}

func (t *Top) LogWriter() (io.Writer, error) {
	return ioutil.Discard, nil
}
