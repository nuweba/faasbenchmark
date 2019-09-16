package json

import (
	"encoding/json"
	"github.com/nuweba/faasbenchmark/report"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	jsonName           = "result.json"
)

type Function struct {
	upperLevel         *Test
	functionName       string
	functionResultPath string
	jsonFile           *os.File
	json               *functionJson
}

type functionJson struct {
	FunctionName string        `json:"functionName"`
	Description  string        `json:"description"`
	Runtime      string        `json:"runtime"`
	MemorySize   string        `json:"memorySize"`
	Results      []requestJson `json:"results"`
}

func (fj *functionJson) AddResult(result requestJson) {
	fj.Results = append(fj.Results, result)
}

func (test *Test) Function(functionName, description, runtime, memorySize string) (report.Function, error) {
	f := &Function{upperLevel: test, functionName: functionName}

	testResultDir := filepath.Join(f.upperLevel.testResultPath, f.upperLevel.ProviderName)
	err := os.MkdirAll(testResultDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, errors.Wrap(err, "could not create function log dir")
	}
	f.functionResultPath = testResultDir

	//json
	f.json = &functionJson{
		FunctionName: functionName,
		Description:  description,
		Runtime:      runtime,
		MemorySize:   memorySize,
	}

	//json file
	jsonFilePath := filepath.Join(f.functionResultPath, jsonName)
	jsonFile, err := os.OpenFile(jsonFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "json result file")
	}

	f.jsonFile = jsonFile

	return f, nil
}

func (f *Function) LogWriter() (io.Writer, error) {
	return ioutil.Discard, nil
}

func (f *Function) BenchResult(bresult string) error {
	f.upperLevel.json.AddFunction(f.json)

	b, err := json.MarshalIndent(f.upperLevel.json, "", "\t")
	if err != nil {
		return err
	}

	_, err = f.jsonFile.Write(b)
	if err != nil {
		return err
	}
	return err
}

func (f *Function) HttpTestConfig(config string) error {
	rawConfig := json.RawMessage(config)
	f.upperLevel.json.HttpConfig = &rawConfig
	return nil
}

func (f *Function) StackDescription(sdesc string) error {
	f.upperLevel.json.StackDescription = sdesc
	return nil
}
