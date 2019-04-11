package report

type Test interface {
	Description(desc string) error
	Function(functionName string) (Function, error)
}
