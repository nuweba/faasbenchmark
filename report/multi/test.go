package multi

import (
	"github.com/pkg/errors"
	"io"
	"github.com/nuweba/faasbenchmark/report"
)

type multiTest struct {
	test []report.Test
}

func (mt *multiTest) DescriptionWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mt.test {
		writer, err := t.DescriptionWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi test description writer error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
}

func (mt *multiTest) Function(functionName string) (report.Function, error) {
	multiFunction := &multiFunction{}
	for _, t := range mt.test {
		function, err := t.Function(functionName)
		if err != nil {
			return nil, errors.Wrap(err, "multi function error")
		}
		multiFunction.function = append(multiFunction.function, function)
	}

	return multiFunction, nil
}
