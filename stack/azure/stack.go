package azure

import (
	"github.com/nuweba/azure-stack"
	"github.com/nuweba/faasbenchmark/stack"
)

type Stack struct {
	azurestack.AzureStack
}

func New(path string) (*Stack, error) {
	stack, err := azurestack.New(path)

	if err != nil {
		return nil, err
	}

	return &Stack{stack}, nil
}

func (s *Stack) ListFunctions() []stack.Function {

	var functions []stack.Function

	for _, f := range s.Functions {
		nf := &Function{
			name:        f.Name,
			handler:     f.Handler,
			description: f.Description,
			runtime:     f.Runtime,
			memorySize:  f.MemorySize,
		}
		functions = append(functions, nf)
	}

	return functions
}
