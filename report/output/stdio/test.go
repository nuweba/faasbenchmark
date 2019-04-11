package stdio

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/report"
)

type Test struct {
	upperLevel   *Top
	testId       string
	ProviderName string
}

func (t *Top) Test(testId string, providerName string) (report.Test, error) {
	test := &Test{upperLevel: t, testId: testId, ProviderName: providerName}
	fmt.Fprintln(t.stdoutWriter, testId)
	return test, nil
}

func (test *Test) Description(desc string) error {
	_, err := test.upperLevel.stdoutWriter.Write([]byte(desc))
	return err
}
