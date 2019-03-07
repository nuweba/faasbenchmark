package file

import (
	"github.com/pkg/errors"
	"io"
	"github.com/nuweba/faasbenchmark/report"
	"os"
	"path/filepath"
)

const (
	FunctionTestResult = "summary.log"
	StackDescription   = "stack-description.txt"
	HttpTestConfig     = "http-test.conf"
	StackLogName = "test.log"
)

type Function struct {
	upperLevel         *Test
	functionName       string
	functionResultFile *os.File
	descriptionFile    *os.File
	httpTestConfigFile    *os.File
	logFile   *os.File
	functionResultPath string
}


func (test *Test) Function(functionName string) (report.Function, error) {
	f := &Function{upperLevel: test, functionName: functionName}

	testResultDir := filepath.Join(f.upperLevel.testResultPath, f.upperLevel.ProviderName)
	err := os.MkdirAll(testResultDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "function result dir should be unique")
	}
	f.functionResultPath = testResultDir
	return f, nil
}

func (f *Function) LogWriter() (io.Writer, error) {
	logPath := filepath.Join(f.functionResultPath, StackLogName)
	file, err := os.OpenFile(logPath,  os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	f.logFile = file
	return f.logFile, nil
}

func (f *Function) ResultWriter() (io.Writer, error) {
	functionResultFile := filepath.Join(f.functionResultPath, f.functionName + "_" + FunctionTestResult)
	file, err := os.OpenFile(functionResultFile,  os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "function result file already exists")
	}

	f.functionResultFile = file
	return f.functionResultFile, nil
}

func (f *Function) HttpTestConfigWriter() (io.Writer, error) {
	httpTestConfigFile := filepath.Join(f.functionResultPath, HttpTestConfig)
	file, err := os.OpenFile(httpTestConfigFile,  os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "http test config result file already exists")
	}

	f.httpTestConfigFile = file
	return f.httpTestConfigFile, nil
}

func (f *Function) DescriptionWriter() (io.Writer, error) {

	functionResultFile := filepath.Join(f.functionResultPath, StackDescription)
	file, err := os.Create(functionResultFile)
	if err != nil {
		return nil, errors.Wrap(err, "test description file")
	}
	f.descriptionFile = file
	return f.descriptionFile, nil
}