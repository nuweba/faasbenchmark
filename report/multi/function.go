package multi

import (
	"github.com/pkg/errors"
	"io"
	"github.com/nuweba/faasbenchmark/report"
)

type multiFunction struct {
	function []report.Function
}

func (mf *multiFunction) LogWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mf.function {
		writer, err := t.LogWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi function log writer error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
}

func (mf *multiFunction) HttpTestConfig(config string) error {
	for _, t := range mf.function {
		err := t.HttpTestConfig(config)
		if err != nil {
			return errors.Wrap(err, "multi function http test config writer error")
		}
	}

	return nil
}

func (mf *multiFunction) BenchResult(bresult string) error {
	for _, t := range mf.function {
		err := t.BenchResult(bresult)
		if err != nil {
			return errors.Wrap(err, "multi function result writer error")
		}
	}

	return nil
}

func (mf *multiFunction) StackDescription(sdesc string) error {
	for _, t := range mf.function {
		err := t.StackDescription(sdesc)
		if err != nil {
			return  errors.Wrap(err, "multi function description writer error")
		}
	}

	return nil
}

func (mf *multiFunction) Request() (report.Request, error) {
	multiRequest := &multiRequest{}
	for _, t := range mf.function {
		function, err := t.Request()
		if err != nil {
			return nil, errors.Wrap(err, "multi request error")
		}
		multiRequest.request = append(multiRequest.request, function)
	}

	return multiRequest, nil
}
