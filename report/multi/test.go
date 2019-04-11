package multi

import (
	"github.com/nuweba/faasbenchmark/report"
	"github.com/pkg/errors"
)

type multiTest struct {
	test []report.Test
}

func (mt *multiTest) Description(desc string) error {
	for _, t := range mt.test {
		err := t.Description(desc)
		if err != nil {
			return errors.Wrap(err, "multi test description writer error")
		}
	}

	return nil
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
