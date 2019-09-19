package graph

import (
	"github.com/nuweba/faasbenchmark/report"
)

type Test struct {
	upperLevel   *Top
	testId       string
	ProviderName string
}

func (t *Top) Test(testId string, providerName string) (report.Test, error) {
	test := &Test{upperLevel: t, testId: testId, ProviderName: providerName}
	return test, nil
}

func (test *Test) Description(desc string) error {
	return nil
}
