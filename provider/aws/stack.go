package aws

import (
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/faasbenchmark/stack/sls"
	"path/filepath"
)

type Stack struct {
	*sls.Stack
}

func (aws *Aws) NewStack(stackPath string) (stack.Stack, error) {
	slsYamlDirPath := filepath.Join(stackPath, aws.Name())
	stack, err := sls.New(aws.Name(), slsYamlDirPath)

	if err != nil {
		return nil, err
	}

	return &Stack{stack}, nil
}
