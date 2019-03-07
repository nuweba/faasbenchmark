package report

import (
	"io"
)

type Test interface {
	DescriptionWriter() (io.Writer, error)
	Function(functionName string) (Function, error)
}
