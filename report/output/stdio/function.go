package stdio

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/report"
	"io"
	"io/ioutil"
)

type Function struct {
	upperLevel   *Test
	functionName string
}

func (test *Test) Function(functionName string) (report.Function, error) {
	f := &Function{upperLevel: test, functionName: functionName}

	fmt.Fprintln(test.upperLevel.stdoutWriter, functionName)
	return f, nil
}

func (f *Function) LogWriter() (io.Writer, error) {
	return ioutil.Discard, nil
}

func (f *Function) BenchResult(bresult string) error {
	return nil
}

func (f *Function) HttpTestConfig(config string) error {
	return nil
}

func (f *Function) StackDescription(sdesc string) error {
	return nil
}
