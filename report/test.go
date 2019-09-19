package report

type Test interface {
	Description(desc string) error
	Function(functionName, description, runtime, memorySize string) (Function, error)
}
