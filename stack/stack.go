package stack

type Stack interface {
	DeployStack() error
	RemoveStack() error
	GetStackId() string
	ListFunctions() []Function
}

type Function interface {
	Name() string
	Description() string
}
