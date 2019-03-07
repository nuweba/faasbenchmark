package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/nuweba/faasbenchmark/report"
)

type Test struct {
	Config          *Global
	Stack           *Stack
	Report          report.Test
	TestId          string
	TestDescription string
}

func (c *Global) NewTest(stack *Stack, testId string, testDescription string) (*Test, error) {
	testReport, err := c.report.Test(testId, c.Provider.Name())

	if err != nil {
		return nil, errors.Wrap(err, "test report failed")
	}

	t := &Test{
		Config: c,
		Stack:  stack,
		Report: testReport,
		TestId: testId,
		TestDescription: testDescription,
	}

	descriptionWriter, err := t.Report.DescriptionWriter()
	if err != nil {
		return nil, errors.Wrap(err, "test report description writer")
	}

	fmt.Fprintln(descriptionWriter, t.TestDescription)
	t.Config.Logger.Debug("wrote test description")

	return t, nil
}
