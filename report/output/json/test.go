package json

import (
	"encoding/json"
	"github.com/nuweba/faasbenchmark/report"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

type Test struct {
	upperLevel     *Top
	testResultPath string
	testId         string
	ProviderName   string
	json           *testJson
}

type testJson struct {
	Provider         string           `json:"provider"`
	TestName         string           `json:"testName"`
	TestDescription  string           `json:"testDescription"`
	StackDescription string           `json:"stackDescription"`
	HttpConfig       *json.RawMessage `json:"httpConfig"`
	Functions        []functionJson   `json:"functions"`
}

func (tj *testJson) AddFunction(function *functionJson) {
	tj.Functions = append(tj.Functions, *function)
}

func (t *Top) Test(testId string, providerName string) (report.Test, error) {
	test := &Test{upperLevel: t, testId: testId, ProviderName: providerName}

	testResultDir := filepath.Join(test.upperLevel.reportDir, test.testId)
	err := os.MkdirAll(testResultDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, errors.Wrap(err, "could not create test log dir")
	}

	test.testResultPath = testResultDir

	//json
	test.json = &testJson{
		Provider: providerName,
		TestName: testId,
	}

	return test, nil
}

func (test *Test) Description(desc string) error {
	test.json.TestDescription = desc
	return nil
}
