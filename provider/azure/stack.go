package azure

import (
	azurestack "github.com/nuweba/azure-stack"
	"github.com/nuweba/faasbenchmark/stack"
	az "github.com/nuweba/faasbenchmark/stack/azure"
	"path/filepath"
)

func (azure *Azure) NewStack(stackPath string) (stack.Stack, error) {
	stackPath = filepath.Join(stackPath, azure.Name())
	stack, err := azurestack.New(stackPath)

	if err != nil {
		return nil, err
	}

	return &az.Stack{AzureStack: stack}, nil
}
