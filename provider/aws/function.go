package aws

type Function struct {
	name        string
	description string
}

func (f *Function) Name() string {
	return f.name
}

func (f *Function) Description() string {
	return f.description
}
