package graph

import (
	"io"
	"io/ioutil"
	"github.com/nuweba/faasbenchmark/report"
)

type Function struct {
	upperLevel         *Test
	functionName       string
}


func (test *Test) Function(functionName string) (report.Function, error) {
	f := &Function{upperLevel: test, functionName: functionName}

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
