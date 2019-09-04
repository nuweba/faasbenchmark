package google

import (
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/faasbenchmark/stack/sls"
	"path/filepath"
)

type Stack struct {
	*sls.Stack
}

func (google *Google) NewStack(stackPath string) (stack.Stack, error) {
	slsYamlDirPath := filepath.Join(stackPath, google.Name())
	stack, err := sls.New(google.Name(), slsYamlDirPath)

	if err != nil {
		return nil, err
	}

	return &Stack{stack}, nil
}

func (s *Stack) ListFunctions() []stack.Function {

	var functions []stack.Function

	funcs := s.ListFunctionsFromYaml()

	for _, f := range funcs {
		functions = append(functions, &Function{name: f.Name,
			handler:     f.Handler,
			description: f.Description,
			memorySize:  f.MemorySize,
			runtime:     f.Runtime})
	}

	return functions
}
