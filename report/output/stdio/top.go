package stdio

import (
	"io"
	"io/ioutil"
	"github.com/nuweba/faasbenchmark/report"
)

type Top struct {
	stdoutWriter   io.Writer
}

func New(stdoutWriter io.Writer) (report.Top, error) {

	t := &Top{stdoutWriter: stdoutWriter}

	return t, nil
}

func (t *Top) LogWriter() (io.Writer, error) {
	return ioutil.Discard, nil
}
