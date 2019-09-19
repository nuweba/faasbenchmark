package tests

import (
	"errors"
	"fmt"
	"github.com/nuweba/faasbenchmark/config"
	"os"
)

type Test struct {
	Id            string
	Fn            func(*config.Test)
	RequiredStack string
	Description   string
}

type testSuites struct {
	TestFunctions map[string]Test
}

func (ts *testSuites) Register(testFunction ...Test) {
	for _, test := range testFunction {
		if _, ok := ts.TestFunctions[test.Id]; ok {
			fmt.Printf("duplicate Test functions, %s", test.Id)
			os.Exit(1)
		}
		ts.TestFunctions[test.Id] = test
	}
}

func (ts *testSuites) GetTestSuite(testSuiteId string) (*Test, error) {
	if test, ok := ts.TestFunctions[testSuiteId]; ok {
		return &test, nil
	}

	return nil, errors.New(fmt.Sprintf("Test suite not found %s", testSuiteId))
}

var Tests = &testSuites{make(map[string]Test)}
