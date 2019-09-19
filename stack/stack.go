package stack

type Stack interface {
	DeployStack() error
	RemoveStack() error
	StackId() string
	Project() string
	Stage() string
	ListFunctions() []Function
}

type Function interface {
	Name() string
	Handler() string
	Description() string
	Runtime() string
	MemorySize() string
}
