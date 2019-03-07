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

func (mf *multiFunction) HttpTestConfigWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mf.function {
		writer, err := t.HttpTestConfigWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi function http test config writer error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
}

func (mf *multiFunction) ResultWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mf.function {
		writer, err := t.ResultWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi function result writer error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
}

func (mf *multiFunction) DescriptionWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mf.function {
		writer, err := t.DescriptionWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi function description writer error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
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
