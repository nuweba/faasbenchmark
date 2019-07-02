package file

import (
	"github.com/nuweba/faasbenchmark/report"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

const (
	FunctionTestResult = "summary.log"
	StackDescription   = "stack-description.txt"
	HttpTestConfig     = "http-test.conf"
	StackLogName       = "test.log"
)

type Function struct {
	upperLevel         *Test
	functionName       string
	functionResultFile *os.File
	descriptionFile    *os.File
	httpTestConfigFile *os.File
	logFile            *os.File
	functionResultPath string
}

func (test *Test) Function(functionName, description, runtime, memorySize string) (report.Function, error) {
	f := &Function{upperLevel: test, functionName: functionName}

	//create provider dir inside the test dir
	testResultDir := filepath.Join(f.upperLevel.testResultPath, f.upperLevel.ProviderName)
	err := os.MkdirAll(testResultDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "function result dir should be unique")
	}
	f.functionResultPath = testResultDir

	//bench result
	functionResultFilePath := filepath.Join(f.functionResultPath, f.functionName+"_"+FunctionTestResult)
	functionResultFile, err := os.OpenFile(functionResultFilePath, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "function result file already exists")
	}

	f.functionResultFile = functionResultFile

	//stack description
	functionStackDescriptionFilePath := filepath.Join(f.functionResultPath, StackDescription)
	functionStackDescriptionFile, err := os.Create(functionStackDescriptionFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "test description file")
	}
	f.descriptionFile = functionStackDescriptionFile

	//http config
	httpTestConfigFilePath := filepath.Join(f.functionResultPath, HttpTestConfig)
	httpTestConfigFile, err := os.OpenFile(httpTestConfigFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "http test config result file already exists")
	}

	f.httpTestConfigFile = httpTestConfigFile

	return f, nil
}

func (f *Function) LogWriter() (io.Writer, error) {
	logPath := filepath.Join(f.functionResultPath, StackLogName)
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	f.logFile = file
	return f.logFile, nil
}

func (f *Function) BenchResult(bresult string) error {
	_, err := f.functionResultFile.WriteString(bresult)
	return err
}

func (f *Function) HttpTestConfig(config string) error {
	_, err := f.httpTestConfigFile.WriteString(config)
	return err
}

func (f *Function) StackDescription(sdesc string) error {
	_, err := f.descriptionFile.WriteString(sdesc)
	return err
}
