package sls

import (
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/sls"
)

type Stack struct {
	*sls.Wrapper
}

func New(provider string, slsDirPath string) (*Stack, error) {
	stack, err := sls.New(provider, slsDirPath)

	if err != nil {
		return nil, err
	}

	return &Stack{stack}, nil
}

func (s *Stack) ListFunctions() []stack.Function {

	var functions []stack.Function

	funcs := s.ListFunctionsFromYaml()

	for _, f := range funcs {
		functions = append(functions, &Function{name: f.Name, description: f.Description})
	}

	return functions
}
