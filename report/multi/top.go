package multi

import (
	"github.com/pkg/errors"
	"io"
	"github.com/nuweba/faasbenchmark/report"
)

type multiReporter struct {
	top []report.Top
}

func Report(t ...report.Top) report.Top {
	return &multiReporter{t}
}

func (mr *multiReporter) LogWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mr.top {
		writer, err := t.LogWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi report logwriter error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
}

func (mr *multiReporter) Test(testId string, providerName string) (report.Test, error) {
	multiTest := &multiTest{}
	for _, t := range mr.top {
		test, err := t.Test(testId, providerName)
		if err != nil {
			return nil, errors.Wrap(err, "multi test reporter error")
		}
		multiTest.test = append(multiTest.test, test)
	}

	return multiTest, nil
}
