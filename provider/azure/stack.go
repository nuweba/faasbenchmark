package azure

import (
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/faasbenchmark/stack/sls"
	"path/filepath"
)

type Stack struct {
	*sls.Stack
}

func (azure *Azure) NewStack(stackPath string) (stack.Stack, error) {
	slsYamlDirPath := filepath.Join(stackPath, azure.Name())
	stack, err := sls.New(azure.Name(), slsYamlDirPath)

	if err != nil {
		return nil, err
	}

	return &Stack{stack}, nil
}

func (s *Stack) ListFunctions() []stack.Function {

	var functions []stack.Function

	funcs := s.ListFunctionsFromYaml()

	for _, f := range funcs {
		functions = append(functions, &Function{name: f.Name, handler: f.Handler, description: f.Description})
	}

	return functions
}
