package file

import (
	"github.com/pkg/errors"
	"io"
	"github.com/nuweba/faasbenchmark/report"
	"os"
	"path/filepath"
)

const (
	TestDescription = "description.txt"
)

type Test struct {
	upperLevel        *Top
	testResultPath    string
	testId            string
	ProviderName      string
	descriptionWriter *os.File
}

func (t *Top) Test(testId string, providerName string) (report.Test, error) {
	test := &Test{upperLevel: t, testId: testId, ProviderName: providerName}

	testResultDir := filepath.Join(test.upperLevel.reportDir, test.testId)
	err := os.MkdirAll(testResultDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "test dir should be unique")
	}

	test.testResultPath = testResultDir
	return test, nil
}

func (test *Test) DescriptionWriter() (io.Writer, error) {

	functionResultFile := filepath.Join(test.testResultPath, TestDescription)
	f, err := os.Create(functionResultFile)
	if err != nil {
		return nil, errors.Wrap(err, "test description file")
	}
	test.descriptionWriter = f
	return f, nil
}
