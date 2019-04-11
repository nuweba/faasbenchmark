package report

import (
	"io"
)

type Function interface {
	LogWriter() (io.Writer, error)
	BenchResult(bresult string) error
	StackDescription(sdesc string) error
	HttpTestConfig(config string) error
	Request() (Request, error)
}
