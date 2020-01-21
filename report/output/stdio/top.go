package stdio

import (
	"github.com/nuweba/faasbenchmark/report"
	"io"
	"os"
)

type Top struct {
	stdoutWriter io.Writer
}

func New(stdoutWriter io.Writer) (report.Top, error) {

	t := &Top{stdoutWriter: stdoutWriter}

	return t, nil
}

func (t *Top) LogWriter() (io.Writer, error) {
	return os.Stdout, nil
}
