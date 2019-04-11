package sls

type Function struct {
	name        string
	handler     string
	description string
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
