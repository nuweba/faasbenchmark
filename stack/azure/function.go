package azure

type Function struct {
	name        string
	handler     string
	description string
	runtime     string
	memorySize  string
}

func (f *Function) Name() string {
	return f.name
}

func (f *Function) Handler() string {
	return f.handler
}

func (f *Function) Description() string {
	return f.description
}

func (f *Function) Runtime() string {
	return f.runtime
}

func (f *Function) MemorySize() string {
	return f.memorySize
}
