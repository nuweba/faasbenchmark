package config

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/nuweba/faasbenchmark/provider"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	DescriptionFile = "description.txt"
)

type Stack struct {
	stack.Stack
	Description string
}

type Stacks struct {
	Stacks map[string]*Stack
}

func newStacks(provider provider.FaasProvider, arsenalPath string) (*Stacks, error) {
	stackPaths, err := getStackPaths(arsenalPath, provider.Name())
	if err != nil {
		return nil, errors.Wrap(err, "failed getting stack paths")
	}

	stacks := &Stacks{
		Stacks: make(map[string]*Stack),
	}

	for _, stackPath := range stackPaths {
		stack, err := newStack(provider, stackPath)

		if err != nil {
			return nil, errors.Wrap(err, "failed reading new stack")
		}

		if _, err := stacks.GetStack(stack.StackId()); err == nil {
			return nil, errors.New(fmt.Sprintf("duplicate stacks, %s", stack.StackId()))
		}

		stacks.Stacks[stack.StackId()] = stack
	}

	return stacks, nil
}

func newStack(provider provider.FaasProvider, stackPath string) (*Stack, error) {
	stack, err := provider.NewStack(stackPath)

	if err != nil {
		return nil, errors.Wrap(err, "failed parsing stack")
	}

	description, err := readDescription(stackPath)

	if err != nil {
		return nil, errors.Wrap(err, "read description failed")
	}

	return &Stack{Stack: stack, Description: description}, nil

}

func (stacks *Stacks) GetStack(stackId string) (*Stack, error) {
	if stack, ok := stacks.Stacks[stackId]; ok {
		return stack, nil
	}

	return nil, errors.New(fmt.Sprintf("stack not found %s", stackId))
}

func isStackFolder(info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	if info.Name() != DescriptionFile {
		return false
	}

	return true
}

func isProviderStackExists(files []os.FileInfo, provider string) bool {
	for _, file := range files {
		if file.Name() == provider {
			return true
		}
	}

	return false
}

func getStackPaths(arsenalPath string, provider string) ([]string, error) {
	var stackPaths []string

	err := filepath.Walk(arsenalPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !isStackFolder(info) {
				return nil
			}
			files, err := ioutil.ReadDir(filepath.Dir(path))
			if err != nil {
				return err
			}
			if !isProviderStackExists(files, provider) {
				return nil
			}
			stackPaths = append(stackPaths, filepath.Dir(path))

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return stackPaths, nil
}

func sanitizeDescription(data []byte) []byte {
	p := bluemonday.StrictPolicy()
	return p.SanitizeBytes(data)
}

func readDescription(testDirPath string) (string, error) {
	descriptionData, err := ioutil.ReadFile(filepath.Join(testDirPath, DescriptionFile))
	if err != nil {
		return "", err
	}

	Sanitized := sanitizeDescription(descriptionData)

	return string(Sanitized), nil
}
